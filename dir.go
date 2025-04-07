package restache

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

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

type componentEntry struct {
	path string // full path (e.g. /foo/Bar.stache)
	ext  string // extension (e.g. .stache)
	stem string // filename without extension (e.g. Bar)
	tag  string // lowercase stem (e.g. bar)

	doc    *Node          // the root component node
	lookup map[string]int // dependency lookup table
	afters []int          // collected dependency indexes
}

// componentEntry parses a template with dependencies; it implements toposort.Vertex.
func newComponentEntry(dir, path string) (*componentEntry, error) {
	if filepath.Base(path) != path {
		return nil, fmt.Errorf("input %q must be a filename", path)
	}
	entry := &componentEntry{ext: filepath.Ext(path)}
	if entry.ext != "" {
		entry.stem = path[:len(path)-len(entry.ext)]
	} else {
		entry.stem = path
	}
	if len(entry.stem) == 0 {
		return nil, fmt.Errorf("input filename %q is not valid", path)
	}
	var hasUpper bool
	hasUpper, err := validateTagName(entry.stem)
	if err != nil {
		return nil, fmt.Errorf("input filename %q %v", path, err)
	}

	if hasUpper {
		entry.tag = strings.ToLower(entry.stem)
	} else {
		entry.tag = entry.stem
	}

	entry.path = filepath.Join(dir, path)

	return entry, nil
}

func (e *componentEntry) parse() error {
	f, err := os.Open(e.path)
	if err != nil {
		return fmt.Errorf("could not open file %q: %w", e.path, err)
	}
	defer f.Close()

	p := newParser(f, e.lookup)
	if err := p.parse(); err != nil {
		return fmt.Errorf("failed to parse file %q: %w", e.path, err)
	}

	e.doc = p.doc
	e.doc.Data = e.tag
	e.doc.Path = []PathSegment{{Key: e.stem}, {Key: e.ext}}

	e.afters = slices.Collect(maps.Keys(p.afters))
	e.doc.Attr = make([]Attribute, len(e.afters))

	return nil
}

func (fp *componentEntry) Afters() []int {
	return fp.afters
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

// validateTagName checks if the provided name is valid and returns true if it contains
// at least one upper case letter.
func validateTagName(name string) (bool, error) {
	var hasUpper bool
	c := rune(name[0])
	if 'A' <= c && c <= 'Z' {
		hasUpper = true
	} else if 'a' > c || c > 'z' {
		return false, fmt.Errorf("must start with a letter")
	}
	for _, c = range name[1:] {
		if 'A' <= c && c <= 'Z' {
			hasUpper = true
		} else if !('a' <= c && c <= 'z' || '0' <= c && c <= '9') {
			return false, fmt.Errorf("contains invalid character %q", c)
		}
	}
	return hasUpper, nil
}
