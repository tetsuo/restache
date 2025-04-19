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
		r := &renderer{w: x, indent: 2}
		if err := r.render(n); err != nil {
			return 0, err
		}
		return r.written, nil
	}
	buf := bufio.NewWriter(w)
	r := &renderer{w: buf, indent: 2}
	if err := r.render(n); err != nil {
		return 0, err
	}
	if err := buf.Flush(); err != nil {
		return 0, err
	}
	return r.written, nil
}

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

type renderer struct {
	w writer

	indent  int
	written int
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

func (r *renderer) printf(format string, args ...any) error {
	if n, err := r.w.WriteString(fmt.Sprintf(format, args...)); err != nil {
		return err
	} else {
		r.written += n
	}
	return nil
}

func (r *renderer) renderText(n *Node) error {
	return r.print(n.Data)
}

func (r *renderer) renderComponent(n *Node) error {
	if err := r.print("import * as React from 'react';"); err != nil {
		return err
	}
	for _, attr := range n.Attr {
		if err := r.printf("import %s from \"./%s.jsx\";\n", attr.Key, attr.Val); err != nil {
			return err
		}
	}
	if err := r.print1('\n'); err != nil {
		return err
	}
	if err := r.print("export default function"); err != nil {
		return err
	}
	if n.Data != "" {
		if err := r.print1(' '); err != nil {
			return err
		}
		if err := r.print(n.Data); err != nil {
			return err
		}
	}
	if err := r.print("(props) {\n"); err != nil {
		return err
	}
	if err := r.print("  return ("); err != nil {
		return err
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := r.render(c); err != nil {
			return err
		}
	}
	if err := r.print("  );"); err != nil {
		return err
	}
	if err := r.print("}"); err != nil {
		return err
	}
	return nil
}

func (r *renderer) writeIndent() error {
	for range r.indent {
		if err := r.print("  "); err != nil {
			return err
		}
	}
	return nil
}

func (r *renderer) render(n *Node) error {
	switch n.Type {
	case ErrorNode:
		return errors.New("html: cannot render an ErrorNode node")
	case TextNode:
		return r.renderText(n)
	case ComponentNode:
		return r.renderComponent(n)
	case ElementNode:
		return r.renderElement(n)
	case CommentNode:
		return r.renderComment(n)
	case VariableNode:
		return r.renderVariable(n)
	default:
		return errors.New("html: unknown node type")
	}
}

func (r *renderer) renderTextAttribute(a Attribute) error {
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
	if err := r.printf(`="%s"`, a.Val); err != nil {
		return err
	}
	return nil
}

// renderElementNode renders a normal element.
func (r *renderer) renderElement(n *Node) error {
	if err := r.writeIndent(); err != nil {
		return err
	}
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
				if err := r.renderTextAttribute(a); err != nil {
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
				if err := r.renderTextAttribute(a); err != nil {
					return err
				}
			}
		}
	}
	if _, ok := voidElements[n.DataAtom]; ok {
		if n.FirstChild != nil {
			return fmt.Errorf("html: void element <%s> has child nodes", n.Data)
		}
		err := r.print("/>")
		return err
	}
	if err := r.print(">\n"); err != nil {
		return err
	}

	// Add initial newline where there is danger of a newline being ignored.
	if c := n.FirstChild; c != nil && c.Type == TextNode && strings.HasPrefix(c.Data, "\n") {
		switch n.DataAtom {
		case atom.Pre, atom.Listing, atom.Textarea:
			if err := r.print1('\n'); err != nil {
				return err
			}
		}
	}

	r.indent++

	for c := n.FirstChild; c != nil; {
		if c.Type == WhenNode || c.Type == UnlessNode {
			group := []*Node{}
			cond := c.Data
			for c != nil && (c.Type == WhenNode || c.Type == UnlessNode) && c.Data == cond {
				group = append(group, c)
				c = c.NextSibling
			}
			if err := r.renderConditionalGroup(group); err != nil {
				return err
			}
		} else {
			if err := r.render(c); err != nil {
				return err
			}
			c = c.NextSibling
		}
	}

	r.indent--
	if err := r.writeIndent(); err != nil {
		return err
	}
	// Render the </xxx> closing tag.
	if err := r.print("</"); err != nil {
		return err
	}
	if err := r.print(n.TagName()); err != nil {
		return err
	}
	if err := r.print(">\n"); err != nil {
		return err
	}
	return nil
}

// renderVariable renders a variable JSX expression.
// Example:
//
//	{props.someValue}
func (r *renderer) renderVariable(n *Node) error {
	if err := r.writeIndent(); err != nil {
		return err
	}
	if err := r.print("{props."); err != nil {
		return err
	}
	if err := r.print(n.Data); err != nil {
		return err
	}
	if err := r.print("}\n"); err != nil {
		return err
	}
	return nil
}

// renderComment renders a comment JSX expression.
// Example:
//
//	{/* comment */}
func (r *renderer) renderComment(n *Node) error {
	if err := r.writeIndent(); err != nil {
		return err
	}
	if err := r.print("{/* "); err != nil {
		return err
	}
	if err := escapeComment(r.w, n.Data); err != nil {
		return err
	}
	if err := r.print(" */}\n"); err != nil {
		return err
	}
	return nil
}

// conditionalBlock is a small struct for when/unless blocks in a group.
type conditionalBlock struct {
	Cond   string
	Negate bool
	Body   *Node
}

// renderConditionalGroup renders a conditional JSX expression based on the node group.
// The function selects the appropriate rendering method (single, ternary, or multi-ternary).
func (r *renderer) renderConditionalGroup(group []*Node) error {
	// Convert each node into a conditionalBlock
	blocks := make([]conditionalBlock, len(group))
	for i, n := range group {
		blocks[i] = conditionalBlock{
			Cond:   n.Data,
			Negate: (n.Type == UnlessNode),
			Body:   n,
		}
	}

	switch len(blocks) {
	case 1:
		return r.renderSingleConditionBlock(&blocks[0])

	case 2:
		// Possibly an if/else if they share the same .Cond but differ in Negate
		a, b := blocks[0], blocks[1]
		if a.Cond == b.Cond &&
			((!a.Negate && b.Negate) || (a.Negate && !b.Negate)) {
			// Same condition, one if, one unless; single ternary
			if a.Negate {
				// Switch them so the positive is first
				return r.renderTwoWayTernaryBlock(b, a)
			}
			return r.renderTwoWayTernaryBlock(a, b)
		} else {
			// Different conditions; just two separate blocks
			if err := r.renderSingleConditionBlock(&a); err != nil {
				return err
			}
			return r.renderSingleConditionBlock(&b)
		}

	default:
		// More than 2; chain them into nested ternaries
		return r.renderMultiTernaryBlocks(blocks)
	}
}

// renderSingleConditionBlock renders a conditional JSX expression.
// Example:
//
//	{ props.foo && ( ... ) }
//
// or
//
//	{ !props.foo && ( ... ) }
func (r *renderer) renderSingleConditionBlock(b *conditionalBlock) error {
	if err := r.writeIndent(); err != nil {
		return err
	}
	if err := r.print("{"); err != nil {
		return err
	}
	if b.Negate {
		if err := r.print("!"); err != nil {
			return err
		}
	}
	if err := r.print("props."); err != nil {
		return err
	}
	if err := r.print(b.Cond); err != nil {
		return err
	}
	if err := r.print(" && (\n"); err != nil {
		return err
	}

	r.indent++
	for c := b.Body.FirstChild; c != nil; c = c.NextSibling {
		if err := r.render(c); err != nil {
			return err
		}
	}
	r.indent--

	if err := r.writeIndent(); err != nil {
		return err
	}
	if err := r.print(")}\n"); err != nil {
		return err
	}

	return nil
}

// renderTwoWayTernaryBlock renders a two-way ternary JSX expression.
// Example:
//
//	{ props.foo ? ( ... ) : ( ... ) }
func (r *renderer) renderTwoWayTernaryBlock(a, b conditionalBlock) error {
	if err := r.writeIndent(); err != nil {
		return err
	}
	if err := r.print("{"); err != nil {
		return err
	}
	if a.Negate {
		if err := r.print("!"); err != nil {
			return err
		}
	}
	if err := r.print("props."); err != nil {
		return err
	}
	if err := r.print(a.Cond); err != nil {
		return err
	}
	if err := r.print(" ? (\n"); err != nil {
		return err
	}

	r.indent++
	for c := a.Body.FirstChild; c != nil; c = c.NextSibling {
		if err := r.render(c); err != nil {
			return err
		}
	}
	r.indent--

	if err := r.writeIndent(); err != nil {
		return err
	}
	if err := r.print(") : (\n"); err != nil {
		return err
	}

	r.indent++
	for c := b.Body.FirstChild; c != nil; c = c.NextSibling {
		if err := r.render(c); err != nil {
			return err
		}
	}
	r.indent--

	if err := r.writeIndent(); err != nil {
		return err
	}
	if err := r.print(")}\n"); err != nil {
		return err
	}

	return nil
}

// renderMultiTernaryBlocks renders a multi-way ternary JSX expression.
// Example:
//
//	{ props.cond1 ? (...) : props.cond2 ? (...) : ... : null }
func (r *renderer) renderMultiTernaryBlocks(blocks []conditionalBlock) error {
	for i, b := range blocks {
		if i == 0 {
			if err := r.writeIndent(); err != nil {
				return err
			}
			if err := r.print("{"); err != nil {
				return err
			}
		}
		if b.Negate {
			if err := r.print("!props."); err != nil {
				return err
			}
		} else {
			if err := r.print("props."); err != nil {
				return err
			}
		}
		if err := r.print(b.Cond); err != nil {
			return err
		}
		if err := r.print(" ? (\n"); err != nil {
			return err
		}

		r.indent++
		for c := b.Body.FirstChild; c != nil; c = c.NextSibling {
			if err := r.render(c); err != nil {
				return err
			}
		}
		r.indent--

		if err := r.writeIndent(); err != nil {
			return err
		}
		if err := r.print(") : "); err != nil {
			return err
		}
	}
	// final fallback
	if err := r.print("null}\n"); err != nil {
		return err
	}
	return nil
}
