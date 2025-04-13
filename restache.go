package restache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tetsuo/toposort"
	"golang.org/x/sync/errgroup"
)

type Option func(*config)

type config struct {
	includes    []string
	parallelism int
	onEmit      func(Artifact)
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

func WithCallback(onEmit func(Artifact)) Option {
	return func(cfg *config) {
		cfg.onEmit = onEmit
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

	if cfg.parallelism == 0 {
		cfg.parallelism = runtime.NumCPU()
	}

	return n, nil
}

func parseModule(inputDir string, includes []string, parallelism int) ([]*Node, error) {
	var err error
	n := len(includes)
	entries := make([]*fileParser, n)
	lookupTable := make(map[string]int, n)
	for i, path := range includes {
		entries[i], err = newFileParser(filepath.Join(inputDir, path))
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
			j := 0
			for i, other := range e.afters {
				// Recursive?
				if entries[other].tag == e.tag {
					e.doc.DataAtom = 1
					continue
				}
				e.afters[j] = e.afters[i]
				e.doc.Attr[j] = Attribute{
					Key: entries[other].tag,
					Val: entries[other].stem,
				}
				j++
			}
			if e.doc.DataAtom == 1 {
				e.afters = e.afters[:j]
				e.doc.Attr = e.doc.Attr[:j]
			}
			if j > 0 {
				sort = true
			}
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

type Artifact struct {
	Source  string
	Path    string
	Bytes   int
	Elapsed time.Duration
}

func TranspileFile(inputFile, outputFile string) (Artifact, error) {
	if inputFile == "" {
		return Artifact{}, fmt.Errorf("input file path is empty")
	}

	absInputFile, err := toAbsPath(inputFile)
	if err != nil {
		return Artifact{}, fmt.Errorf("failed to resolve input file %q: %v", inputFile, err)
	}

	absInputDir := filepath.Dir(absInputFile)

	start := time.Now()

	e, err := newFileParser(absInputFile)
	if err != nil {
		return Artifact{}, err // returns descriptive error
	}
	e.lookup = map[string]int{e.tag: 0}

	if err := e.parse(); err != nil {
		return Artifact{}, err // returns descriptive error
	}

	var absOutputFile string
	if outputFile == "" {
		absOutputFile = filepath.Join(absInputDir, e.stem+".jsx")
	} else {
		absOutputFile, err = toAbsPath(outputFile)
		if err != nil {
			return Artifact{}, fmt.Errorf("failed to resolve absolute path for output file %q: %v", outputFile, err)
		}
	}

	if written, err := renderToFile(absOutputFile, e.doc); err != nil {
		return Artifact{}, err
	} else {
		return Artifact{
			Path:    absOutputFile,
			Bytes:   written,
			Source:  absInputFile,
			Elapsed: time.Since(start),
		}, nil
	}
}

func renderToFile(absPath string, n *Node) (int, error) {
	f, err := os.Create(absPath)
	if err != nil {
		return 0, fmt.Errorf("could not create output file %q: %v", absPath, err)
	}
	defer f.Close()

	if written, err := Render(f, n); err != nil {
		return 0, fmt.Errorf("failed to write output file %q: %v", absPath, err)
	} else {
		return written, nil
	}
}

func TranspileModule(inputDir string, outputDir string, opts ...Option) ([]Artifact, error) {
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

	parallelism := min(n, cfg.parallelism)

	var absOutputDir string
	if outputDir == "" {
		absOutputDir = absInputDir
	} else {
		absOutputDir, err = toAbsPath(outputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve absolute path for output directory %q: %v", outputDir, err)
		}
	}

	start := time.Now()

	nodes, err := parseModule(inputDir, cfg.includes, parallelism)
	if err != nil {
		return nil, fmt.Errorf("parse module %q: %v", inputDir, err)
	}

	var mu sync.Mutex
	artifacts := make([]Artifact, n)

	var g errgroup.Group
	g.SetLimit(min(n, cfg.parallelism))

	for i, node := range nodes {
		node := node
		i := i
		g.Go(func() error {
			// parseModule guarantees that node.Path is always at least length 2 (stem, ext),
			// otherwise this might panic on certain input files:
			outfile := filepath.Join(absOutputDir, node.Path[:len(node.Path)-1][0].Key) + ".jsx"
			dst, err := os.Create(outfile)
			if err != nil {
				return fmt.Errorf("could not create file %q: %v", outfile, err)
			}
			if written, err := Render(dst, node); err != nil {
				dst.Close()
				return fmt.Errorf("failed to write file %q: %v", outfile, err)
			} else {
				dst.Close()
				art := Artifact{Path: outfile, Bytes: written, Elapsed: time.Since(start)}
				mu.Lock()
				artifacts[i] = art
				mu.Unlock()
				if cfg.onEmit != nil {
					cfg.onEmit(art)
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return artifacts, nil
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
