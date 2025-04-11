package restache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
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
		// This also catches when path = "."
		return nil, fmt.Errorf("input filename %q is not valid", path)
	}
	var hasUpper bool
	hasUpper, err := validateTagName(entry.stem)
	if err != nil {
		// This also catches ".."
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

	e.afters = collectKeys(p.afters)
	e.doc.Attr = make([]Attribute, len(e.afters))

	return nil
}

func (fp *componentEntry) Afters() []int {
	return fp.afters
}

// validateTagName checks if the provided name is valid as HTML tag and
// returns true if it contains at least one upper case letter.
func validateTagName(name string) (bool, error) {
	r := rune(name[0])
	if !unicode.IsLetter(r) {
		return false, fmt.Errorf("must start with a letter")
	}

	hasUpper := unicode.IsUpper(r)

	for _, r := range name[1:] {
		if unicode.IsUpper(r) {
			hasUpper = true
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return false, fmt.Errorf("contains invalid character %q", r)
		}
	}

	return hasUpper, nil
}

func collectKeys[K comparable, V any](m map[K]V) (keys []K) {
	for key := range m {
		keys = append(keys, key)
	}
	return
}
