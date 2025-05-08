package restache

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type pluginConfig struct {
	extName     string
	tagPrefixes map[string]string
	tagMappings map[string]string
}

type PluginOption func(*pluginConfig)

func WithTagPrefixes(tagPrefixes map[string]string) PluginOption {
	return func(cfg *pluginConfig) {
		cfg.tagPrefixes = tagPrefixes
	}
}

func WithTagMappings(tagMappings map[string]string) PluginOption {
	return func(cfg *pluginConfig) {
		cfg.tagMappings = tagMappings
	}
}

func WithExtensionName(extName string) PluginOption {
	extName = sanitizeExtensionName(extName)
	return func(cfg *pluginConfig) {
		cfg.extName = extName
	}
}

func readPluginConfig(cfg *pluginConfig, opts ...PluginOption) {
	for _, opt := range opts {
		opt(cfg)
	}
}

type plugin struct {
	cfg          *pluginConfig
	buildOptions *api.BuildOptions
	resolveFunc  func(path string, options api.ResolveOptions) api.ResolveResult
}

func (p *plugin) resolvePath(path string, resolveDir string) (string, bool, error) {
	result := p.resolveFunc(path, api.ResolveOptions{
		Kind:       api.ResolveJSImportStatement,
		ResolveDir: resolveDir,
	})
	if len(result.Errors) > 0 {
		return "", false, fmt.Errorf("error building static files: %v", result.Errors)
	}
	if len(result.Warnings) > 0 {
		return "", false, fmt.Errorf("error building static files: %v", result.Warnings)
	}
	return result.Path, result.External, nil
}

func (p *plugin) resolvePathAny(resolveDir string, paths ...string) (string, bool, error) {
	var (
		resolved   string
		isExternal bool
		err        error
	)
	for _, path := range paths {
		if resolved, isExternal, err = p.resolvePath(path, resolveDir); err == nil {
			break
		}
	}
	return resolved, isExternal, err
}

type importsInfo struct {
	byID     map[string]string // local ident  to import path
	byImport map[string]string // import path to local ident
}

func (db *importsInfo) existsByID(id string) bool {
	_, exists := db.byID[id]
	return exists
}

func (db *importsInfo) setByTag(tag, path string) string {
	if existingIdent, ok := db.byImport[path]; ok {
		return existingIdent
	} else {
		id := pascalize(tag)
		db.byID[id] = path
		db.byImport[path] = id
		return id
	}
}

func (db *importsInfo) nextIdent(pascal string) string {
	ident := pascal
	for i := 2; ; i++ {
		if exists := db.existsByID(ident); !exists {
			break
		}
		ident = pascal + strconv.Itoa(i)
	}
	return ident
}

func (db *importsInfo) setByID(id, path string) string {
	if existingIdent, ok := db.byImport[path]; ok {
		return existingIdent
	} else {
		db.byID[id] = path
		db.byImport[path] = id
		return id
	}
}

func copyAttrs(src *importsInfo, targ *Node) {
	keys := make([]string, 0, len(src.byID))
	for k := range src.byID {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, local := range keys {
		targ.Attr = append(targ.Attr, Attribute{Key: local, Val: src.byID[local]})
	}
}

var (
	fileSep     = string(filepath.Separator)
	currentPath = "." + fileSep
)

func (p *plugin) buildImports(idb *importsInfo, root *Node, resolveDir string) (map[string]string, error) {
	rewrites := make(map[string]string) // orig tag to local ident

	for _, tag := range collectElementData(root) {
		if tag == "React.Fragment" {
			continue
		}

		if path, ok := p.cfg.tagMappings[tag]; ok {
			var (
				resolved string
				err      error
			)
			if !filepath.IsAbs(path) {
				if resolved, _, err = p.resolvePath(path, resolveDir); err != nil {
					return nil, err
				}
			} else {
				resolved = path
			}
			rewrites[tag] = idb.setByTag(tag, resolved)
		} else {
			prefix, baseName := tagNameParts(tag)

			// unique local id (ButtonGroup, ButtonGroup2, ...)
			pascal := pascalize(baseName)
			ident := idb.nextIdent(pascal)

			if prefix != "" {
				if basePath, ok := p.cfg.tagPrefixes[prefix]; ok {
					if !strings.HasSuffix(basePath, fileSep) {
						basePath += fileSep
					}
					if resolved, _, err := p.resolvePathAny(
						resolveDir,
						[]string{basePath + pascal, basePath + baseName}...,
					); err != nil {
						return nil, err
					} else {
						rewrites[tag] = idb.setByID(ident, resolved)
					}
				} else {
					prefixedPascal := pascalize(prefix) + pascal
					if resolved, _, err := p.resolvePathAny(
						resolveDir, []string{
							currentPath + sanitizeFileName(tag),
							currentPath + prefixedPascal,
							currentPath + filepath.Join(prefix, sanitizeFileName(baseName)),
							currentPath + filepath.Join(prefix, pascal),
						}...); err != nil {
						return nil, err
					} else {
						ident = idb.nextIdent(prefixedPascal)
						rewrites[tag] = idb.setByID(ident, resolved)
					}
				}
			} else {
				if resolved, _, err := p.resolvePathAny(
					resolveDir, []string{
						currentPath + sanitizeFileName(tag),
						currentPath + pascal,
					}...); err != nil {
					return nil, err
				} else {
					rewrites[tag] = idb.setByID(ident, resolved)
				}
			}
		}
	}
	return rewrites, nil
}

func (p *plugin) rewriteImports(root *Node, resolveDir string) error {
	idb := &importsInfo{
		byID:     make(map[string]string),
		byImport: make(map[string]string),
	}

	rewrites, err := p.buildImports(idb, root, resolveDir)
	if err != nil {
		return err
	}

	copyAttrs(idb, root)

	rewriteElementData(root, rewrites)

	return nil
}

func (p *plugin) onLoad(args api.OnLoadArgs) (api.OnLoadResult, error) {
	root, err := parseFile(args.Path)
	if err != nil {
		return api.OnLoadResult{}, err
	}
	resolveDir := filepath.Dir(args.Path)

	if root.FirstChild == nil {
		contents := "export default function() { return null; }"
		return api.OnLoadResult{
			Contents:   &contents,
			Loader:     api.LoaderJSX,
			ResolveDir: resolveDir,
		}, nil
	}

	componentName := strings.TrimSuffix(filepath.Base(args.Path), filepath.Ext(args.Path))
	root.Data = pascalize(componentName)

	if err := p.rewriteImports(root, resolveDir); err != nil {
		return api.OnLoadResult{}, err
	}

	var buf bytes.Buffer
	if _, err := buf.WriteString("import * as React from 'react';\n"); err != nil {
		return api.OnLoadResult{}, err
	}

	if _, err := Render(&buf, root); err != nil {
		return api.OnLoadResult{}, err
	}
	contents := buf.String()

	return api.OnLoadResult{
		Contents:   &contents,
		Loader:     api.LoaderJSX,
		ResolveDir: resolveDir,
	}, nil
}

func parseFile(path string) (*Node, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	node, err := Parse(f)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func Plugin(opts ...PluginOption) api.Plugin {
	cfg := pluginConfig{}
	readPluginConfig(&cfg, opts...)
	if cfg.extName == "" {
		cfg.extName = ".stache"
	}
	filter := regexp.QuoteMeta(cfg.extName) + "$"
	return api.Plugin{
		Name: "stache-loader",
		Setup: func(pb api.PluginBuild) {
			p := &plugin{cfg: &cfg, buildOptions: pb.InitialOptions, resolveFunc: pb.Resolve}
			pb.OnLoad(api.OnLoadOptions{Filter: filter}, p.onLoad)
		},
	}
}

func pascalize(s string) string {
	var result []rune
	upperNext := true

	for _, c := range s {
		if c == '-' {
			upperNext = true
			continue
		}
		if upperNext {
			if 'a' <= c && c <= 'z' {
				c = c - 'a' + 'A'
			}
			upperNext = false
		}
		result = append(result, c)
	}
	return string(result)
}

func sanitizeFileName(tagName string) string {
	result := make([]byte, len(tagName))
	for i := range len(tagName) {
		c := tagName[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result[i] = c
		} else {
			result[i] = '-'
		}
	}
	return string(result)
}

func sanitizeExtensionName(extName string) string {
	extName = strings.TrimSpace(extName)
	n := len(extName)

	if n < 2 || extName[0] != '.' {
		panic("invalid extension name: must start with a dot and be at least 2 characters")
	}

	valid := false
	for i := 1; i < n; i++ {
		c := extName[i]

		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') {
			valid = true
		} else if c != '-' && c != '.' {
			panic("invalid extension name: only letters, digits, dash and dot allowed after leading dot")
		}
	}

	if !valid {
		panic("invalid extension name: must contain letters or digits after the leading dot")
	}

	return extName
}

func tagNameParts(tag string) (string, string) {
	prefix, baseName := "", tag
	if i := strings.Index(tag, ":"); i != -1 {
		prefix, baseName = tag[:i], tag[i+1:]
	}
	return prefix, baseName
}
