package main

type Config struct {
	// Minify enables JavaScript minification in production builds.
	Minify bool
	// SourceMap enables source map output in development builds.
	SourceMap bool
}
