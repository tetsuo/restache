package stache

import (
	"bytes"
	"io"
	"slices"

	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"golang.org/x/net/html/atom"
)

var voidAtoms = map[atom.Atom]bool{
	atom.Area:   true,
	atom.Br:     true,
	atom.Embed:  true,
	atom.Img:    true,
	atom.Input:  true,
	atom.Wbr:    true,
	atom.Col:    true,
	atom.Hr:     true,
	atom.Link:   true,
	atom.Track:  true,
	atom.Source: true,
} // meta, base, param, keygen not included

type insertionMode func(*parser) bool

type PathSegment struct {
	Key     []byte
	IsRange bool
}

type parser struct {
	z    *Tokenizer
	oe   nodeStack
	doc  *Node // root node
	im   insertionMode
	tt   TokenType
	path []PathSegment
	sc   bool // has self closing token

	depsTrie  *iradix.Tree[int]
	dependsOn []int
}

func initialIM(p *parser) bool {
	p.oe = append(p.oe, p.doc)
	p.im = inBodyIM
	return p.im(p)
}

func inBodyIM(p *parser) bool {
	switch p.tt {
	case TextToken:
		text := bytes.Clone(bytes.TrimSpace(p.z.Raw()))
		if len(text) == 0 {
			return true
		}
		n := &Node{
			Type: TextNode,
			Data: text,
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

		for hasAttr {
			key, val, isExpr, more := p.z.TagAttr()
			elem.Attr = append(elem.Attr, Attribute{
				Key:    bytes.Clone(key),
				Val:    bytes.Clone(val),
				IsExpr: isExpr,
			})
			hasAttr = more
		}

		if p.depsTrie != nil {
			if idx, ok := p.depsTrie.Root().Get(name); ok {
				p.dependsOn = append(p.dependsOn, idx)
			}
		}

		p.oe.top().AppendChild(elem)

		if p.sc || voidAtoms[elem.DataAtom] {
			p.sc = false
			return true
		}

		p.oe = append(p.oe, elem)
		return true

	case EndTagToken:
		name, _ := p.z.TagName()
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
		if n, found = p.oe.popCtrl(name); found && n.Type == RangeNode {
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

func Parse(r io.Reader) (nodes *Node, err error) {
	nodes, _, err = parseWithDependencies(r, nil)
	return
}

func parseWithDependencies(r io.Reader, componentTrie *iradix.Tree[int]) (*Node, []int, error) {
	p := &parser{
		z:  NewTokenizer(r),
		im: initialIM,
		doc: &Node{
			Type: ComponentNode,
		},
		depsTrie: componentTrie,
	}
	if err := p.parse(); err != nil {
		return nil, nil, err
	}
	return p.doc, p.dependsOn, nil
}
