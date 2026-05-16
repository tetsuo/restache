package main

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tetsuo/restache"
)

func Build(cfg Config) error {
	options := api.BuildOptions{
		EntryPoints:       []string{"index.html"},
		Bundle:            true,
		Write:             true,
		Platform:          api.PlatformBrowser,
		ResolveExtensions: []string{".jsx", ".tsx", ".html", ".js", ".mjs", ".ts"},
		Plugins: []api.Plugin{restache.Plugin(
			restache.WithExtensionName(".html"),
			restache.WithTagPrefixes(map[string]string{
				"mui":   "@mui/material",
				"icons": "@mui/icons-material",
			}),
		)},
		Format: api.FormatESModule,
		OutExtension: map[string]string{
			".css": ".min.css",
		},
		Outfile: "static/app.js",
		JSX:     api.JSXTransform,
		External: []string{
			"react",
			"react-dom/client",
			"@mui/material",
			"@mui/icons-material",
			"history",
			"date-fns",
			"date-fns-tz",
		},
	}
	if cfg.Minify {
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
		options.MinifyWhitespace = true
	}
	if cfg.SourceMap {
		options.Sourcemap = api.SourceMapLinked
	} else {
		options.Sourcemap = api.SourceMapNone
	}

	result := api.Build(options)
	if len(result.Errors) > 0 {
		return fmt.Errorf("error building static files: %v", result.Errors)
	}
	if len(result.Warnings) > 0 {
		return fmt.Errorf("error building static files: %v", result.Warnings)
	}
	return nil
}
