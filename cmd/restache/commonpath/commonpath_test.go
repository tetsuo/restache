package commonpath_test

import (
	"testing"

	"github.com/tetsuo/restache/cmd/restache/commonpath"
)

func TestCommonPath(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected string
		platform string
	}{
		// UNIX-style
		{
			name:     "no paths",
			paths:    []string{},
			expected: "",
			platform: "unix",
		},
		{
			name:     "single path",
			paths:    []string{"/a/b/c"},
			expected: "/a/b/c",
			platform: "unix",
		},
		{
			name:     "single root path",
			paths:    []string{"/"},
			expected: "/",
			platform: "unix",
		},
		{
			name:     "root only",
			paths:    []string{"/a", "/b", "/c"},
			expected: "/",
			platform: "unix",
		},
		{
			name:     "identical",
			paths:    []string{"/x/y/z", "/x/y/z", "/x/y/z"},
			expected: "/x/y/z",
			platform: "unix",
		},
		{
			name:     "nested common",
			paths:    []string{"/a/b/c/d", "/a/b/c/e", "/a/b/c/f/g"},
			expected: "/a/b/c",
			platform: "unix",
		},
		{
			name:     "common /a/b",
			paths:    []string{"/a/b/c", "/a/b/d", "/a/b/e/f"},
			expected: "/a/b",
			platform: "unix",
		},
		{
			name:     "two root paths",
			paths:    []string{"/", "/"},
			expected: "/",
			platform: "unix",
		},
		{
			name:     "mixed no overlap",
			paths:    []string{"/foo", "/bar"},
			expected: "/",
			platform: "unix",
		},

		// Windows-style
		{
			name:     "win single path root",
			paths:    []string{`C:\`},
			expected: `C:\`,
			platform: "win",
		},
		{
			name:     "win single segment drive only",
			paths:    []string{`C:`},
			expected: `C:\`,
			platform: "win",
		},
		{
			name:     "win common drive",
			paths:    []string{`C:\x\y\z`, `C:\x\y\m`, `C:\x\y\n`},
			expected: `C:\x\y`,
			platform: "win",
		},
		{
			name:     "win case fold",
			paths:    []string{`C:\A\B\C`, `c:\a\b\d`},
			expected: `C:\A\B`,
			platform: "win",
		},
		{
			name:     "win trailing slash",
			paths:    []string{`C:\foo\`, `C:\foo\bar`},
			expected: `C:\foo`,
			platform: "win",
		},
		{
			name:     "win different drives",
			paths:    []string{`C:\x`, `D:\y`},
			expected: `\`,
			platform: "win",
		},
		{
			name:     "win UNC",
			paths:    []string{`\\server\share\folder1\sub`, `\\server\share\folder2`},
			expected: `\\server\share`,
			platform: "win",
		},
		{
			name:     "win relative plus absolute",
			paths:    []string{`C:\folder\sub`, `C:\folder\sub\child`},
			expected: `C:\folder\sub`,
			platform: "win",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			if tt.platform == "unix" {
				got = commonpath.CommonPathUnix(tt.paths)
			} else {
				got = commonpath.CommonPathWin(tt.paths)
			}
			if got != tt.expected {
				t.Errorf("want %q, got %q", tt.expected, got)
			}
		})
	}
}
