package stache

import (
	"path/filepath"
	"strings"
)

// Template represents a parsed template with a name and a list of nodes.
type Template struct {
	Name string
	Path string
	Root *Node
	Deps []*Template
}

// NewTemplate creates a new Template with the given name.
func NewTemplate(path string) *Template {
	return &Template{
		Path: path,
		Name: strings.ToLower(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))),
	}
}

type sortNode struct {
	deps []int
	tmpl *Template
}

func (t *sortNode) Afters() []int {
	return t.deps
}
