package restache

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/tetsuo/toposort"
	"golang.org/x/sync/errgroup"
)

type Option func(*config)

type config struct {
	parallelism int
}

func WithParallelism(limit int) Option {
	return func(cfg *config) {
		cfg.parallelism = limit
	}
}

// Parse parses a single Node with no dependencies.
func Parse(r io.Reader) (node *Node, err error) {
	p := newParser(r, nil)
	if err = p.parse(); err != nil {
		return
	}
	node = p.doc
	return
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

// ParseDir reads the provided directory for files listed in includes, parses each file to
// build a dependency graph, and returns a slice of components in topologically sorted order.
func ParseDir(dir string, includes []string, opts ...Option) ([]*Node, error) {
	n := len(includes)
	if n < 1 {
		return nil, fmt.Errorf("no input files provided")
	}

	var err error
	entries := make([]*componentEntry, n)
	lup := make(map[string]int, n)
	for i, path := range includes {
		entries[i], err = newComponentEntry(dir, path)
		if err != nil {
			return nil, err
		}
		lup[entries[i].tag] = i
	}

	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.parallelism < 1 {
		cfg.parallelism = runtime.NumCPU()
	} else {
		cfg.parallelism = min(n, cfg.parallelism)
	}

	var g errgroup.Group
	g.SetLimit(cfg.parallelism)

	for _, e := range entries {
		e.lookup = lup
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
			return nil, fmt.Errorf("error sorting files in %s: %w", dir, err)
		}
	}

	nodes := make([]*Node, n)
	for i, entry := range entries {
		nodes[i] = entry.doc
	}

	return nodes, nil
}
