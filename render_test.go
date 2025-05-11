package restache_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tetsuo/restache"
)

func TestRender(t *testing.T) {
	const file = "testdata/render_jsx.txt"
	for _, tc := range buildTestcases(t, file) {
		t.Run(fmt.Sprintf("%s L%d", file, tc.line), func(t *testing.T) {
			root, err := restache.Parse(strings.NewReader(tc.data))
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			var sb strings.Builder
			_, err = root.Render(&sb)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			got := sb.String()
			want := "export default function ($0) {return " + tc.expected + ";}"

			if got != want {
				t.Errorf("Render mismatch at line %d:\nwant:\n%s\ngot:\n%s\n", tc.line, want, got)
			}
		})
	}
}
