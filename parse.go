package stache

import (
	"bytes"
	"io"
	"slices"

	"golang.org/x/net/html/atom"
)

// voidElements only have a start tag; ends tags are not specified.
var voidElements = map[atom.Atom]bool{
	atom.Area: true, atom.Br: true, atom.Embed: true, atom.Img: true,
	atom.Input: true, atom.Wbr: true, atom.Col: true, atom.Hr: true,
	atom.Link: true, atom.Track: true, atom.Source: true,
}

type insertionMode func(*parser) bool

type DependencyResolver interface {
	Get([]byte) (int, bool)
}

type parser struct {
	z    *Tokenizer
	oe   nodeStack
	doc  *Node
	im   insertionMode
	tt   TokenType
	path []PathSegment
	sc   bool  // indicates self closing token
	deps []int // found deps
	dt   DependencyResolver
}

// initialIM is the first insertion mode used.
// It switches to inBodyIM after adding the root document.
func initialIM(p *parser) bool {
	p.oe = append(p.oe, p.doc)
	p.im = inBodyIM
	return p.im(p)
}

func inBodyIM(p *parser) bool {
	switch p.tt {
	case TextToken:
		raw := bytes.TrimSpace(p.z.Raw())
		if len(raw) == 0 {
			return true
		}
		n := &Node{
			Type: TextNode,
			Data: bytes.Clone(raw),
		}
		p.oe.top().AppendChild(n)
		return true

	case StartTagToken:
		name, hasAttr := p.z.TagName()

		elem := &Node{
			Type:     ElementNode,
			Data:     name,
			DataAtom: atom.Lookup(name),
			Path:     slices.Clone(p.path),
		}
		// Gather attributes
		for hasAttr {
			key, val, isExpr, more := p.z.TagAttr()
			elem.Attr = append(elem.Attr, Attribute{
				Key:    bytes.Clone(key),
				Val:    bytes.Clone(val),
				IsExpr: isExpr,
			})
			hasAttr = more
		}
		// If a trie of dependencies available, see if tag name is in it:
		if p.dt != nil {
			if idx, ok := p.dt.Get(name); ok {
				// add to dependencies
				p.deps = append(p.deps, idx)
			}
		}
		p.oe.top().AppendChild(elem)

		// If it's self-closing tag, or void element, don't push onto the stack:
		if p.sc || voidElements[elem.DataAtom] {
			p.sc = false
			return true
		}
		// else push onto stack
		p.oe = append(p.oe, elem)
		return true

	case EndTagToken:
		name, _ := p.z.TagName()
		// pop stack until a matching element is found
		p.oe.popUntil(atom.Lookup(name), name)
		return true

	case VariableToken:
		p.oe.top().AppendChild(
			&Node{
				Type: VariableNode,
				Data: bytes.Clone(bytes.TrimSpace(p.z.Raw())),
				Path: slices.Clone(p.path),
			},
		)
		return true

	case WhenToken:
		node := &Node{
			Type: WhenNode,
			Data: bytes.Clone(bytes.TrimSpace(p.z.ControlName())),
			Path: slices.Clone(p.path),
		}
		p.oe.top().AppendChild(node)
		p.oe = append(p.oe, node)
		return true

	case UnlessToken:
		node := &Node{
			Type: UnlessNode,
			Data: bytes.Clone(bytes.TrimSpace(p.z.ControlName())),
			Path: slices.Clone(p.path),
		}
		p.oe.top().AppendChild(node)
		p.oe = append(p.oe, node)
		return true

	case RangeToken:
		node := &Node{
			Type: RangeNode,
			Data: bytes.Clone(bytes.TrimSpace(p.z.ControlName())),
			Path: slices.Clone(p.path),
		}
		parts := bytes.Split(node.Data, []byte("."))
		var (
			i    int
			part []byte
			last = len(parts) - 1
		)
		for i, part = range parts {
			p.path = append(p.path, PathSegment{
				Key:     part,
				IsRange: i == last,
			})
		}
		p.oe.top().AppendChild(node)
		p.oe = append(p.oe, node)
		return true

	case EndControlToken:
		var (
			name  = bytes.TrimSpace(p.z.ControlName())
			n     *Node
			found bool
		)
		// If it's a range node, restore the path
		if n, found = p.oe.popControl(name); found && n.Type == RangeNode {
			p.path = n.Path
		}
		return true

	case CommentToken:
		p.oe.top().AppendChild(
			&Node{
				Type: CommentNode,
				Data: bytes.Clone(bytes.TrimSpace(p.z.Comment())),
			},
		)
		return true
	}

	return false
}

// parseCurrentTokens call the current insertion mode; it sets the self-closing-tag
// mode (p.sc) on, if a SelfClosingTagToken is encountered.
func (p *parser) parseCurrentToken() {
	if p.tt == SelfClosingTagToken {
		p.sc = true
		p.tt = StartTagToken
	}

	consumed := false
	for !consumed {
		consumed = p.im(p)
	}
}

func (p *parser) parse() error {
	var err error
	for err != io.EOF {
		p.tt = p.z.Next()
		if p.tt == ErrorToken {
			if err := p.z.Err(); err != nil && err != io.EOF {
				return err
			}
			break
		}
		p.parseCurrentToken()
	}
	return nil
}

// Parse parses a single Node with no dependencies.
func Parse(r io.Reader) (nodes *Node, err error) {
	nodes, _, err = ParseWithDependencies(r, nil)
	return
}

// ParseWithDependencies also tracks references to component indexes in the given trie.
func ParseWithDependencies(r io.Reader, resolver DependencyResolver) (*Node, []int, error) {
	p := &parser{
		z:  NewTokenizer(r),
		im: initialIM,
		doc: &Node{
			Type: ComponentNode,
		},
		dt: resolver,
	}
	if err := p.parse(); err != nil {
		return nil, nil, err
	}
	return p.doc, p.deps, nil
}
