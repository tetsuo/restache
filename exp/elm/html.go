// Package elm provides an Elm-style HTML renderer intended for use in parser tests.
// It is not intended for general use.
package elm

import (
	"strings"

	"github.com/tetsuo/stache"
)

func Dump(buf *strings.Builder, rootName string, n *stache.Node) {
	buf.WriteString(rootName + " [] [ ")
	dumpNode(buf, n)
	buf.WriteString(" ]")
}

func dumpNode(buf *strings.Builder, n *stache.Node) {
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
			dumpNode(buf, c)
			buf.WriteString("]")

		case stache.WhenNode:
			buf.WriteString(`when "` + string(c.Data) + `" [`)
			dumpNode(buf, c)
			buf.WriteString("]")

		case stache.UnlessNode:
			buf.WriteString(`unless "` + string(c.Data) + `" [`)
			dumpNode(buf, c)
			buf.WriteString("]")

		case stache.RangeNode:
			buf.WriteString(`range "` + string(c.Data) + `" [`)
			dumpNode(buf, c)
			buf.WriteString("]")

		default:
			buf.WriteString("???")
		}
	}
}

func NormalizeWhitespace(s string) string {
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
