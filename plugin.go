package restache

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"golang.org/x/net/html/atom"
)

func Plugin() api.Plugin {
	return api.Plugin{
		Name:  "stache-loader",
		Setup: pluginSetup,
	}
}

func pluginSetup(build api.PluginBuild) {
	build.OnResolve(api.OnResolveOptions{Filter: "\\.stache$"}, pluginBuildResolveHandler)
	build.OnLoad(api.OnLoadOptions{Filter: "\\.stache$"}, pluginBuildLoadHandler)
}

func pluginBuildResolveHandler(args api.OnResolveArgs) (api.OnResolveResult, error) {
	path := filepath.Join(args.ResolveDir, args.Path)
	return api.OnResolveResult{
		Path:      path,
		Namespace: "file",
	}, nil
}

func pluginBuildLoadHandler(args api.OnLoadArgs) (api.OnLoadResult, error) {
	e, err := newFileParser(args.Path)
	if err != nil {
		return api.OnLoadResult{}, err
	}

	if err := e.parse(); err != nil {
		return api.OnLoadResult{}, err
	}

	var buf bytes.Buffer

	if _, err := Render(&buf, e.doc); err != nil {
		return api.OnLoadResult{}, err
	}

	contents := buf.String()
	return api.OnLoadResult{
		Contents: &contents,
		Loader:   api.LoaderJSX,
	}, nil
}

// fileParser parses a template with dependencies; it implements toposort.Vertex.
type fileParser struct {
	path string // full path (e.g. /foo/Bar.stache)
	ext  string // extension (e.g. .stache)
	stem string // filename without extension (e.g. Bar)
	tag  string // lowercase stem (e.g. bar)
	doc  *Node  // the root component node
}

func newFileParser(absPath string) (*fileParser, error) {
	path := filepath.Base(absPath)
	z := &fileParser{ext: filepath.Ext(path), path: absPath}
	if z.ext != "" {
		z.stem = path[:len(path)-len(z.ext)]
	} else {
		z.stem = path
	}
	if len(z.stem) == 0 {
		// This also catches when path = "."
		return nil, fmt.Errorf("input filename %q is not valid", path)
	}
	a := atom.Lookup([]byte(z.stem))
	if a != 0 {
		if _, ok := commonElements[a]; ok {
			return nil, fmt.Errorf("invalid filename %q: <%s> conflicts with standard element", path, z.stem)
		}
	}
	var hasUpper bool
	hasUpper, err := isTagNameValidAndUpper(z.stem)
	if err != nil {
		// This also catches ".."
		return nil, fmt.Errorf("input filename %q %v", path, err)
	}
	if hasUpper {
		z.tag = strings.ToLower(z.stem)
	} else {
		z.tag = z.stem
	}
	return z, nil
}

func (s *fileParser) parse() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	p := newParser(f)
	if err := p.parse(); err != nil {
		return err
	}

	s.doc = p.doc
	s.doc.Data = s.tag
	s.doc.Path = []PathComponent{{Key: s.stem}, {Key: s.ext}}

	return nil
}

// isTagNameValidAndUpper checks if the provided name is valid as HTML tag and
// returns true if it contains at least one upper case letter.
func isTagNameValidAndUpper(name string) (bool, error) {
	if len(name) == 0 {
		return false, fmt.Errorf("name is empty")
	}

	first := name[0]
	if !isLetter(first) {
		return false, fmt.Errorf("must start with an ASCII letter")
	}

	hasUpper := isUpper(first)
	for i := 1; i < len(name); i++ {
		c := name[i]
		if isUpper(c) {
			hasUpper = true
		} else if !isLetter(c) && !isDigit(c) && c != '-' {
			return false, fmt.Errorf("contains invalid character %q", c)
		}
	}

	return hasUpper, nil
}

func isLetter(c byte) bool {
	return isUpper(c) || isLower(c)
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return (c >= 'a' && c <= 'z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
