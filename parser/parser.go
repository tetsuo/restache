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

type VariableNode struct {
	Name string
}

func (c VariableNode) Kind() NodeKind {
	return Variable
}

type CommentNode struct {
	Comment string
}

func (c CommentNode) Kind() NodeKind {
	return Comment
}

type SectionNode struct {
	Name     string
	Inverted bool
	Children []Node
}

func (c SectionNode) Value() []Node {
	return c.Children
}

func (c SectionNode) Kind() NodeKind {
	return Section
}

type TagNode struct {
	Name     string
	Children []Node
	Attrs    map[string][]Node
}

func (c TagNode) Kind() NodeKind {
	return Tag
}

func (c TagNode) Value() []Node {
	return c.Children
}

type Node interface {
	Kind() NodeKind
}

type Tree interface {
	Value() []Node
}

type Parent interface {
	Children() []Node
}

var (
	ErrSyntax = errors.New("stache: syntax error")
)

type parser struct {
	stack [][]Node
	err   error
}

func parse(p *parser, cb func(*TagNode) bool) func(lexer.Token) bool {
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
					return false
				}
				children := p.stack[i+1]
				if len(children) > 0 {
					treeNode.Children = p.stack[i+1]
				}
				p.stack = p.stack[:i+1]
				if len(p.stack) == 1 && cb != nil {
					cb(node.(*TagNode))
				}
				break
			}
			p.err = fmt.Errorf("%w: ...<%s>", ErrSyntax, t.Body)
			return false
		case lexer.Text:
			i := len(p.stack) - 1
			p.stack[i] = append(p.stack[i], &TextNode{
				Text: t.Body,
			})
		case lexer.Variable:
			i := len(p.stack) - 1
			p.stack[i] = append(p.stack[i], &VariableNode{
				Name: t.Body,
			})
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
					return false
				}
				children := p.stack[i+1]
				if len(children) > 0 {
					sectionNode.Children = p.stack[i+1]
				}
				p.stack = p.stack[:i+1]
				break
			}
			p.err = fmt.Errorf("%w: ...{/%s}", ErrSyntax, t.Body)
			return false
		case lexer.Comment:
			i := len(p.stack) - 1
			p.stack[i] = append(p.stack[i], &CommentNode{
				Comment: t.Body,
			})
		}
		return true
	}
}

func Parse(r io.Reader, cb func(*TagNode) bool) error {
	p := new(parser)
	p.stack = [][]Node{{}}
	if err := lexer.Tokenize(r, parse(p, cb)); err != nil {
		return err
	}
	return p.err
}
