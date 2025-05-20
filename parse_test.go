package restache_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tetsuo/restache"
)

func TestParse(t *testing.T) {
	const file = "testdata/parser_ast.txt"
	for _, tc := range buildTestcases(t, file) {
		t.Run(fmt.Sprintf("%s L%d", file, tc.line), func(t *testing.T) {
			r := strings.NewReader(tc.data)
			root, err := restache.Parse(r)
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
		node, err := restache.Parse(&errorReader{})
		if err.Error() != "test error" {
			t.Errorf("expected test error, got %v", err)
		}
		if node != nil {
			t.Errorf("expected nil Node, got %v", node)
		}
	})

	t.Run("empty text becomes fragment", func(t *testing.T) {
		root, err := restache.Parse(strings.NewReader("  "))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		res := normalizeSpaces(dumpResolvedTree(root))
		expected := "[[][]]"
		if res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	})
}

func TestNodePanic(t *testing.T) {
	checkPanic := func(expected string, actual any) {
		if msg, ok := actual.(string); ok {
			if msg != expected {
				t.Errorf("expected panic message %q, got %q", expected, msg)
			}
		} else {
			t.Errorf("expected string panic, got %T: %v", actual, actual)
		}
	}

	t.Run("Node; RemoveChild called for a non-child Node", func(t *testing.T) {
		didPanic := false

		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
					checkPanic("restache: RemoveChild called for a non-child Node", r)
				}
			}()

			n1 := &restache.Node{}
			n2 := &restache.Node{}
			n1.RemoveChild(n2)
		}()

		if !didPanic {
			t.Error("expected panic, but function did not panic")
		}
	})

	t.Run("Node; AppendChild called for an attached child Node", func(t *testing.T) {
		didPanic := false

		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
					checkPanic("restache: AppendChild called for an attached child Node", r)
				}
			}()

			p1 := &restache.Node{}
			c1 := &restache.Node{}
			p1.AppendChild(c1)
			p2 := &restache.Node{}
			p2.AppendChild(c1)
		}()

		if !didPanic {
			t.Error("expected panic, but function did not panic")
		}
	})
}

func dumpResolvedTree(n *restache.Node) string {
	var b strings.Builder
	b.WriteString("[ ")
	dumpResolvedNode(&b, n, 0)
	b.WriteString(" ]")
	return b.String()
}

func dumpResolvedNode(b *strings.Builder, n *restache.Node, indent int) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		indentLine(b, indent)

		switch c.Type {
		case restache.TextNode:
			b.WriteString(`text "`)
			b.WriteString(c.Data)
			b.WriteString(`"`)

		case restache.CommentNode:
			b.WriteString(`comment "`)
			b.WriteString(c.Data)
			b.WriteString(`"`)

		case restache.VariableNode:
			b.WriteString(`var `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.WriteString(c.Data)

		case restache.WhenNode:
			b.WriteString(`when `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.WriteString(c.Data)
			b.WriteString(` [`)
			b.WriteByte('\n')
			dumpResolvedNode(b, c, indent+2)
			indentLine(b, indent)
			b.WriteString(`]`)

		case restache.UnlessNode:
			b.WriteString(`unless `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.WriteString(c.Data)
			b.WriteString(` [`)
			b.WriteByte('\n')
			dumpResolvedNode(b, c, indent+2)
			indentLine(b, indent)
			b.WriteString(`]`)

		case restache.RangeNode:
			b.WriteString(`range `)
			writeResolvedPath(b, c)
			b.WriteByte('.')
			b.WriteString(c.Data)
			b.WriteString(` [`)
			b.WriteByte('\n')
			dumpResolvedNode(b, c, indent+2)
			indentLine(b, indent)
			b.WriteString(`]`)

		case restache.ElementNode:
			b.WriteString(c.TagName()) // tag name
			b.WriteString(" [")        // attributes
			for j, attr := range c.Attr {
				if j > 0 {
					b.WriteString(", ")
				}
				b.WriteString(attrKeyName(attr))
				if len(attr.Val) > 0 {
					b.WriteString(" ")
					if attr.IsExpr {
						b.WriteString(`var `)
						for i, seg := range c.Path {
							if i > 0 {
								b.WriteByte('.')
							}
							b.WriteString(seg.Key)
							if seg.IsRange {
								b.WriteString(".#")
							}
						}

						segments := strings.Split(attr.Val, ".")
						for _, seg := range segments {
							b.WriteByte('.')
							b.WriteString(seg)
						}
					} else {
						b.WriteString(`text "`)
						b.WriteString(attr.Val)
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

func attrKeyName(a restache.Attribute) string {
	if a.KeyAtom != 0 {
		return a.KeyAtom.String()
	}
	return a.Key
}

func writeResolvedPath(b *strings.Builder, n *restache.Node) {
	for i, seg := range n.Path {
		if i > 0 {
			b.WriteByte('.')
		}
		b.WriteString(seg.Key)
		if seg.IsRange {
			b.WriteString(".#")
		}
	}
}

func indentLine(b *strings.Builder, n int) {
	for range n {
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
