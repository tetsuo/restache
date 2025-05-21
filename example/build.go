package main

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tetsuo/restache"
)

func Build(cfg Config) error {
	options := api.BuildOptions{
		EntryPoints:       []string{"static/index.mjs"},
		Bundle:            cfg.Bundle,
		Outdir:            "static",
		Write:             true,
		Platform:          api.PlatformBrowser,
		ResolveExtensions: []string{".stache", ".jsx", ".js"},
		Plugins:           []api.Plugin{restache.Plugin()},
		Format:            api.FormatESModule,
		OutExtension: map[string]string{
			".css": ".min.css",
		},
		JSX:      api.JSXTransform,
		External: []string{"*.svg", "react", "react-dom/client", "@mui/material"},
	}
	if cfg.Minify {
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
		options.MinifyWhitespace = true
	}
	if cfg.Watch {
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
