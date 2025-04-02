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

type Module struct {
	Name      string
	Path      string
	Templates []*Template
}

func ParseGlob(wd string, maxParallelism int, patterns ...string) (modules []*Module, err error) {
	fileMap := make(map[string][]string)
	if err = findFiles(wd, fileMap, patterns); err != nil {
		return
	}
	var (
		sem chan struct{}
		vz  []*sortNode
	)
	for dir, files := range fileMap {
		sem = make(chan struct{}, min(len(files), maxParallelism))
		vz, err = parseFiles(sem, files...)
		if err != nil {
			return
		}
		n := len(vz)
		if n < 1 {
			continue
		}
		templates := make([]*Template, n)
		for i, v := range vz {
			templates[i] = v.tmpl
			templates[i].Deps = make([]*Template, len(v.deps))
			for j, dep := range v.deps {
				templates[i].Deps[j] = vz[dep].tmpl
			}
		}
		if err = toposort.BFS(vz); err != nil {
			return
		}
		for i, v := range vz {
			templates[i] = v.tmpl
		}
		modules = append(modules, &Module{
			Name:      filepath.Base(dir),
			Path:      dir,
			Templates: templates,
		})
	}
	return modules, nil
}

func parseFiles(sem chan struct{}, paths ...string) (vz []*sortNode, err error) {
	n := len(paths)
	if n == 0 {
		return
	}

	names := make([]string, n)
	trie := iradix.New[int]()

	for i, p := range paths {
		p = filepath.Base(p)
		names[i] = strings.ToLower(strings.TrimSuffix(p, filepath.Ext(p)))
		trie, _, _ = trie.Insert([]byte(names[i]), i)
	}

	vz = make([]*sortNode, n)

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

			root, deps, err := parseWithDependencies(r, trie)
			if err != nil {
				fmt.Printf("Error parsing file %s: %v\n", f, err)
				return
			}

			t.Root = root
			t.Name = names[i]

			mu.Lock()
			vz[idx] = &sortNode{deps: deps, tmpl: t}
			mu.Unlock()
		}(i, path)
	}

	wg.Wait()
	return
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
