package restache

import (
	"slices"

	"golang.org/x/net/html/atom"
)

type NodeType uint32

const (
	ErrorNode NodeType = iota
	TextNode
	ComponentNode
	ElementNode
	CommentNode
	VariableNode
	RangeNode
	WhenNode
	UnlessNode
)

type Attribute struct {
	Key     string
	KeyAtom atom.Atom
	Val     string
	IsExpr  bool
}

type PathComponent struct {
	Key     string
	IsRange bool
}

type Node struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	Type     NodeType
	DataAtom atom.Atom
	Data     string
	Attr     []Attribute
	Path     []PathComponent
}

func (n *Node) TagName() string {
	if n.DataAtom != 0 {
		return n.DataAtom.String()
	}
	return n.Data
}

// InsertBefore inserts newChild as a child of n, immediately before oldChild
// in the sequence of n's children. oldChild may be nil, in which case newChild
// is appended to the end of n's children.
//
// It will panic if newChild already has a parent or siblings.
func (n *Node) InsertBefore(newChild, oldChild *Node) {
	if newChild.Parent != nil || newChild.PrevSibling != nil || newChild.NextSibling != nil {
		panic("restache: InsertBefore called for an attached child Node")
	}
	var prev, next *Node
	if oldChild != nil {
		prev, next = oldChild.PrevSibling, oldChild
	} else {
		prev = n.LastChild
	}
	if prev != nil {
		prev.NextSibling = newChild
	} else {
		n.FirstChild = newChild
	}
	if next != nil {
		next.PrevSibling = newChild
	} else {
		n.LastChild = newChild
	}
	newChild.Parent = n
	newChild.PrevSibling = prev
	newChild.NextSibling = next
}

func (n *Node) Render(w writer) (int, error) {
	return Render(w, n)
}

// AppendChild adds a node c as a child of n.
//
// It will panic if c already has a parent or siblings.
func (n *Node) AppendChild(c *Node) {
	if c.Parent != nil || c.PrevSibling != nil || c.NextSibling != nil {
		panic("restache: AppendChild called for an attached child Node")
	}
	last := n.LastChild
	if last != nil {
		last.NextSibling = c
	} else {
		n.FirstChild = c
	}
	n.LastChild = c
	c.Parent = n
	c.PrevSibling = last
}

// RemoveChild removes a node c that is a child of n. Afterwards, c will have
// no parent and no siblings.
//
// It will panic if c's parent is not n.
func (n *Node) RemoveChild(c *Node) {
	if c.Parent != n {
		panic("restache: RemoveChild called for a non-child Node")
	}
	if n.FirstChild == c {
		n.FirstChild = c.NextSibling
	}
	if c.NextSibling != nil {
		c.NextSibling.PrevSibling = c.PrevSibling
	}
	if n.LastChild == c {
		n.LastChild = c.PrevSibling
	}
	if c.PrevSibling != nil {
		c.PrevSibling.NextSibling = c.NextSibling
	}
	c.Parent = nil
	c.PrevSibling = nil
	c.NextSibling = nil
}

func (n *Node) wrapChildrenInFragment() {
	first := n.FirstChild
	if first == nil {
		// range with empty body; <></>
		n.AppendChild(&Node{Type: ElementNode})
		return
	}

	// range with exactly one element child; prepend key attr
	if n.Type == RangeNode &&
		first.NextSibling == nil &&
		first.Type == ElementNode {
		first.Attr = append([]Attribute{{
			Key:    "key",
			Val:    "key",
			IsExpr: true,
		}}, first.Attr...)
		return
	}

	if first.NextSibling == nil && !(first.Type == TextNode || first.Type == CommentNode) {
		return // already a single non-text node; no fragment needed
	}

	// for other scenarios, add fragment
	frag := &Node{
		Type: ElementNode,
		Path: slices.Clone(n.Path),
	}
	if n.Type == RangeNode {
		frag.Data = "React.Fragment"
		frag.Attr = []Attribute{{
			Key:    "key",
			Val:    "key",
			IsExpr: true,
		}}
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		n.RemoveChild(c)
		frag.AppendChild(c)
		c = next
	}
	n.AppendChild(frag)
}

func (n *Node) nameEquals(name []byte) bool {
	if len(n.Data) != len(name) {
		return false
	}
	for i := range n.Data {
		if n.Data[i] != name[i] {
			return false
		}
	}
	return true
}

// extractUnknownElementTags returns the .Data of every ElementNode whose DataAtom == 0,
// without duplicates, in depth-first (pre-order) order.
func (n *Node) extractUnknownElementTags() []string {
	if n == nil {
		return nil
	}

	seen := make(map[string]struct{}, 16) // seen .Data values
	out := make([]string, 0, 8)

	stack := []*Node{n.FirstChild}

	for len(stack) > 0 {
		i := len(stack) - 1
		c := stack[i]
		stack = stack[:i]

		for c != nil {
			if c.Type == ElementNode && c.DataAtom == 0 {
				if _, ok := seen[c.Data]; !ok {
					seen[c.Data] = struct{}{}
					out = append(out, c.Data)
				}
			}
			if nx := c.NextSibling; nx != nil {
				stack = append(stack, nx)
			}
			c = c.FirstChild
		}
	}
	return out
}

func (n *Node) renameUnknownElementTags(rewrites map[string]string) {
	if n == nil {
		return
	}

	var stack nodeStack
	if n.FirstChild != nil {
		stack = append(stack, n.FirstChild)
	}

	for len(stack) > 0 {
		c := stack.pop()

		for c != nil {
			if c.Type == ElementNode && c.DataAtom == 0 {
				if newVal, ok := rewrites[c.Data]; ok {
					c.Data = newVal
				}
			}
			if next := c.NextSibling; next != nil {
				stack = append(stack, next)
			}
			c = c.FirstChild
		}
	}
}

// nodeStack is a stack of nodes.
type nodeStack []*Node

// pop pops the stack. It will panic if s is empty.
func (s *nodeStack) pop() *Node {
	i := len(*s)
	n := (*s)[i-1]
	*s = (*s)[:i-1]
	return n
}

// top returns the most recently pushed node, or nil if s is empty.
func (s *nodeStack) top() *Node {
	if i := len(*s); i > 0 {
		return (*s)[i-1]
	}
	return nil
}

func (s *nodeStack) popUntilAtom(a atom.Atom) bool {
	for i := len(*s) - 1; i >= 0; i-- {
		n := (*s)[i]
		if n.Type == ElementNode && n.DataAtom == a {
			*s = (*s)[:i]
			return true
		}
	}
	return false
}

func (s *nodeStack) popUntilName(name []byte) bool {
	for i := len(*s) - 1; i >= 0; i-- {
		n := (*s)[i]
		if n.Type == ElementNode && n.nameEquals(name) {
			*s = (*s)[:i]
			return true
		}
	}
	return false
}

func (s *nodeStack) popControl(name []byte) (*Node, bool) {
	for len(*s) > 1 {
		n := s.pop()
		if (n.Type == RangeNode || n.Type == WhenNode || n.Type == UnlessNode) &&
			n.nameEquals(name) {
			return n, true
		}
	}
	return nil, false
}
