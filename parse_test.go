package stache_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tetsuo/stache"
	"github.com/tetsuo/stache/exp/elm"
)

func TestParse(t *testing.T) {
	const file = "testdata/parser_testcases.txt"
	for _, tc := range buildTestcases(t, file) {
		t.Run(fmt.Sprintf("%s L%d", file, tc.line), func(t *testing.T) {
			r := strings.NewReader(tc.data)
			root, err := stache.Parse(r)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			var buf strings.Builder
			elm.Dump(&buf, "main", root)

			expected := elm.NormalizeWhitespace(tc.expected)
			actual := elm.NormalizeWhitespace(buf.String())

			if actual != expected {
				t.Errorf("\nexpected:\n%s\ngot:\n%s\n", expected, actual)
			}
		})
	}

	t.Run("faulty reader", func(t *testing.T) {
		node, err := stache.Parse(&errorReader{})
		if err.Error() != "test error" {
			t.Errorf("expected test error, got %v", err)
		}
		if node != nil {
			t.Errorf("expected nil Node, got %v", node)
		}
	})

	t.Run("skip empty text", func(t *testing.T) {
		root, err := stache.Parse(strings.NewReader("  "))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		var buf strings.Builder
		elm.Dump(&buf, "_", root)
		expected := "_ [] [  ]"
		if buf.String() != expected {
			t.Errorf("expected %q, got %q", expected, buf.String())
		}
	})
}
