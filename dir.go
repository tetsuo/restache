package stache

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/tetsuo/fnvtable"
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

type fileParser struct {
	doc    *Node
	file   string
	lookup lookupFunc
	afters []int
}

func (fp *fileParser) parse() error {
	f, err := os.Open(fp.file)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", fp.file, err)
	}
	defer f.Close()

	p := newParser(f, fp.lookup)
	if err := p.parse(); err != nil {
		return fmt.Errorf("error parsing file %s: %w", fp.file, err)
	}

	fp.afters = slices.Collect(maps.Keys(p.afters))
	fp.doc = p.doc
	return nil
}

func (fp *fileParser) Afters() []int {
	return fp.afters
}

// ParseDir reads the provided directory for files listed in includes, parses each file to
// build a dependency graph, and returns a slice of components in topologically sorted order.
func ParseDir(dir string, includes []string, opts ...Option) ([]*Node, error) {
	n := len(includes)
	if n < 1 {
		return nil, fmt.Errorf("includes must contain at least one file")
	}

	stems := make([][]byte, n) // filenames without extensions
	tags := make([][]byte, n)  // tags are lowercase stems

	for i, file := range includes {
		if filepath.Base(file) != file {
			return nil, fmt.Errorf("include entry '%s' is not a basename", file)
		}
		var (
			tag = file
			err error
		)
		ext := filepath.Ext(file)
		if ext != "" {
			tag = file[:len(file)-len(ext)]
		}
		stems[i] = []byte(tag)
		var hasUpper bool
		hasUpper, err = validateTagName(stems[i])
		if err != nil {
			return nil, fmt.Errorf("include entry '%s' %v", file, err)
		}
		if hasUpper {
			tags[i] = lower(bytes.Clone(stems[i]))
		} else {
			tags[i] = stems[i]
		}
	}

	tbl, err := fnvtable.New(tags)
	if err != nil {
		return nil, err
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

	parsers := make([]*fileParser, n)

	var g errgroup.Group
	g.SetLimit(cfg.parallelism)

	for i, path := range includes {
		i, path := i, filepath.Join(dir, path)
		parsers[i] = &fileParser{file: path, lookup: tbl.Lookup}
		g.Go(parsers[i].parse)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	shouldSort := false

	for i, tagName := range tags {
		doc := parsers[i].doc
		doc.Data = tagName
		afters := parsers[i].afters
		n := len(afters)
		if n > 0 {
			shouldSort = true
			doc.Attr = make([]Attribute, n)
			for j, k := range afters {
				doc.Attr[j] = Attribute{Key: tags[k], Val: stems[k]}
			}
		}
	}

	if shouldSort {
		if err := toposort.BFS(parsers); err != nil {
			return nil, fmt.Errorf("error sorting files in %s: %w", dir, err)
		}
	}

	nodes := make([]*Node, n)
	for i, p := range parsers {
		nodes[i] = p.doc
	}

	return nodes, nil
}

// validateTagName checks if the provided name is valid and returns true if it contains
// at least one upper case letter.
func validateTagName(name []byte) (bool, error) {
	if len(name) == 0 {
		return false, fmt.Errorf("must have a length")
	}
	var hasUpper bool
	c := name[0]
	if 'A' <= c && c <= 'Z' {
		hasUpper = true
	} else if 'a' > c || c > 'z' {
		return false, fmt.Errorf("must start with a letter")
	}
	for _, c = range name[1:] {
		if 'A' <= c && c <= 'Z' {
			hasUpper = true
		} else if !('a' <= c && c <= 'z' || '0' <= c && c <= '9') {
			return false, fmt.Errorf("contains invalid character '%c'", c)
		}
	}
	return hasUpper, nil
}

func lower(b []byte) []byte {
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return b
}
