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
	ErrErrorNode       = errors.New("cannot render an ErrorNode node")
	ErrUnknownNode     = errors.New("unknown node type")
	ErrVoidChildren    = errors.New("void element has child nodes")
	ErrTooManyChildren = errors.New("element must have a single child")
)

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

type renderer struct {
	w writer

	written int
	scope   int
}

func (r *renderer) print1(c byte) (err error) {
	if err = r.w.WriteByte(c); err == nil {
		r.written += 1
	}
	return
}

func (r *renderer) print(s string) (err error) {
	var n int
	if n, err = r.w.WriteString(s); err == nil {
		r.written += n
	}
	return
}

func (r *renderer) printf(format string, args ...any) (err error) {
	var n int
	if n, err = fmt.Fprintf(r.w, format, args...); err == nil {
		r.written += n
	}
	return
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

func (r *renderer) renderVariable(n *Node) (err error) {
	if err = r.printf("$%d.", r.scope); err == nil {
		err = r.print(n.Data)
	}
	return
}

func (r *renderer) renderComponent(n *Node) error {
	for _, attr := range n.Attr {
		if err := r.printf("import %s from '%s';\n", attr.Key, attr.Val); err != nil {
			return err
		}
	}
	if err := r.printf("export default function %s($%d) {", n.Data, r.scope); err != nil {
		return err
	}
	first := n.FirstChild
	if first != nil {
		if err := r.print("return "); err != nil {
			return err
		}
		if err := r.render(first); err != nil {
			return err
		}
		if err := r.print(";}"); err != nil {
			return err
		}
	} else {
		if err := r.print("return null;"); err != nil {
			return err
		}
		if err := r.print("}"); err != nil {
			return err
		}
	}
	return nil
}

func (r *renderer) renderWhen(n *Node, negate bool) error {
	if err := r.print1('('); err != nil {
		return err
	}
	if negate {
		if err := r.print1('!'); err != nil {
			return err
		}
	}
	if err := r.printf("$%d.%s && ", r.scope, n.Data); err != nil {
		return err
	}
	if n.FirstChild != nil && n.FirstChild == n.LastChild {
		if err := r.render(n.FirstChild); err != nil {
			return err
		}
	} else {
		// TODO: raise error
	}
	if err := r.print1(')'); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderRange(n *Node) error {
	if err := r.printf("$%d.%s.map(", r.scope, n.Data); err != nil {
		return err
	}
	r.scope++
	if err := r.printf("$%d => ", r.scope); err != nil {
		return err
	}
	if n.FirstChild != nil && n.FirstChild == n.LastChild {
		if err := r.render(n.FirstChild); err != nil {
			return err
		}
	} else {
		// TODO: raise too many
	}
	r.scope--
	if err := r.print1(')'); err != nil {
		return err
	}
	return nil
}

func (r *renderer) renderAttribute(a Attribute, key string) error {
	if err := r.print1(' '); err != nil {
		return err
	}
	if err := r.print(key); err != nil {
		return err
	}
	if a.Val == "" && a.KeyAtom != 0 {
		if _, ok := boolAttrs[a.KeyAtom]; ok {
			return nil
		}
	}
	if a.IsExpr {
		return r.printf(`={ $%d.%s }`, r.scope, a.Val)
	}
	return r.printf(`="%s"`, a.Val)
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
				if err := r.renderAttribute(a, a.Key); err != nil {
					return err
				}
			}
		} else {
			for _, a := range n.Attr {
				var key string
				if a.KeyAtom != 0 {
					if alias, ok := globalCamelAttrTable[a.KeyAtom]; ok {
						key = alias
					} else {
						key = a.KeyAtom.String()
					}
				} else {
					key = a.Key
				}
				if err := r.renderAttribute(a, key); err != nil {
					return err
				}
			}
		}
	}

	if n.DataAtom != 0 {
		if _, ok := voidElements[n.DataAtom]; ok {
			if n.FirstChild != nil {
				return ErrVoidChildren
			}
			err := r.print(" />")
			return err
		}
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
		for c != nil {
			if err := r.render(c); err != nil {
				return err
			}
			c = c.NextSibling
		}
	} else {
		if err := r.print1('>'); err != nil {
			return err
		}
	}
	// Render the </xxx> closing tag.
	if err := r.print("</"); err != nil {
		return err
	}
	if err := r.print(tagName); err != nil {
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
		if err := escapeComment(r.w, n.Data); err != nil {
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

func escapeComment(w writer, s string) error {
	if len(s) == 0 {
		return nil
	}

	i := 0
	for j := 0; j < len(s)-1; j++ {
		if s[j] == '*' && s[j+1] == '/' {
			if i < j {
				if _, err := w.WriteString(s[i:j]); err != nil {
					return err
				}
			}
			if _, err := w.WriteString("*\\/"); err != nil { // escape the '/'
				return err
			}
			i = j + 2
			j++ // skip the '/'
		}
	}

	if i < len(s) {
		if _, err := w.WriteString(s[i:]); err != nil {
			return err
		}
	}
	return nil
}
