// Package template provides support for retrieving and parsing template files from disk.
package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tetsuo/stache"
)

// Template represents a parsed template with a name and a list of nodes.
type Template struct {
	Name   string
	Source string
	Root   *stache.Node
}

// NewTemplate creates a new Template with the given name.
func NewTemplate(source string) *Template {
	t := &Template{
		Source: source,
		Name:   strings.ToLower(strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))),
	}
	return t
}

// ParseGlob finds files matching the given pattern and parses them into templates.
func ParseGlob(pattern string, maxWorkers int) ([]*Template, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	filesLen := len(files)
	if filesLen == 0 {
		return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
	}
	if maxWorkers == 0 {
		maxWorkers = 4
	}
	sem := make(chan struct{}, min(filesLen, maxWorkers))
	return parseFiles(sem, files...)
}

// parseFiles reads and parses the specified files into templates using a WaitGroup and limited workers.
func parseFiles(sem chan struct{}, files ...string) ([]*Template, error) {
	if len(files) == 0 {
		return []*Template{}, nil
	}

	var ts []*Template
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			t := NewTemplate(f)
			r, err := os.Open(f)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", f, err)
				return
			}
			defer r.Close()

			root, err := stache.Parse(r)
			if err != nil {
				fmt.Printf("Error parsing file %s: %v\n", f, err)
				return
			}

			t.Root = root

			mu.Lock()
			ts = append(ts, t)
			mu.Unlock()
		}(file)
	}

	wg.Wait()
	return ts, nil
}
