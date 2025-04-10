package restache

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

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
