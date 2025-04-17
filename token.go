package restache

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// TokenType represents the type of token.
type TokenType uint32

const (
	ErrorToken TokenType = iota
	StartTagToken
	EndTagToken
	SelfClosingTagToken
	TextToken
	CommentToken
	VariableToken
	WhenToken   // if
	UnlessToken // if not
	RangeToken  // for
	EndControlToken
)

// Tokenizer holds state for parsing.
type Tokenizer struct {
	z        *html.Tokenizer
	tt       TokenType
	err      error
	buf      []byte // the last chunk from z.Raw()
	pos      int    // current parse position in buf
	bufEnd   int    // end of buf
	tokBegin int    // start offset of current token in buf
	tokEnd   int    // end offset of current token in buf
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		z: html.NewTokenizer(r),
	}
}

// Err returns the last error encountered by the tokenizer.
func (t *Tokenizer) Err() error {
	return t.err
}

// Raw returns the raw byte slice of the current token.
func (t *Tokenizer) Raw() []byte {
	return t.buf[t.tokBegin:t.tokEnd]
}

// TagName returns the name of the current HTML tag, if applicable.
func (t *Tokenizer) TagName() ([]byte, bool) {
	return t.z.TagName()
}

// Comment extracts the content of a comment.
func (t *Tokenizer) Comment() []byte {
	// Find the first comment symbol, and return the rest
	b := t.Raw()
	i := bytes.IndexByte(b, '!')
	return b[i+1:]
}

// ControlName extracts the name of a section, inverted section, or end section.
func (t *Tokenizer) ControlName() []byte {
	// Find the first section symbol, and return the rest
	b := t.Raw()
	var i int
	switch t.tt {
	case WhenToken:
		i = bytes.IndexByte(b, '?')
	case RangeToken:
		i = bytes.IndexByte(b, '#')
	case EndControlToken:
		i = bytes.IndexByte(b, '/')
	case UnlessToken:
		i = bytes.IndexByte(b, '^')
	}
	return b[i+1:]
}

// TagAttr retrieves the next attribute key and value from an HTML start tag.
func (t *Tokenizer) TagAttr() (key []byte, val []byte, isExpr bool, moreAttr bool) {
	key, val, moreAttr = t.z.TagAttr()
	if !(len(key) > 5 &&
		key[4] == '-' && key[3] == 'a' &&
		((key[0] == 'd' && key[1] == 'a' && key[2] == 't') ||
			(key[0] == 'a' && key[1] == 'r' && key[2] == 'i'))) {
		key = kebabToCamel(key)
	}
	i := 0
	n := len(val)
	for i < n && spaceTable[val[i]] {
		i++
	}
	if i == n || val[i] != '{' { // No starting '{' found; text node
		return
	}
	rpos := bytes.IndexByte(val[i:], '}')
	if rpos < 0 { // It is text node if no closing '}' found
		return
	}
	if i+rpos+1 == n {
		expr := val[i:][1:rpos]
		val = expr
		isExpr = true
		return
	}
	end := i + rpos + 1
	for end < n {
		if !spaceTable[val[end]] { // If anything other than space; text node
			return
		}
		end++
	}
	val = val[i:][1:rpos]
	isExpr = true
	return
}

// Next advances the tokenizer to the next token and returns its type.
func (t *Tokenizer) Next() TokenType {
	// Seen any errors?
	if t.err != nil {
		t.tt = ErrorToken
		return t.tt
	}

	// Still have leftover text in buf?
	if t.pos < t.bufEnd {
		t.parseTextSegment()
		return t.tt
	}

consume:
	for {
		tt := t.z.Next()
		switch tt {
		case html.ErrorToken:
			t.err = t.z.Err()
			t.tt = ErrorToken
			return t.tt
		case html.TextToken:
			t.buf = t.z.Raw()
			t.pos = 0
			t.bufEnd = len(t.buf)
			t.parseTextSegment()
			return t.tt
		case html.StartTagToken:
			t.tt = StartTagToken
			break consume
		case html.SelfClosingTagToken:
			t.tt = SelfClosingTagToken
			break consume
		case html.EndTagToken:
			t.tt = EndTagToken
			break consume
		}
	}
	t.buf = t.z.Raw()
	length := len(t.buf)
	t.pos = length // This entire chunk consumed
	t.bufEnd = length
	t.tokBegin = 0
	t.tokEnd = length
	return t.tt
}

func (t *Tokenizer) parseTextSegment() {
	b := t.buf
	start := t.pos

	// Find the next '{'
	lpos := bytes.IndexByte(b[start:], '{')
	if lpos < 0 {
		// No '{' found => everything left is normal text
		t.tt = TextToken
		t.tokBegin = start
		t.tokEnd = t.bufEnd
		t.pos = t.bufEnd
		return
	}
	lpos += start // adjust lpos to absolute index in b

	// If there's text before the '{', return that text first
	if lpos > start {
		t.tt = TextToken
		t.tokBegin = start
		t.tokEnd = lpos
		t.pos = lpos // Next time we call parseTextSegment, we handle the '{'
		return
	}

	// If we get here, it means b[start] == '{', find the matching '}'
	rpos := bytes.IndexByte(b[lpos+1:], '}')
	if rpos < 0 {
		// No closing '}' => treat the rest as text
		t.tt = TextToken
		t.tokBegin = lpos
		t.tokEnd = t.bufEnd
		t.pos = t.bufEnd
		return
	}
	rpos += (lpos + 1)

	t.tt = identifyKeyword(b[lpos+1 : rpos])

	t.tokBegin = lpos + 1
	t.tokEnd = rpos

	t.pos = rpos + 1
}

// identifyKeyword looks at the content inside {...} and decides the token type.
func identifyKeyword(chunk []byte) TokenType {
	// Skip leading spaces
	i := 0
	length := len(chunk)
	for i < length && spaceTable[chunk[i]] {
		i++
	}
	if i >= length {
		return VariableToken
	}
	switch chunk[i] {
	case '?':
		return WhenToken
	case '^':
		return UnlessToken
	case '#':
		return RangeToken
	case '/':
		return EndControlToken
	case '!':
		return CommentToken
	default:
		return VariableToken
	}
}

var spaceTable = [256]bool{
	' ': true, '\t': true, '\r': true, '\n': true,
}

func kebabToCamel(b []byte) []byte {
	n := 0
	upperNext := false
	for i := range b {
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
