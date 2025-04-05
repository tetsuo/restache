package stache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

type Component struct {
	Name     string
	FullPath string
	RelPath  string
	Doc      *Node
	Deps     []*Component
}

type depsInfo struct {
	afters    []int
	component *Component
}

func (n depsInfo) Afters() []int {
	return n.afters
}

func parseWithDeps(file string, tbl *fnvtable.Table, info *depsInfo) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", file, err)
	}
	defer f.Close()

	p := newParser(f)
	p.dtbl = tbl
	if err := p.parse(); err != nil {
		return fmt.Errorf("error parsing file %s: %w", file, err)
	}

	*info = depsInfo{afters: p.deps, component: &Component{Doc: p.doc, FullPath: file}}
	return nil
}

// ParseDir reads the provided directory for files listed in includes, parses each file to
// build a dependency graph, and returns a slice of components in topologically sorted order.
func ParseDir(dir string, includes []string, opts ...Option) ([]*Component, error) {
	if len(includes) < 1 {
		return nil, fmt.Errorf("includes must contain at least one file")
	}

	tagIndex := make(map[string]int, len(includes))
	tagNames := make([][]byte, len(includes))

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
		tag, err = sanitizeTagName(tag)
		if err != nil {
			return nil, fmt.Errorf("include entry '%s': %w", file, err)
		}

		if prev, exists := tagIndex[tag]; exists {
			return nil, fmt.Errorf("tag name collision: '%s' from both '%q' and '%s'", tag, includes[prev], file)
		}

		tagIndex[tag] = i
		tagNames[i] = []byte(tag)
	}

	tbl, err := fnvtable.New(tagNames)
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
		cfg.parallelism = min(len(includes), cfg.parallelism)
	}

	depsInfos := make([]depsInfo, len(includes))

	var g errgroup.Group
	g.SetLimit(cfg.parallelism)

	for i, path := range includes {
		i, path := i, filepath.Join(dir, path)
		g.Go(func() error {
			return parseWithDeps(path, tbl, &depsInfos[i])
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	shouldsort := false

	for tag, idx := range tagIndex {
		c := depsInfos[idx].component
		c.Name = tag
		c.RelPath = includes[idx]
		afters := depsInfos[idx].afters
		n := len(afters)
		if n > 0 {
			shouldsort = true
			c.Deps = make([]*Component, n)
			for j, k := range afters {
				c.Deps[j] = depsInfos[k].component
			}
		}
	}

	if shouldsort {
		if err := toposort.BFS(depsInfos); err != nil {
			return nil, fmt.Errorf("error sorting files in %s: %w", dir, err)
		}
	}

	components := make([]*Component, len(depsInfos))
	for i, n := range depsInfos {
		components[i] = n.component
	}

	return components, nil
}

func validTagRune(r rune) rune {
	switch {
	case 'a' <= r && r <= 'z',
		'A' <= r && r <= 'Z',
		'0' <= r && r <= '9',
		r == '-', r == '_':
		return r
	default:
		return -1
	}
}

func sanitizeTagName(name string) (string, error) {
	if name == "" {
		return "", errors.New("tag name cannot be empty")
	}

	sanitized := strings.Map(validTagRune, name)

	if sanitized == "" {
		return "", errors.New("tag name contains only invalid characters")
	}

	first := sanitized[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z')) {
		return "", errors.New("tag name must start with a letter")
	}

	return strings.ToLower(sanitized), nil
}
