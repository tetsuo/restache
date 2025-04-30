package restache

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html/atom"
)

func Render(w io.Writer, n *Node) (int, error) {
	if x, ok := w.(writer); ok {
		r := &renderer{w: x}
		if err := r.render(n); err != nil {
			return 0, err
		}
		return r.written, nil
	}
	buf := bufio.NewWriter(w)
	r := &renderer{w: buf}
	if err := r.render(n); err != nil {
		return 0, err
	}
	if err := buf.Flush(); err != nil {
		return 0, err
	}
	return r.written, nil
}

var (
	ErrErrorNode    = errors.New("cannot render an ErrorNode node")
	ErrUnknownNode  = errors.New("unknown node type")
	ErrVoidChildren = errors.New("void element has child nodes")
)

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

type renderer struct {
	w writer

	indent  int
	written int
	rlvl    int
}

func (r *renderer) print1(c byte) error {
	if err := r.w.WriteByte(c); err != nil {
		return err
	} else {
		r.written += 1
	}
	return nil
}

func (r *renderer) print(s string) error {
	if n, err := r.w.WriteString(s); err != nil {
		return err
	} else {
		r.written += n
	}
	return nil
}

func (r *renderer) println(s string) error {
	err := r.print(s)
	if err == nil {
		return r.print1('\n')
	}
	return err
}

func (r *renderer) printf(format string, args ...any) error {
	if n, err := r.w.WriteString(fmt.Sprintf(format, args...)); err != nil {
		return err
	} else {
		r.written += n
	}
	return nil
}

func (r *renderer) lineBreak() error {
	if err := r.print1('\n'); err != nil {
		return err
	}
	if r.indent < len(indentStrings) {
		return r.print(indentStrings[r.indent])
	}
	return r.print(strings.Repeat("  ", r.indent))
}

func (r *renderer) renderText(n *Node) error {
	s := n.Data
	if s[0] == ' ' && n.PrevSibling == nil {
		s = s[1:]
	}
	x := len(s)
	if s[x-1] == ' ' && n.NextSibling == nil {
		s = s[:x-1]
	}
	if err := r.print(s); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderVariable(n *Node) error {
	if err := r.printf("{ d%d.", r.rlvl); err != nil {
		return err
	}
	if err := r.print(n.Data); err != nil {
		return err
	}
	if err := r.print(" }"); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderComponent(n *Node) error {
	if n.FirstChild == nil {
		return r.println("export default function() { return null; }")
	}
	if err := r.println("import * as React from 'react';"); err != nil {
		return err
	}
	for _, attr := range n.Attr {
		if err := r.printf("import %s from \"./%s.jsx\";\n", attr.Key, attr.Val); err != nil {
			return err
		}
	}
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.printf("export default function %s(d%d) {", n.Data, r.rlvl); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print("return ("); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if n.FirstChild == n.LastChild && n.FirstChild.Type == ElementNode {
		if err := r.renderElement(n.FirstChild); err != nil {
			return err
		}
	} else {
		n.Type = ElementNode
		n.Data = ""
		if err := r.renderElement(n); err != nil {
			return err
		}
	}
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.println(");\n}"); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderWhen(n *Node) error {
	if err := r.print1('{'); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.printf("d%d.%s && (", r.rlvl, n.Data); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if n.FirstChild == n.LastChild && n.FirstChild.Type == ElementNode {
		if err := r.renderElement(n.FirstChild); err != nil {
			return err
		}
	} else {
		n.Type = ElementNode
		n.Data = ""
		if err := r.renderElement(n); err != nil {
			return err
		}
	}
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print1(')'); err != nil {
		return err
	}
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print1('}'); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderUnless(n *Node) error {
	if err := r.print1('{'); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.printf("!d%d.%s && (", r.rlvl, n.Data); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if n.FirstChild == n.LastChild && n.FirstChild.Type == ElementNode {
		if err := r.renderElement(n.FirstChild); err != nil {
			return err
		}
	} else {
		n.Type = ElementNode
		n.Data = ""
		if err := r.renderElement(n); err != nil {
			return err
		}
	}
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print1(')'); err != nil {
		return err
	}
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print1('}'); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderRange(n *Node) error {
	if err := r.print1('{'); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.printf("d%d.%s.map(", r.rlvl, n.Data); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}

	r.rlvl++
	if err := r.printf("d%d => (", r.rlvl); err != nil {
		return err
	}
	r.indent++
	if err := r.lineBreak(); err != nil {
		return err
	}
	if n.FirstChild == n.LastChild && n.FirstChild.Type == ElementNode {
		n.FirstChild.Attr = append(n.FirstChild.Attr, Attribute{Key: "key", Val: "key", IsExpr: true})

		if err := r.renderElement(n.FirstChild); err != nil {
			return err
		}
	} else {
		n.Type = ElementNode
		n.Data = "React.Fragment"
		n.Attr = append(n.Attr, Attribute{Key: "key", Val: "key", IsExpr: true})
		if err := r.renderElement(n); err != nil {
			return err
		}
	}
	r.rlvl--
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print1(')'); err != nil {
		return err
	}
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print1(')'); err != nil {
		return err
	}
	r.indent--
	if err := r.lineBreak(); err != nil {
		return err
	}
	if err := r.print1('}'); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderAttribute(a Attribute) error {
	if err := r.print1(' '); err != nil {
		return err
	}
	if err := r.print(a.Key); err != nil {
		return err
	}
	if a.Val == "" && a.KeyAtom != 0 {
		if _, ok := boolAttrs[a.KeyAtom]; ok {
			return nil
		}
	}
	if a.IsExpr {
		if err := r.printf(`={ d%d.%s }`, r.rlvl, a.Val); err != nil {
			return err
		}
	} else {
		if err := r.printf(`="%s"`, a.Val); err != nil {
			return err
		}
	}
	return nil
}

func (r *renderer) renderElement(n *Node) error {
	// Render the <xxx> opening tag.
	if err := r.print1('<'); err != nil {
		return err
	}

	var tagName string
	if n.DataAtom != 0 {
		tagName = n.DataAtom.String()
	} else {
		tagName = n.Data
	}

	if err := r.print(tagName); err != nil {
		return err
	}

	if len(n.Attr) > 0 {
		if _, found := camelAttrTags[n.DataAtom]; found {
			searchPrefix := uint64(n.DataAtom) << 32
			for _, a := range n.Attr {
				if a.KeyAtom != 0 {
					if alias, ok := globalCamelAttrTable[a.KeyAtom]; ok {
						a.Key = alias
					} else if alias, ok := camelAttrTable[searchPrefix|uint64(a.KeyAtom)]; ok {
						a.Key = alias
					} else {
						a.Key = a.KeyAtom.String()
					}
				}
				if err := r.renderAttribute(a); err != nil {
					return err
				}
			}
		} else {
			for _, a := range n.Attr {
				if a.KeyAtom != 0 {
					if alias, ok := globalCamelAttrTable[a.KeyAtom]; ok {
						a.Key = alias
					} else {
						a.Key = a.KeyAtom.String()
					}
				}
				if err := r.renderAttribute(a); err != nil {
					return err
				}
			}
		}
	}
	if _, ok := voidElements[n.DataAtom]; ok {
		if n.FirstChild != nil {
			return ErrVoidChildren
		}
		err := r.print(" />")
		return err
	}

	c := n.FirstChild
	if c != nil {
		if err := r.print1('>'); err != nil {
			return err
		}
		// Add initial newline where there is danger of a newline being ignored.
		if c.Type == TextNode && strings.HasPrefix(c.Data, "\n") {
			switch n.DataAtom {
			case atom.Pre, atom.Listing, atom.Textarea:
				if err := r.print1('\n'); err != nil {
					return err
				}
			}
		}

		r.indent++
		for c != nil {
			if err := r.render(c); err != nil {
				return err
			}
			c = c.NextSibling
		}
		r.indent--
	} else {
		if err := r.print1('>'); err != nil {
			return err
		}
	}
	if n.FirstChild != nil {
		if err := r.lineBreak(); err != nil {
			return err
		}
	}
	// Render the </xxx> closing tag.
	if err := r.print("</"); err != nil {
		return err
	}
	if err := r.print(n.TagName()); err != nil {
		return err
	}
	if err := r.print1('>'); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderComment(n *Node) error {
	if len(n.Data) < 80 {
		if err := r.print("{ /* "); err != nil {
			return err
		}
		if err := escapeComment(r.w, n.Data); err != nil {
			return err
		}
		if err := r.print(" */ }"); err != nil {
			return err
		}
	} else {
		if err := r.print("{/*"); err != nil {
			return err
		}
		if err := r.lineBreak(); err != nil {
			return err
		}
		if err := escapeComment(r.w, n.Data); err != nil {
			return err
		}
		if err := r.lineBreak(); err != nil {
			return err
		}
		if err := r.print("*/}"); err != nil {
			return err
		}
	}
	return nil
}

func (r *renderer) render(n *Node) error {
	switch n.Type {
	case ErrorNode:
		return ErrErrorNode
	case TextNode:
		if n.PrevSibling == nil || !(n.PrevSibling.Type == TextNode || n.PrevSibling.Type == VariableNode) {
			if err := r.lineBreak(); err != nil {
				return err
			}
		}
		return r.renderText(n)
	case ElementNode:
		if err := r.lineBreak(); err != nil {
			return err
		}
		return r.renderElement(n)
	case VariableNode:
		if n.PrevSibling == nil || !(n.PrevSibling.Type == TextNode || n.PrevSibling.Type == VariableNode) {
			if err := r.lineBreak(); err != nil {
				return err
			}
		}
		return r.renderVariable(n)
	case WhenNode:
		if err := r.lineBreak(); err != nil {
			return err
		}
		return r.renderWhen(n)
	case UnlessNode:
		if err := r.lineBreak(); err != nil {
			return err
		}
		return r.renderUnless(n)
	case RangeNode:
		if err := r.lineBreak(); err != nil {
			return err
		}
		return r.renderRange(n)
	case CommentNode:
		if err := r.lineBreak(); err != nil {
			return err
		}
		return r.renderComment(n)
	case ComponentNode:
		return r.renderComponent(n)
	default:
		return ErrUnknownNode
	}
}

var indentStrings = [32]string{
	"",
	"  ",
	"    ",
	"      ",
	"        ",
	"          ",
	"            ",
	"              ",
	"                ",
	"                  ",
	"                    ",
	"                      ",
	"                        ",
	"                          ",
	"                            ",
	"                              ",
	"                                ",
	"                                  ",
	"                                    ",
	"                                      ",
	"                                        ",
	"                                          ",
	"                                            ",
	"                                              ",
	"                                                ",
	"                                                  ",
	"                                                    ",
	"                                                      ",
	"                                                        ",
	"                                                          ",
	"                                                            ",
	"                                                              ",
}
