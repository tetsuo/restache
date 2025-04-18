package restache

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/net/html/atom"
)

// Parse parses a single Node with no dependencies.
func Parse(r io.Reader) (node *Node, err error) {
	p := newParser(r, nil)
	if err = p.parse(); err != nil {
		return
	}
	node = p.doc
	return
}

type insertionMode func(*parser) bool

type parser struct {
	z    *Tokenizer
	oe   nodeStack
	doc  *Node
	im   insertionMode
	tt   TokenType
	path []PathComponent
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

func (p *parser) markDependency(data string) {
	if p.lookup != nil {
		if idx, ok := p.lookup[data]; ok {
			if _, ok = p.afters[idx]; !ok {
				p.afters[idx] = true
			}
		}
	}
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

		e := &Node{
			Type:     ElementNode,
			DataAtom: atom.Lookup(name),
			Path:     slices.Clone(p.path),
		}

		if e.DataAtom != 0 {
			if _, ok := commonElements[e.DataAtom]; !ok {
				e.Data = e.DataAtom.String()
				e.DataAtom = 0
				p.markDependency(e.Data)
			}
		} else {
			e.Data = string(name)
			p.markDependency(e.Data)
		}

		if hasAttr {
			if _, found := nonSpecCamelAttrTags[e.DataAtom]; found {
				searchPrefix := uint64(e.DataAtom) << 32
				for hasAttr {
					key, val, isExpr, more := p.z.TagAttr()
					x := Attribute{
						KeyAtom: atom.Lookup(key),
						Val:     string(val),
						IsExpr:  isExpr,
					}
					if x.KeyAtom == 0 {
						if attrIsNotDataOrAria(key) {
							h, camelSafe, i := fnv(hash0, key)
							if camelSafe && i < len(key)-1 {
								x.Key = string(camelize(key, i))
							} else if match, known := nonSpecCamelAttrTable[searchPrefix|uint64(h)]; known {
								x.Key = match
							} else {
								x.Key = string(key)
							}
						} else {
							x.Key = string(key)
						}
					}
					e.Attr = append(e.Attr, x)
					hasAttr = more
				}
			} else {
				for hasAttr {
					key, val, isExpr, more := p.z.TagAttr()
					x := Attribute{
						KeyAtom: atom.Lookup(key),
						Val:     string(val),
						IsExpr:  isExpr,
					}
					if x.KeyAtom == 0 {
						if attrIsNotDataOrAria(key) {
							x.Key = string(camelize(key, 0))
						} else {
							x.Key = string(key)
						}
					}
					e.Attr = append(e.Attr, x)
					hasAttr = more
				}
			}
		}

		p.oe.top().AppendChild(e)

		// If it's self-closing tag, or void element, don't push onto the stack:
		if p.sc {
			p.sc = false
			return true
		}
		if _, ok := voidElements[e.DataAtom]; ok {
			p.sc = false
			return true
		}
		// else push onto stack
		p.oe = append(p.oe, e)
		return true

	case EndTagToken:
		name, _ := p.z.TagName()
		// pop stack until a matching element is found
		a := atom.Lookup(name)
		if a != 0 {
			if _, ok := commonElements[a]; ok {
				p.oe.popUntilAtom(atom.Lookup(name))
				return true
			}
		}
		p.oe.popUntilName(name)
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
			p.path = append(p.path, PathComponent{
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

const hash0 = 0x84f70e16

func fnv(h uint32, s []byte) (uint32, bool, int) {
	for i := range s {
		if s[i] == '-' {
			return 0, true, i // hyphen encountered; short-circuit
		}
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h, false, 0
}

func attrIsNotDataOrAria(s []byte) bool {
	return !(len(s) > 5 &&
		s[4] == '-' && s[3] == 'a' &&
		((s[0] == 'd' && s[1] == 'a' && s[2] == 't') ||
			(s[0] == 'a' && s[1] == 'r' && s[2] == 'i')))
}

func camelize(b []byte, offset int) []byte {
	n := offset
	upperNext := false
	for i := offset; i < len(b); i++ {
		c := b[i]
		if c == '-' {
			upperNext = true
			continue
		}
		if upperNext && 'a' <= c && c <= 'z' {
			b[n] = c - 'a' + 'A'
			upperNext = false
		} else {
			b[n] = c
			upperNext = false
		}
		n++
	}
	return b[:n]
}
