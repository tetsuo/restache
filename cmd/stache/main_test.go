package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestResolveGlobs(t *testing.T) {
	tmp := t.TempDir()

	subDir := filepath.Join(tmp, "dir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	files := []string{
		filepath.Join(tmp, "a.stache"),
		filepath.Join(tmp, "b.log"),
		filepath.Join(subDir, "c.stache"),
	}
	for _, f := range files {
		if err := os.WriteFile(f, []byte("data"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		desc     string
		patterns []string
		expected map[string][]string
	}{
		{
			desc:     "match top-level stache files",
			patterns: []string{"*.stache"},
			expected: map[string][]string{
				tmp: {filepath.Join(tmp, "a.stache")},
			},
		},
		{
			desc:     "filepath.Glob with **/*.stache still matches one level deep",
			patterns: []string{"**/*.stache"},
			expected: map[string][]string{
				filepath.Join(tmp, "dir"): {
					filepath.Join(tmp, "dir", "c.stache"),
				},
			},
		},
		{
			desc:     "match all .log",
			patterns: []string{"*.log"},
			expected: map[string][]string{
				tmp: {filepath.Join(tmp, "b.log")},
			},
		},
		{
			desc:     "match all .stache everywhere (manual)",
			patterns: []string{"*.stache", "dir/*.stache"},
			expected: map[string][]string{
				tmp:    {filepath.Join(tmp, "a.stache")},
				subDir: {filepath.Join(subDir, "c.stache")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := resolveGlobs(tmp, tt.patterns)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected:\n%v\ngot:\n%v", tt.expected, got)
			}
		})
	}

	t.Run("directory match is skipped", func(t *testing.T) {
		tmp := t.TempDir()
		dir := filepath.Join(tmp, "matchdir")
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatal(err)
		}

		got, err := resolveGlobs(tmp, []string{"matchdir"})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Errorf("expected no files, got: %v", got)
		}
	})

	t.Run("glob pattern error", func(t *testing.T) {
		tmp := t.TempDir()

		_, err := resolveGlobs(tmp, []string{"["})
		if err == nil {
			t.Fatal("expected error from filepath.Glob, got nil")
		}
		if !strings.Contains(err.Error(), "failed to match files") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
