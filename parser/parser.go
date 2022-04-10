package parser

import (
	"errors"
	"fmt"
	"io"

	"github.com/onur1/stache/lexer"
)

type NodeKind int

const (
	Tag NodeKind = iota + 1
	Section
	Variable
	InvertedSection
	Comment
	Text
)

type TextNode struct {
	Text string
}

func (c TextNode) Kind() NodeKind {
	return Text
}

func (c TextNode) Serialize() interface{} {
	return c.Text
}

type VariableNode struct {
	Name string
}

func (c VariableNode) Kind() NodeKind {
	return Variable
}

func (c VariableNode) Serialize() interface{} {
	return []interface{}{
		Variable,
		c.Name,
	}
}

type CommentNode struct {
	Comment string
}

func (c CommentNode) Kind() NodeKind {
	return Comment
}

func (c CommentNode) Serialize() interface{} {
	return []interface{}{
		Comment,
		c.Comment,
	}
}

type SectionNode struct {
	Name     string
	Inverted bool
	Children []Node
}

func (c SectionNode) Kind() NodeKind {
	return Section
}

func (c SectionNode) Serialize() interface{} {
	var kind NodeKind
	if c.Inverted {
		kind = InvertedSection
	} else {
		kind = Section
	}
	children := make([]interface{}, len(c.Children))
	for i, v := range c.Children {
		children[i] = v.Serialize()
	}
	return []interface{}{kind, c.Name, children}
}

type TagNode struct {
	Name     string
	Children []Node
	Attrs    map[string][]Node
}

func (c TagNode) Kind() NodeKind {
	return Tag
}

func (c TagNode) Serialize() interface{} {
	children := make([]interface{}, len(c.Children))
	for i, v := range c.Children {
		children[i] = v.Serialize()
	}
	attrs := make(map[string]interface{}, len(c.Attrs))
	for k, list := range c.Attrs {
		xs := make([]interface{}, len(list))
		for i, v := range list {
			xs[i] = v.Serialize()
		}
		attrs[k] = xs
	}
	return []interface{}{c.Name, attrs, children}
}

type Node interface {
	Kind() NodeKind
	Serialize() interface{}
}

var (
	ErrSyntax = errors.New("stache: syntax error")
)

type parser struct {
	stack [][]Node
	err   error
}

func parse(p *parser, cb func(Node) bool) func(lexer.Token) bool {
	return func(t lexer.Token) bool {
		switch t.Kind {
		case lexer.Open:
			children := []Node{}
			p.stack = append(p.stack, children)
			i := len(p.stack) - 2
			node := &TagNode{Name: t.Body}
			p.stack[i] = append(p.stack[i], node)
			attrs := make(map[string][]Node, len(t.Attrs))
			for key, tks := range t.Attrs {
				attrParser := new(parser)
				attrParser.stack = [][]Node{children}
				next := parse(attrParser, nil)
				var tk lexer.Token
				for len(tks) > 0 {
					tk, tks = tks[0], tks[1:]
					if ok := next(tk); !ok {
						p.err = attrParser.err
						return false
					}
				}
				attrs[key] = attrParser.stack[0]
			}
			if len(attrs) > 0 {
				node.Attrs = attrs
			}
		case lexer.Close:
			i := len(p.stack) - 2
			if i < 0 {
				p.err = fmt.Errorf("%w: tag not initialized: %s", ErrSyntax, t.Body)
				return false
			}
			node := p.stack[i][len(p.stack[i])-1]
			if treeNode, ok := node.(*TagNode); ok {
				if treeNode.Name != t.Body {
					p.err = fmt.Errorf("%w: <%s>...</%s>", ErrSyntax, treeNode.Name, t.Body)
				} else {
					children := p.stack[i+1]
					if len(children) > 0 {
						treeNode.Children = p.stack[i+1]
					}
					p.stack = p.stack[:i+1]
					if len(p.stack) == 1 && cb != nil {
						cb(node)
					}
					break
				}
			}
			return false
		case lexer.SectionOpen, lexer.InvertedSectionOpen:
			p.stack = append(p.stack, make([]Node, 0))
			i := len(p.stack) - 2
			p.stack[i] = append(p.stack[i], &SectionNode{
				Name:     t.Body,
				Inverted: t.Kind == lexer.InvertedSectionOpen,
			})
		case lexer.SectionClose:
			i := len(p.stack) - 2
			if i < 0 {
				p.err = fmt.Errorf("%w: section not initialized: %s", ErrSyntax, t.Body)
				return false
			}
			node := p.stack[i][len(p.stack[i])-1]
			if sectionNode, ok := node.(*SectionNode); ok {
				if sectionNode.Name != t.Body {
					p.err = fmt.Errorf("%w: {#%s}...{/%s}", ErrSyntax, sectionNode.Name, t.Body)
				} else {
					children := p.stack[i+1]
					if len(children) > 0 {
						sectionNode.Children = p.stack[i+1]
					}
					p.stack = p.stack[:i+1]
					if len(p.stack) == 1 && cb != nil {
						cb(node)
					}
					break
				}
			}
			return false
		case lexer.Text:
			i := len(p.stack) - 1
			node := &TextNode{
				Text: t.Body,
			}
			if i == 0 && cb != nil {
				cb(node)
			}
			p.stack[i] = append(p.stack[i], node)
		case lexer.Variable:
			i := len(p.stack) - 1
			node := &VariableNode{
				Name: t.Body,
			}
			if i == 0 && cb != nil {
				cb(node)
			}
			p.stack[i] = append(p.stack[i], node)
		case lexer.Comment:
			i := len(p.stack) - 1
			node := &CommentNode{
				Comment: t.Body,
			}
			if i == 0 && cb != nil {
				cb(node)
			}
			p.stack[i] = append(p.stack[i], node)
		}
		return true
	}
}

func Parse(r io.Reader, cb func(Node) bool) error {
	p := new(parser)
	p.stack = [][]Node{{}}
	if err := lexer.Tokenize(r, parse(p, cb)); err != nil {
		return err
	}
	p.stack = nil
	return p.err
}
