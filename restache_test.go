package restache_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/tetsuo/restache"
)

func TestTranspileFileInvalidInputFile(t *testing.T) {
	for _, tt := range []struct {
		desc      string
		inputFile string
		wantErr   string
	}{
		{
			"empty input path",
			"",
			"input file path is empty",
		},
		{
			"filename with invalid char",
			"/components/bad$name.stache",
			"contains invalid character '$'",
		},
		{
			"empty stem",
			"/components/.stache",
			"is not valid",
		},
		{
			"starts with symbol",
			"/components/_hidden.stache",
			"must start with a letter",
		},
		{
			"input is root slash",
			"/",
			"must start with a letter",
		},
		{
			"input is series of slashes",
			"///",
			"must start with a letter",
		},
		{
			"input is '.' (current directory)",
			".",
			": is a directory", // handled during parse()
		},
		{
			"input is '..' (parent directory)",
			"..",
			": is a directory", // handled during parse()
		},
		{
			"input doesn't exist",
			"/i/am/not/exist.stache",
			"no such file or directory",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := restache.TranspileFile(tt.inputFile, "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got err = %v, want err containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestParseModuleDeps(t *testing.T) {
	tests := []struct {
		name      string
		recursive bool
	}{
		{name: "card_grid"},
		{name: "deep_tree"},
		{name: "diamond"},
		{name: "fan_out_and_in"},
		{name: "fruits_basic"},
		{name: "many_to_one"},
		{name: "mixed_recursive", recursive: true},
		{name: "nested_panels"},
		{name: "one_to_many"},
		{name: "recursive_tree", recursive: true},
		{name: "standalone"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			baseDir := filepath.Join("testdata", "module_deps", tc.name)
			_, relPaths := mustGlobFiles(t, baseDir, "*.stache")
			expectedAttrs := mustParseDepsFile(t, baseDir)

			recursiveMap := make(map[string]bool)
			for from, attrs := range expectedAttrs {
				filtered := attrs[:0]
				for _, attr := range attrs {
					if attr.Key == from {
						recursiveMap[from] = true
						continue
					}
					filtered = append(filtered, attr)
				}
				expectedAttrs[from] = filtered
			}

			nodes, err := restache.ParseDir(baseDir, restache.WithIncludes(relPaths))
			if err != nil {
				t.Fatalf("ParseModule failed: %v", err)
			}
			for _, node := range nodes {
				if node.Type != restache.ComponentNode {
					t.Errorf("expected node %q to be ComponentNode, got %v", node.Data, node.Type)
				}
				if len(node.Path) != 2 {
					t.Errorf("expected node %q Path len=2, got %d", node.Data, len(node.Path))
				} else {
					if node.Path[0].Key != node.Data {
						t.Errorf("expected Path[0] == Data (%q), got %q", node.Data, node.Path[0].Key)
					}
					if node.Path[1].Key != ".stache" {
						t.Errorf("expected Path[1] == .stache, got %q", node.Path[1].Key)
					}
				}
				isRecursive := recursiveMap[node.Data]
				if isRecursive {
					if node.DataAtom != 1 {
						t.Errorf("expected recursive node %q to have DataAtom=1, got %v", node.Data, node.DataAtom)
					}
					for _, attr := range node.Attr {
						if attr.Key == node.Data {
							t.Errorf("recursive node %q should not have itself in Attrs", node.Data)
						}
					}
				} else {
					if node.DataAtom != 0 {
						t.Errorf("expected non-recursive node %q to have DataAtom=0, got %v", node.Data, node.DataAtom)
					}
				}
				expected := expectedAttrs[node.Data]
				if len(node.Attr) != len(expected) {
					t.Errorf("node %q: expected %d Attrs, got %d", node.Data, len(expected), len(node.Attr))
					continue
				}
				sort.Slice(expected, func(i, j int) bool {
					return expected[i].Key < expected[j].Key
				})
				sort.Slice(node.Attr, func(i, j int) bool {
					return node.Attr[i].Key < node.Attr[j].Key
				})
				for i := range node.Attr {
					if node.Attr[i] != expected[i] {
						t.Errorf("node %q Attr[%d] mismatch: expected %+v, got %+v", node.Data, i, expected[i], node.Attr[i])
					}
				}
			}
		})
	}
}

func mustGlobFiles(t *testing.T, dir, pattern string) ([]string, []string) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		t.Fatalf("failed to glob stache files: %v", err)
	}
	if len(matches) == 0 {
		t.Fatalf("no .stache files found in %s", dir)
	}
	// Build relative paths
	var relPaths []string
	for _, m := range matches {
		rel, err := filepath.Rel(dir, m)
		if err != nil {
			t.Fatalf("failed to get relative path: %v", err)
		}
		relPaths = append(relPaths, rel)
	}
	return matches, relPaths
}

func mustParseDepsFile(t *testing.T, dir string) map[string][]restache.Attribute {
	depsPath := filepath.Join(dir, "deps.txt")
	expectedAttrs := make(map[string][]restache.Attribute) // key = component name
	if content, err := os.ReadFile(depsPath); err == nil {
		lines := strings.SplitSeq(string(content), "\n")
		for line := range lines {
			fields := strings.Fields(line)
			if len(fields) == 2 {
				from, to := fields[0], fields[1]
				expectedAttrs[from] = append(expectedAttrs[from], restache.Attribute{
					Key:    to,
					Val:    to,
					IsExpr: false,
				})
			}
		}
	} else if !os.IsNotExist(err) {
		t.Fatalf("failed to read deps.txt: %v", err)
	}
	return expectedAttrs
}
