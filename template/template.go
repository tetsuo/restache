package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tetsuo/stache/parser"
)

type Template struct {
	name  string
	trees []parser.Node
}

func New(name string) *Template {
	t := &Template{
		name:  name,
		trees: []parser.Node{},
	}
	return t
}

func (t *Template) Name() string {
	return t.name
}

func (t *Template) Serialize() interface{} {
	c := make([]interface{}, len(t.trees))
	for i, tree := range t.trees {
		c[i] = tree.Serialize()
	}
	return []interface{}{t.name, map[string]interface{}{}, c}
}

func (t *Template) Trees() []parser.Node {
	return t.trees
}

func ParseGlob(pattern string) ([]*Template, error) {
	return parseGlob(pattern)
}

func parseGlob(pattern string) ([]*Template, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
	}
	return parseFiles(filenames...)
}

func parseFiles(filenames ...string) ([]*Template, error) {
	if len(filenames) == 0 {
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	var ts []*Template
	for _, filename := range filenames {
		t := New(strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)))
		r, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		parser.Parse(r, func(tree parser.Node) bool {
			t.trees = append(t.trees, tree)
			return true
		})
		ts = append(ts, t)
	}
	return ts, nil
}
