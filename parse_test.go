package stache_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tetsuo/stache"
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
			dumpElmRoot(&buf, "main", root)

			expected := normalizeSpaces(tc.expected)
			actual := normalizeSpaces(buf.String())

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
		dumpElmRoot(&buf, "_", root)
		expected := "_ [] [  ]"
		if buf.String() != expected {
			t.Errorf("expected %q, got %q", expected, buf.String())
		}
	})
}

func dumpElmRoot(buf *strings.Builder, rootName string, n *stache.Node) {
	buf.WriteString(rootName + " [] [ ")
	dumpElmNode(buf, n)
	buf.WriteString(" ]")
}

func dumpElmNode(buf *strings.Builder, n *stache.Node) {
	for c, i := n.FirstChild, 0; c != nil; c, i = c.NextSibling, i+1 {
		if i > 0 {
			buf.WriteString(", ")
		}

		switch c.Type {
		case stache.TextNode:
			buf.WriteString(`text "` + string(c.Data) + `"`)

		case stache.CommentNode:
			buf.WriteString(`comment "` + string(c.Data) + `"`)

		case stache.VariableNode:
			buf.WriteString(`var "` + string(c.Data) + `"`)

		case stache.ElementNode:
			buf.WriteString(string(c.Data)) // tag name
			buf.WriteString(" [")           // attributes
			for j, attr := range c.Attr {
				if j > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(string(attr.Key))
				if len(attr.Val) > 0 {
					buf.WriteString(" ")
					if attr.IsExpr {
						buf.WriteString(`var "`)
					} else {
						buf.WriteString(`text "`)
					}
					buf.WriteString(string(attr.Val))
					buf.WriteString(`"`)
				}
			}
			buf.WriteString("] [") // children
			dumpElmNode(buf, c)
			buf.WriteString("]")

		case stache.WhenNode:
			buf.WriteString(`when "` + string(c.Data) + `" [`)
			dumpElmNode(buf, c)
			buf.WriteString("]")

		case stache.UnlessNode:
			buf.WriteString(`unless "` + string(c.Data) + `" [`)
			dumpElmNode(buf, c)
			buf.WriteString("]")

		case stache.RangeNode:
			buf.WriteString(`range "` + string(c.Data) + `" [`)
			dumpElmNode(buf, c)
			buf.WriteString("]")

		default:
			buf.WriteString("???")
		}
	}
}

func normalizeSpaces(s string) string {
	var buf strings.Builder
	inQuotes := false
	previousCharWasSpace := false

	for i := range len(s) {
		ch := s[i]

		switch ch {
		case '"':
			inQuotes = !inQuotes
			buf.WriteByte(ch)
			previousCharWasSpace = false

		case '[', ']', ',', ':':
			if !inQuotes {
				if previousCharWasSpace && buf.Len() > 0 && buf.String()[buf.Len()-1] == ' ' {
					b := buf.String()
					buf.Reset()
					buf.WriteString(b[:len(b)-1])
				}
				buf.WriteByte(ch)
				buf.WriteByte(' ')
				previousCharWasSpace = true
			} else {
				buf.WriteByte(ch)
			}

		case ' ', '\t', '\n', '\r':
			if inQuotes {
				buf.WriteByte(ch)
			} else if !previousCharWasSpace {
				buf.WriteByte(' ')
				previousCharWasSpace = true
			}

		default:
			buf.WriteByte(ch)
			previousCharWasSpace = false
		}
	}

	return strings.TrimSpace(buf.String())
}
