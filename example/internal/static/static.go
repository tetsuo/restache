// Package static builds static assets for the frontend.
package static

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tetsuo/restache"
)

func Build(config Config) error {
	files, err := getEntry(config.EntryPoint, config.Bundle)
	if err != nil {
		return err
	}
	options := api.BuildOptions{
		EntryPoints: files,
		Bundle:      config.Bundle,
		Outdir:      config.EntryPoint,
		Write:       true,
		Platform:    api.PlatformBrowser,
		Plugins:     []api.Plugin{restache.Plugin()},
		Format:      api.FormatESModule,
		OutExtension: map[string]string{
			".css": ".min.css",
		},
		External: []string{"*.svg", "react", "react-dom/client"},
	}
	if config.Minify {
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
		options.MinifyWhitespace = true
	}
	if config.Watch {
		options.Sourcemap = api.SourceMapLinked
		ctx, err := api.Context(options)
		if err != nil {
			return err
		}
		return ctx.Watch(api.WatchOptions{})
	}
	options.Sourcemap = api.SourceMapNone
	result := api.Build(options)
	if len(result.Errors) > 0 {
		return fmt.Errorf("error building static files: %v", result.Errors)
	}
	if len(result.Warnings) > 0 {
		return fmt.Errorf("error building static files: %v", result.Warnings)
	}
	return nil
}

// getEntry walks the given directory and collects entry file paths
// for esbuild. It ignores test files and files prefixed with an underscore.
// Underscore prefixed files are assumed to be imported by and bundled together
// with the output of an entry file.
func getEntry(dir string, bundle bool) ([]string, error) {
	var matches []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		basePath := filepath.Base(path)
		notPartial := !strings.HasPrefix(basePath, "_")
		notTest := !strings.HasSuffix(basePath, ".test.mjs")
		isMJS := strings.HasSuffix(basePath, ".mjs")
		isCSS := strings.HasSuffix(basePath, ".css") && !strings.HasSuffix(basePath, ".min.css")
		if notPartial && notTest && (isMJS || (bundle && isCSS)) {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
