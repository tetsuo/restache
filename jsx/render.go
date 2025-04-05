package jsx

import (
	"bytes"
	"io"

	"github.com/tetsuo/stache"
)

type Renderer struct {
	w        io.Writer
	indent   int
	doc      *stache.Node // top-level component node
	cur      *stache.Node // pointer to current child
	started  bool         // fragment open
	finished bool         // fragmend closed
}

func NewRenderer(w io.Writer, indent int, doc *stache.Node) *Renderer {
	return &Renderer{
		w:      w,
		indent: indent,
		doc:    doc,
	}
}

func (r *Renderer) RenderNext() bool {
	if r.doc == nil || r.finished {
		return false
	}

	if !r.started {
		// First call; open fragment if needed
		if numChildren(r.doc) > 1 {
			r.writeIndent()
			r.w.Write([]byte("<>\n"))
			r.indent++
		}
		r.cur = r.doc.FirstChild
		r.started = true
	}

	// Done iterating
	if r.cur == nil {
		if numChildren(r.doc) > 1 {
			r.indent--
			r.writeIndent()
			r.w.Write([]byte("</>\n"))
		}
		r.finished = true
		return false
	}

	// Render this node
	node := r.cur
	r.cur = r.cur.NextSibling // advance pointer

	// Handle conditional group logic
	if node.Type == stache.WhenNode || node.Type == stache.UnlessNode {
		group := []*stache.Node{}
		cond := node.Data
		for node != nil && (node.Type == stache.WhenNode || node.Type == stache.UnlessNode) && bytes.Equal(node.Data, cond) {
			group = append(group, node)
			node = node.NextSibling
			r.cur = node // update current after grouped block
		}
		r.renderConditionalGroup(group)
		return true
	}

	// Otherwise render normally
	r.renderNode(node)
	return true
}

// conditionalBlock is a small struct for when/unless blocks in a group.
type conditionalBlock struct {
	Cond   []byte
	Negate bool
	Body   *stache.Node
}

// renderSingleConditionBlock renders a conditional JSX expression.
// Example:
//
//	{ props.foo && ( ... ) }
//
// or
//
//	{ !props.foo && ( ... ) }
func (r *Renderer) renderSingleConditionBlock(b *conditionalBlock) {
	r.writeIndent()
	r.w.Write([]byte("{"))
	if b.Negate {
		r.w.Write([]byte("!"))
	}
	r.w.Write([]byte("props."))
	r.w.Write(b.Cond)
	r.w.Write([]byte(" && (\n"))

	r.indent++
	r.renderChildren(b.Body)
	r.indent--

	r.writeIndent()
	r.w.Write([]byte(")}\n"))
}

// renderTwoWayTernaryBlock renders a two-way ternary JSX expression.
// Example:
//
//	{ props.foo ? ( ... ) : ( ... ) }
func (r *Renderer) renderTwoWayTernaryBlock(a, b conditionalBlock) {
	r.writeIndent()
	r.w.Write([]byte("{"))
	if a.Negate {
		r.w.Write([]byte("!"))
	}
	r.w.Write([]byte("props."))
	r.w.Write(a.Cond)
	r.w.Write([]byte(" ? (\n"))

	r.indent++
	r.renderChildren(a.Body)
	r.indent--

	r.writeIndent()
	r.w.Write([]byte(") : (\n"))

	r.indent++
	r.renderChildren(b.Body)
	r.indent--

	r.writeIndent()
	r.w.Write([]byte(")}\n"))
}

// renderMultiTernaryBlocks renders a multi-way ternary JSX expression.
// Example:
//
//	{ props.cond1 ? (...) : props.cond2 ? (...) : ... : null }
func (r *Renderer) renderMultiTernaryBlocks(blocks []conditionalBlock) {
	for i, b := range blocks {
		if i == 0 {
			r.writeIndent()
			r.w.Write([]byte("{"))
		}

		if b.Negate {
			r.w.Write([]byte("!props."))
		} else {
			r.w.Write([]byte("props."))
		}
		r.w.Write(b.Cond)
		r.w.Write([]byte(" ? (\n"))

		r.indent++
		r.renderChildren(b.Body)
		r.indent--

		r.writeIndent()
		r.w.Write([]byte(") : "))
	}
	// final fallback
	r.w.Write([]byte("null}\n"))
}

// renderConditionalGroup renders a conditional JSX expression based on the node group.
// The function selects the appropriate rendering method (single, ternary, or multi-ternary).
func (r *Renderer) renderConditionalGroup(group []*stache.Node) {
	// Convert each node into a conditionalBlock
	blocks := make([]conditionalBlock, len(group))
	for i, n := range group {
		blocks[i] = conditionalBlock{
			Cond:   n.Data,
			Negate: (n.Type == stache.UnlessNode),
			Body:   n,
		}
	}

	switch len(blocks) {
	case 1:
		r.renderSingleConditionBlock(&blocks[0])

	case 2:
		// Possibly an if/else if they share the same .Cond but differ in Negate
		a, b := blocks[0], blocks[1]
		if bytes.Equal(a.Cond, b.Cond) &&
			((!a.Negate && b.Negate) || (a.Negate && !b.Negate)) {
			// Same condition, one if, one unless; single ternary
			if a.Negate {
				// Switch them so the positive is first
				r.renderTwoWayTernaryBlock(b, a)
			} else {
				r.renderTwoWayTernaryBlock(a, b)
			}
		} else {
			// Different conditions; just two separate blocks
			r.renderSingleConditionBlock(&a)
			r.renderSingleConditionBlock(&b)
		}

	default:
		// More than 2; chain them into nested ternaries
		r.renderMultiTernaryBlocks(blocks)
	}
}

// renderComponentNode renders a top-level component node, wrapping children in a React
// fragment (<>...</>) if multiple children are present.
func (r *Renderer) renderComponentNode(n *stache.Node) {
	needsFragment := numChildren(n) > 1
	if needsFragment {
		r.writeIndent()
		r.w.Write([]byte("<>\n"))
		r.indent++
	}

	for c := n.FirstChild; c != nil; {
		// Check if we have a consecutive group of conditionals with the same .Data
		if c.Type == stache.WhenNode || c.Type == stache.UnlessNode {
			group := []*stache.Node{}
			cond := c.Data

			for c != nil {
				// Skip pure-whitespace text nodes
				if c.Type == stache.TextNode && isAllWhitespace(c.Data) {
					c = c.NextSibling
					continue
				}
				// If next is same condition; group it
				if (c.Type == stache.WhenNode || c.Type == stache.UnlessNode) && bytes.Equal(c.Data, cond) {
					group = append(group, c)
					c = c.NextSibling
				} else {
					break
				}
			}
			r.renderConditionalGroup(group)
			continue
		}

		// Otherwise just render normally
		r.renderNode(c)
		c = c.NextSibling
	}

	if needsFragment {
		r.indent--
		r.writeIndent()
		r.w.Write([]byte("</>\n"))
	}
}

// renderElementNode renders a normal element.
func (r *Renderer) renderElementNode(n *stache.Node) {
	r.writeIndent()
	r.w.Write([]byte("<"))
	r.w.Write(n.Data)
	// TODO: render attrs
	if n.FirstChild == nil {
		r.w.Write([]byte(" />\n"))
		return
	}

	r.w.Write([]byte(">\n"))
	r.indent++

	for c := n.FirstChild; c != nil; {
		if c.Type == stache.WhenNode || c.Type == stache.UnlessNode {
			group := []*stache.Node{}
			cond := c.Data
			for c != nil && (c.Type == stache.WhenNode || c.Type == stache.UnlessNode) && bytes.Equal(c.Data, cond) {
				group = append(group, c)
				c = c.NextSibling
			}
			r.renderConditionalGroup(group)
		} else {
			r.renderNode(c)
			c = c.NextSibling
		}
	}

	r.indent--
	r.writeIndent()
	r.w.Write([]byte("</"))
	r.w.Write(n.Data)
	r.w.Write([]byte(">\n"))
}

// renderConditionalNode renders a single When or Unless conditional JSX expression.
// Example:
//
//	{ props.foo && ( ... ) }
//	{ !props.foo && ( ... ) }
func (r *Renderer) renderConditionalNode(n *stache.Node) {
	r.writeIndent()
	r.w.Write([]byte("{"))
	if n.Type == stache.UnlessNode {
		r.w.Write([]byte("!"))
	}
	r.w.Write([]byte("props."))
	r.w.Write(n.Data)
	r.w.Write([]byte(" && (\n"))

	r.indent++
	r.renderChildren(n)
	r.indent--

	r.writeIndent()
	r.w.Write([]byte(")}\n"))
}

// renderTextNode renders a text JSX expression.
// Example:
//
//	"Hello, world!"
func (r *Renderer) renderTextNode(n *stache.Node) {
	r.writeIndent()
	r.w.Write(n.Data)
	r.w.Write([]byte("\n"))
}

// renderVariableNode renders a variable JSX expression.
// Example:
//
//	{props.someValue}
func (r *Renderer) renderVariableNode(n *stache.Node) {
	r.writeIndent()
	r.w.Write([]byte("{props."))
	r.w.Write(n.Data)
	r.w.Write([]byte("}\n"))
}

// renderCommentNode renders a comment JSX expression.
// Example:
//
//	{/* comment */}
func (r *Renderer) renderCommentNode(n *stache.Node) {
	r.writeIndent()
	r.w.Write([]byte("{/* "))
	r.w.Write(n.Data)
	r.w.Write([]byte(" */}\n"))
}

// renderNode is the central switch that dispatches by node type.
func (r *Renderer) renderNode(n *stache.Node) {
	switch n.Type {
	case stache.TextNode:
		r.renderTextNode(n)
	case stache.VariableNode:
		r.renderVariableNode(n)
	case stache.ElementNode:
		r.renderElementNode(n)
	case stache.WhenNode, stache.UnlessNode:
		r.renderConditionalNode(n)
	case stache.ComponentNode:
		r.renderComponentNode(n)
	case stache.CommentNode:
		r.renderCommentNode(n)
	case stache.RangeNode:
		// TODO: Range/loop not implemented yet
	default:
		// No-op
	}
}

// renderChildren just calls renderNode on all children in order.
func (r *Renderer) renderChildren(n *stache.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.renderNode(c)
	}
}

// writeIndent writes "  " * indentLevel.
func (r *Renderer) writeIndent() {
	for range r.indent {
		r.w.Write([]byte("  "))
	}
}

// isAllWhitespace checks if b consists only of whitespace.
func isAllWhitespace(b []byte) bool {
	for _, c := range b {
		switch c {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return false
		}
	}
	return true
}

// numChildren counts the number of direct children of a node.
func numChildren(n *stache.Node) int {
	i := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		i++
	}
	return i
}
