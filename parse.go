package restache

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/net/html/atom"
)

type insertionMode func(*parser) bool

type parser struct {
	z    *Tokenizer
	oe   nodeStack
	doc  *Node
	im   insertionMode
	tt   TokenType
	path []PathSegment
	sc   bool // indicates self closing token

	lookup map[string]int // dependency lookup table
	afters map[int]bool   // marked dependency indexes
}

func newParser(r io.Reader, lookup map[string]int) *parser {
	p := &parser{
		z:      NewTokenizer(r),
		im:     initialIM,
		lookup: lookup,
		doc: &Node{
			Type: ComponentNode,
		},
	}
	if p.lookup != nil {
		p.afters = make(map[int]bool)
	}
	return p
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
			Data: string(raw),
		}
		p.oe.top().AppendChild(n)
		return true

	case StartTagToken:
		name, hasAttr := p.z.TagName()

		elem := &Node{
			Type:     ElementNode,
			Data:     string(name),
			DataAtom: atom.Lookup(name),
			Path:     slices.Clone(p.path),
		}
		// Gather attributes
		for hasAttr {
			key, val, isExpr, more := p.z.TagAttr()
			elem.Attr = append(elem.Attr, Attribute{
				Key:    string(key),
				Val:    string(val),
				IsExpr: isExpr,
			})
			hasAttr = more
		}

		p.oe.top().AppendChild(elem)

		if elem.DataAtom == 0 && p.lookup != nil {
			if idx, ok := p.lookup[elem.Data]; ok {
				if _, ok = p.afters[idx]; !ok {
					p.afters[idx] = true // Mark dependency
				}
			}
		}

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
				Data: string(bytes.TrimSpace(p.z.Raw())),
				Path: slices.Clone(p.path),
			},
		)
		return true

	case WhenToken:
		node := &Node{
			Type: WhenNode,
			Data: string(bytes.TrimSpace(p.z.ControlName())),
			Path: slices.Clone(p.path),
		}
		p.oe.top().AppendChild(node)
		p.oe = append(p.oe, node)
		return true

	case UnlessToken:
		node := &Node{
			Type: UnlessNode,
			Data: string(bytes.TrimSpace(p.z.ControlName())),
			Path: slices.Clone(p.path),
		}
		p.oe.top().AppendChild(node)
		p.oe = append(p.oe, node)
		return true

	case RangeToken:
		node := &Node{
			Type: RangeNode,
			Data: string(bytes.TrimSpace(p.z.ControlName())),
			Path: slices.Clone(p.path),
		}
		parts := strings.Split(node.Data, ".")
		var (
			i    int
			part string
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
				Data: string(bytes.TrimSpace(p.z.Comment())),
			},
		)
		return true
	}

	panic(fmt.Sprintf("unknown token type: %d", int(p.tt)))
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

// voidElements only have a start tag; ends tags are not specified.
var voidElements = map[atom.Atom]bool{
	atom.Area: true, atom.Br: true, atom.Embed: true, atom.Img: true,
	atom.Input: true, atom.Wbr: true, atom.Col: true, atom.Hr: true,
	atom.Link: true, atom.Track: true, atom.Source: true,
}
