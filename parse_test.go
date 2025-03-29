package stache_test

import (
	"bytes"
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

			actual := normalizeSpaces(dumpResolvedTree(root))
			expected := normalizeSpaces(tc.expected)

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
		res := dumpResolvedTree(root)
		expected := "[] [  ]"
		if res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	})
}

func dumpResolvedTree(n *stache.Node) string {
	var b strings.Builder
	b.WriteString("[] [ ")
	dumpResolvedNode(&b, n, 0)
	b.WriteString(" ]")
	return b.String()
}

func dumpResolvedNode(b *strings.Builder, n *stache.Node, indent int) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		indentLine(b, indent)

		switch c.Type {
		case stache.TextNode:
			b.WriteString(`text "`)
			b.Write(c.Data)
			b.WriteString(`"`)

		case stache.CommentNode:
			b.WriteString(`comment "`)
			b.Write(c.Data)
			b.WriteString(`"`)

		case stache.VariableNode:
			b.WriteString(`var `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.Write(c.Data)

		case stache.WhenNode:
			b.WriteString(`when `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.Write(c.Data)
			b.WriteString(` [`)
			b.WriteByte('\n')
			dumpResolvedNode(b, c, indent+2)
			indentLine(b, indent)
			b.WriteString(`]`)

		case stache.UnlessNode:
			b.WriteString(`unless `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.Write(c.Data)
			b.WriteString(` [`)
			b.WriteByte('\n')
			dumpResolvedNode(b, c, indent+2)
			indentLine(b, indent)
			b.WriteString(`]`)

		case stache.RangeNode:
			b.WriteString(`range `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.Write(c.Data)
			b.WriteString(` [`)
			b.WriteByte('\n')
			dumpResolvedNode(b, c, indent+2)
			indentLine(b, indent)
			b.WriteString(`]`)

		case stache.ElementNode:
			b.Write(c.Data)     // tag name
			b.WriteString(" [") // attributes
			for j, attr := range c.Attr {
				if j > 0 {
					b.WriteString(", ")
				}
				b.Write(attr.Key)
				if len(attr.Val) > 0 {
					b.WriteString(" ")
					if attr.IsExpr {
						b.WriteString(`var `)
						for i, seg := range c.Path {
							if i > 0 {
								b.WriteByte('.')
							}
							b.Write(seg)
						}

						segments := bytes.Split(attr.Val, []byte("."))
						for _, seg := range segments {
							b.WriteByte('.')
							b.Write(seg)
						}
					} else {
						b.WriteString(`text "`)
						b.Write(attr.Val)
						b.WriteString(`"`)
					}
				}
			}
			b.WriteString("] [") // children
			dumpResolvedNode(b, c, indent+2)
			indentLine(b, indent)
			b.WriteString("]")

		default:
			b.WriteString("???")
		}

		if c.NextSibling != nil {
			b.WriteString(",\n")
		} else {
			b.WriteByte('\n')
		}
	}
}

func writeResolvedPath(b *strings.Builder, n *stache.Node) {
	for i, seg := range n.Path {
		if i > 0 {
			b.WriteByte('.')
		}
		b.Write(seg)
	}
}

func indentLine(b *strings.Builder, n int) {
	for i := 0; i < n; i++ {
		b.WriteByte(' ')
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
