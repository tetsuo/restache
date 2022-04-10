package lexer

import (
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type TokenKind int

const (
	Open TokenKind = iota + 1
	Close
	Text
	Variable
	SectionOpen
	SectionClose
	InvertedSectionOpen
	Comment
)

type Token struct {
	Kind  TokenKind
	Body  string
	Attrs map[string][]Token
}

func Tokenize(r io.Reader, cb func(Token) bool) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if err := z.Err(); err != io.EOF {
				return err
			}
			return nil
		case html.TextToken:
			if err := scan(string(z.Text()), func(tk TokenKind, s string) bool {
				ok := cb(Token{Kind: tk, Body: s})
				return ok
			}); err != nil {
				return err
			}
		case html.SelfClosingTagToken, html.StartTagToken:
			tk := z.Token()
			t := Token{
				Kind: Open,
				Body: tk.Data,
			}
			if length := len(tk.Attr); length > 0 {
				attrs := make(map[string][]Token)
				for _, v := range tk.Attr {
					var sub []Token
					if err := scan(v.Val, func(tk TokenKind, s string) bool {
						sub = append(sub, Token{Kind: tk, Body: s})
						return true
					}); err != nil {
						return err
					}
					attrs[v.Key] = sub
				}
				t.Attrs = attrs
			}
			cb(t)
			if tt == html.SelfClosingTagToken {
				cb(Token{Kind: Close, Body: tk.Data})
			}
		case html.EndTagToken:
			data := z.Token().Data
			cb(Token{Kind: Close, Body: data})
		}
	}
}

const (
	tagSymbolVar                 = ""
	tagSymbolSectionOpen         = "#"
	tagSymbolSectionClose        = "/"
	tagSymbolInvertedSectionOpen = "^"
	tagSymbolComment             = "!"
)

func kindFromSymbol(s string) TokenKind {
	var kind TokenKind
	switch s {
	case tagSymbolVar:
		kind = Variable
	case tagSymbolSectionOpen:
		kind = SectionOpen
	case tagSymbolSectionClose:
		kind = SectionClose
	case tagSymbolInvertedSectionOpen:
		kind = InvertedSectionOpen
	case tagSymbolComment:
		kind = Comment
	}
	return kind
}

var (
	tokenMatcher = regexp.MustCompile("{[^}]+}")
	tagMatcher   = regexp.MustCompile(`^{\s*([#\/\^\!]?)([\s\S]+)s*}$`)
)

func findAllIndex(s string, re *regexp.Regexp) [][]int {
	xs := re.FindAllStringSubmatchIndex(s, -1)
	i, w := 0, 0
	var cur []int
	for ; i < len(xs); i++ {
		cur = xs[i]
		if cur[0] > w {
			xs = append(xs[:i], append([][]int{{w, cur[0]}}, xs[i:]...)...)
			i++
		}
		w = cur[1]
	}
	if w < len(s) {
		xs = append(xs, []int{w, len(s)})
	}

	return xs
}

func scan(s string, cb func(TokenKind, string) bool) error {
	tokens := findAllIndex(s, tokenMatcher)
	length := len(tokens)

	var (
		cur   []int
		t     string
		match [][]string
		r     int
		w     int
		kind  TokenKind
		name  string
	)

	i := 0
	for i < length {
		cur = tokens[i]
		r, w = cur[0], cur[1]
		t = s[r:w]

		match = tagMatcher.FindAllStringSubmatch(t, -1)
		if len(match) == 0 {
			kind, name = Text, t
		} else {
			kind, name = kindFromSymbol(match[0][1]), match[0][2]
			if kind != Comment {
				name = strings.TrimSpace(name)
			}
		}
		if ok := cb(kind, name); !ok {
			return nil
		}

		i++
	}

	return nil
}
