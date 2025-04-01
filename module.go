package stache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"github.com/tetsuo/toposort"
)

// Template represents a parsed template with a name and a list of nodes.
type Template struct {
	Name string
	Path string
	Root *Node

	dependsOn []int
}

// NewTemplate creates a new Template with the given name.
func NewTemplate(path string) *Template {
	t := &Template{
		Path: path,
		Name: strings.ToLower(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))),
	}
	return t
}

func (t *Template) Afters() []int {
	return t.dependsOn
}

type Module struct {
	Name      string
	Path      string
	Templates []*Template
}

func findFiles(wd string, fileMap map[string][]string, patterns []string) (err error) {
	var (
		m    []string
		p    string
		info os.FileInfo
		dir  string
	)

	for _, p = range patterns {
		m, err = filepath.Glob(filepath.Join(wd, p))
		if err != nil {
			return
		}
		for _, p = range m {
			info, err = os.Stat(p)
			if err != nil {
				return
			}
			if info.IsDir() {
				continue
			}
			dir = filepath.Dir(p)
			fileMap[dir] = append(fileMap[dir], p)
		}
	}

	return nil
}

func ParseGlob(wd string, maxParallelism int, patterns ...string) (modules []*Module, err error) {
	fileMap := make(map[string][]string)

	if err = findFiles(wd, fileMap, patterns); err != nil {
		return
	}

	var (
		sem       chan struct{}
		templates []*Template
	)
	for dir, files := range fileMap {
		sem = make(chan struct{}, min(len(files), maxParallelism))
		templates, err = parseFiles(sem, files...)
		if err != nil {
			return
		}
		if len(templates) < 1 {
			continue
		}
		if err = toposort.BFS(templates); err != nil {
			return
		}
		modules = append(modules, &Module{
			Name:      filepath.Base(dir),
			Path:      dir,
			Templates: templates,
		})
	}
	return modules, nil
}

func parseFiles(sem chan struct{}, paths ...string) ([]*Template, error) {
	n := len(paths)
	if n == 0 {
		return []*Template{}, nil
	}

	names := make([]string, n)
	trie := iradix.New[int]()

	for i, p := range paths {
		p = filepath.Base(p)
		names[i] = strings.TrimSuffix(p, filepath.Ext(p))
		trie, _, _ = trie.Insert([]byte(names[i]), i)
	}

	ts := make([]*Template, n)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, path := range paths {
		wg.Add(1)
		go func(idx int, f string) {
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

			root, dependsOn, err := parseWithDependencies(r, trie)
			if err != nil {
				fmt.Printf("Error parsing file %s: %v\n", f, err)
				return
			}

			t.Root = root
			t.Name = filepath.Base(path)
			t.Name = names[i]
			t.dependsOn = dependsOn

			mu.Lock()
			ts[idx] = t
			mu.Unlock()
		}(i, path)
	}

	wg.Wait()
	return ts, nil
}
