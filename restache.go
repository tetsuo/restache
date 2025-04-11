package restache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tetsuo/toposort"
	"golang.org/x/sync/errgroup"
)

type Option func(*config)

type config struct {
	includes    []string
	parallelism int
}

func WithParallelism(parallelism int) Option {
	return func(cfg *config) {
		cfg.parallelism = parallelism
	}
}

func WithIncludes(includes []string) Option {
	return func(cfg *config) {
		cfg.includes = includes
	}
}

func ParseFile(path string) (node *Node, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		err = fmt.Errorf("error opening file %s: %w", path, err)
		return
	}
	defer f.Close()
	node, err = Parse(f)
	if err != nil {
		err = fmt.Errorf("error parsing file %s: %w", path, err)
	}
	return
}

func readConfig(cfg *config, inputDir string, opts ...Option) (int, error) {
	for _, opt := range opts {
		opt(cfg)
	}

	includes := cfg.includes

	var err error
	n := len(includes)
	if n == 0 {
		pat := filepath.Join(inputDir, "*")
		includes, err = filepath.Glob(pat)
		if err != nil {
			return 0, fmt.Errorf("invalid glob %q: %v", pat, err)
		}
		n = len(includes)
		if n == 0 {
			return 0, fmt.Errorf("no input files found in directory %q", inputDir)
		}
		for i, f := range includes {
			includes[i] = filepath.Base(f)
		}
	}

	return n, nil
}

func parseModule(inputDir string, includes []string, parallelism int) ([]*Node, error) {
	var err error
	n := len(includes)
	entries := make([]*componentEntry, n)
	lookupTable := make(map[string]int, n)
	for i, path := range includes {
		entries[i], err = newComponentEntry(inputDir, path)
		if err != nil {
			return nil, err
		}
		lookupTable[entries[i].tag] = i
	}

	var g errgroup.Group
	g.SetLimit(parallelism)

	for _, e := range entries {
		e.lookup = lookupTable
		g.Go(e.parse)
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	sort := false

	// Before sorting collect imports:
	for _, e := range entries {
		if e.doc.Attr != nil {
			for i, other := range e.afters {
				e.doc.Attr[i] = Attribute{Key: entries[other].tag, Val: entries[other].stem}
			}
			sort = true
		}
	}

	if sort {
		if err := toposort.BFS(entries); err != nil {
			// TODO: will detect recursive entries earlier; this only returns ErrCircular,
			// should rather panic after recursive detection.
			return nil, fmt.Errorf("error sorting files in %s: %w", inputDir, err)
		}
	}

	nodes := make([]*Node, n)
	for i, entry := range entries {
		nodes[i] = entry.doc
	}

	return nodes, nil
}

// ParseModule reads the provided directory for files listed in includes, parses each file to
// build a dependency graph, and returns a slice of components in topologically sorted order.
func ParseModule(inputDir string, opts ...Option) ([]*Node, error) {
	if inputDir == "" {
		return nil, fmt.Errorf("input directory path is empty")
	}

	absInputDir, err := toAbsPath(inputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve input directory %q: %v", inputDir, err)
	}

	cfg := config{}
	n, err := readConfig(&cfg, absInputDir, opts...)
	if err != nil {
		return nil, err
	}

	return parseModule(inputDir, cfg.includes, min(n, cfg.parallelism))
}

func TranspileFile(inputFile, outputFile string) error {
	if inputFile == "" {
		return fmt.Errorf("input file path is empty")
	}

	absInputFile, err := toAbsPath(inputFile)
	if err != nil {
		return fmt.Errorf("failed to resolve input file %q: %v", inputFile, err)
	}

	absInputDir := filepath.Dir(absInputFile)

	e, err := newComponentEntry(absInputDir, filepath.Base(absInputFile))
	if err != nil {
		return err // returns descriptive error
	}
	e.lookup = map[string]int{e.tag: 0}

	if err := e.parse(); err != nil {
		return err // returns descriptive error
	}

	var absOutputFile string
	if outputFile == "" {
		absOutputFile = filepath.Join(absInputDir, e.stem+".jsx")
	} else {
		absOutputFile, err = toAbsPath(outputFile)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path for output file %q: %v", outputFile, err)
		}
	}

	return renderToFile(absOutputFile, e.doc)
}

func renderToFile(absPath string, n *Node) error {
	f, err := os.Create(absPath)
	if err != nil {
		return fmt.Errorf("could not create output file %q: %v", absPath, err)
	}
	defer f.Close()

	if err := Render(f, n); err != nil {
		return fmt.Errorf("failed to write output file %q: %v", absPath, err)
	}

	return nil
}

func TranspileModule(inputDir string, outputDir string, opts ...Option) error {
	if inputDir == "" {
		return fmt.Errorf("input directory path is empty")
	}

	absInputDir, err := toAbsPath(inputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve input directory %q: %v", inputDir, err)
	}

	cfg := config{}
	n, err := readConfig(&cfg, absInputDir, opts...)
	if err != nil {
		return err
	}

	parallelism := min(n, cfg.parallelism)

	nodes, err := parseModule(inputDir, cfg.includes, parallelism)
	if err != nil {
		return fmt.Errorf("parse module %q: %v", inputDir, err)
	}

	var absOutputDir string
	if outputDir == "" {
		absOutputDir = absInputDir
	} else {
		absOutputDir, err = toAbsPath(outputDir)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path for output directory %q: %v", outputDir, err)
		}
	}

	var g errgroup.Group
	g.SetLimit(min(n, cfg.parallelism))

	for _, node := range nodes {
		node := node
		g.Go(func() error {
			// parseModule guarantees that node.Path is always at least length 2 (stem, ext),
			// otherwise this might panic on certain input files:
			outfile := filepath.Join(absOutputDir, node.Path[:len(node.Path)-1][0].Key) + ".jsx"
			dst, err := os.Create(outfile)
			if err != nil {
				return fmt.Errorf("could not create file %q: %v", outfile, err)
			}
			defer dst.Close()
			if err = Render(dst, node); err != nil {
				return fmt.Errorf("failed to write file %q: %v", outfile, err)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func toAbsPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		} else {
			return absPath, nil
		}
	}
	return path, nil
}
