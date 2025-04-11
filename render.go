package restache

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

type renderer struct {
	w writer

	indent int
}

func (r *renderer) renderComponent(n *Node) error {
	w := r.w
	if _, err := w.WriteString("import * as React from 'react';"); err != nil {
		return err
	}
	for _, attr := range n.Attr {
		if _, err := w.WriteString(fmt.Sprintf("import %s from \"./%s.jsx\";\n", attr.Key, attr.Val)); err != nil {
			return err
		}
	}
	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	if _, err := w.WriteString("export default function"); err != nil {
		return err
	}
	if n.Data != "" {
		if err := w.WriteByte(' '); err != nil {
			return err
		}
		if _, err := w.WriteString(n.Data); err != nil {
			return err
		}
	}
	if _, err := w.WriteString("(props) {\n"); err != nil {
		return err
	}
	if _, err := w.WriteString("  return ("); err != nil {
		return err
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := r.render(c); err != nil {
			return err
		}
	}
	if _, err := w.WriteString("  );"); err != nil {
		return err
	}
	if _, err := w.WriteString("}"); err != nil {
		return err
	}
	return nil
}

func (r *renderer) writeIndent() error {
	for range r.indent {
		if _, err := r.w.WriteString("  "); err != nil {
			return err
		}
	}
	return nil
}

func (r *renderer) render(n *Node) error {
	w := r.w
	switch n.Type {
	case ErrorNode:
		return errors.New("html: cannot render an ErrorNode node")
	case TextNode:
		return escape(w, n.Data)
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

// renderElementNode renders a normal element.
func (r *renderer) renderElement(n *Node) error {
	w := r.w
	if err := r.writeIndent(); err != nil {
		return err
	}
	// Render the <xxx> opening tag.
	if err := w.WriteByte('<'); err != nil {
		return err
	}
	if _, err := w.WriteString(n.Data); err != nil {
		return err
	}
	for _, a := range n.Attr {
		if err := w.WriteByte(' '); err != nil {
			return err
		}
		if _, err := w.WriteString(a.Key); err != nil {
			return err
		}
		if _, err := w.WriteString(`="`); err != nil {
			return err
		}
		if err := escape(w, a.Val); err != nil {
			return err
		}
		if err := w.WriteByte('"'); err != nil {
			return err
		}
	}
	// TODO: check void elements
	if _, err := w.WriteString(">\n"); err != nil {
		return err
	}

	// Add initial newline where there is danger of a newline being ignored.
	if c := n.FirstChild; c != nil && c.Type == TextNode && strings.HasPrefix(c.Data, "\n") {
		switch n.Data {
		case "pre", "listing", "textarea":
			if err := w.WriteByte('\n'); err != nil {
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
	if _, err := w.WriteString("</"); err != nil {
		return err
	}
	if _, err := w.WriteString(n.Data); err != nil {
		return err
	}
	if _, err := w.WriteString(">\n"); err != nil {
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
	if _, err := r.w.WriteString("{props."); err != nil {
		return err
	}
	if _, err := r.w.WriteString(n.Data); err != nil {
		return err
	}
	if _, err := r.w.WriteString("}\n"); err != nil {
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
	if _, err := r.w.WriteString("{/* "); err != nil {
		return err
	}
	if err := escapeComment(r.w, n.Data); err != nil {
		return err
	}
	if _, err := r.w.WriteString(" */}\n"); err != nil {
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
	if _, err := r.w.WriteString("{"); err != nil {
		return err
	}
	if b.Negate {
		if _, err := r.w.WriteString("!"); err != nil {
			return err
		}
	}
	if _, err := r.w.WriteString("props."); err != nil {
		return err
	}
	if _, err := r.w.WriteString(b.Cond); err != nil {
		return err
	}
	if _, err := r.w.WriteString(" && (\n"); err != nil {
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
	if _, err := r.w.WriteString(")}\n"); err != nil {
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
	if _, err := r.w.WriteString("{"); err != nil {
		return err
	}
	if a.Negate {
		if _, err := r.w.WriteString("!"); err != nil {
			return err
		}
	}
	if _, err := r.w.WriteString("props."); err != nil {
		return err
	}
	if _, err := r.w.WriteString(a.Cond); err != nil {
		return err
	}
	if _, err := r.w.WriteString(" ? (\n"); err != nil {
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
	if _, err := r.w.WriteString(") : (\n"); err != nil {
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
	if _, err := r.w.WriteString(")}\n"); err != nil {
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
			if _, err := r.w.WriteString("{"); err != nil {
				return err
			}
		}
		if b.Negate {
			if _, err := r.w.WriteString("!props."); err != nil {
				return err
			}
		} else {
			if _, err := r.w.WriteString("props."); err != nil {
				return err
			}
		}
		if _, err := r.w.WriteString(b.Cond); err != nil {
			return err
		}
		if _, err := r.w.WriteString(" ? (\n"); err != nil {
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
		if _, err := r.w.WriteString(") : "); err != nil {
			return err
		}
	}
	// final fallback
	if _, err := r.w.WriteString("null}\n"); err != nil {
		return err
	}
	return nil
}
