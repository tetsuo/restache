package restache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// fileParser parses a template with dependencies; it implements toposort.Vertex.
type fileParser struct {
	path   string         // full path (e.g. /foo/Bar.stache)
	ext    string         // extension (e.g. .stache)
	stem   string         // filename without extension (e.g. Bar)
	tag    string         // lowercase stem (e.g. bar)
	doc    *Node          // the root component node
	lookup map[string]int // dependency lookup table
	afters []int          // collected dependency indexes
}

func newFileParser(absPath string) (*fileParser, error) {
	path := filepath.Base(absPath)
	z := &fileParser{ext: filepath.Ext(path), path: absPath}
	if z.ext != "" {
		z.stem = path[:len(path)-len(z.ext)]
	} else {
		z.stem = path
	}
	if len(z.stem) == 0 {
		// This also catches when path = "."
		return nil, fmt.Errorf("input filename %q is not valid", path)
	}
	var hasUpper bool
	hasUpper, err := isTagNameValidAndUpper(z.stem)
	if err != nil {
		// This also catches ".."
		return nil, fmt.Errorf("input filename %q %v", path, err)
	}
	if hasUpper {
		z.tag = strings.ToLower(z.stem)
	} else {
		z.tag = z.stem
	}
	return z, nil
}

func (s *fileParser) parse() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	p := newParser(f, s.lookup)
	if err := p.parse(); err != nil {
		return err
	}

	s.doc = p.doc
	s.doc.Data = s.tag
	s.doc.Path = []PathComponent{{Key: s.stem}, {Key: s.ext}}

	s.afters = collectKeys(p.afters)
	s.doc.Attr = make([]Attribute, len(s.afters))

	return nil
}

func (s *fileParser) Afters() []int {
	return s.afters
}

// isTagNameValidAndUpper checks if the provided name is valid as HTML tag and
// returns true if it contains at least one upper case letter.
func isTagNameValidAndUpper(name string) (bool, error) {
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
