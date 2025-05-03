package restache

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
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
	root, err := parseFile(args.Path)
	if err != nil {
		return api.OnLoadResult{}, err
	}
	var buf bytes.Buffer
	if _, err := Render(&buf, root); err != nil {
		return api.OnLoadResult{}, err
	}
	contents := buf.String()
	return api.OnLoadResult{
		Contents: &contents,
		Loader:   api.LoaderJSX,
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
	node.Data = "root"
	return node, nil
}
