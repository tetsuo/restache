package lexer_test

import (
	"log"
	"strings"
	"testing"

	. "github.com/onur1/stache/lexer"
	"github.com/stretchr/testify/assert"
)

var na map[string][]Token

func TestTokenize(t *testing.T) {
	testCases := []struct {
		desc     string
		text     string
		expected []Token
	}{
		{
			desc: "tokenize html",
			text: `<table cols=3>
  <tbody>blah blah blah</tbody>
  <tr><td>there</td></tr>
  <tr><td>it</td></tr>
  <tr><td bgcolor="blue">is</td></tr>
</table>`,
			expected: []Token{
				{1, "table", map[string][]Token{"cols": {{3, "3", na}}}},
				{3, "\n  ", na},
				{1, "tbody", na},
				{3, "blah blah blah", na},
				{2, "tbody", na},
				{3, "\n  ", na},
				{1, "tr", na},
				{1, "td", na},
				{3, "there", na},
				{2, "td", na},
				{2, "tr", na},
				{3, "\n  ", na},
				{1, "tr", na},
				{1, "td", na},
				{3, "it", na},
				{2, "td", na},
				{2, "tr", na},
				{3, "\n  ", na},
				{1, "tr", na},
				{1, "td", map[string][]Token{"bgcolor": {{3, "blue", na}}}},
				{3, "is", na},
				{2, "td", na},
				{2, "tr", na},
				{3, "\n", na},
				{2, "table", na},
			},
		},
		{
			desc: "tokenize stache",
			text: `<table cols=3>
  {#fruits}
    <tr>
      <td bgcolor="blue">{name}</td>
      {#proteins}<td>{name}</td>{/proteins}
    </tr>
  {/fruits}
</table>`,
			expected: []Token{
				{1, "table", map[string][]Token{"cols": {{3, "3", na}}}},
				{3, "\n  ", na},
				{5, "fruits", na},
				{3, "\n    ", na},
				{1, "tr", na},
				{3, "\n      ", na},
				{1, "td", map[string][]Token{"bgcolor": {{3, "blue", na}}}},
				{4, "name", na},
				{2, "td", na},
				{3, "\n      ", na},
				{5, "proteins", na},
				{1, "td", na},
				{4, "name", na},
				{2, "td", na},
				{6, "proteins", na},
				{3, "\n    ", na},
				{2, "tr", na},
				{3, "\n  ", na},
				{6, "fruits", na},
				{3, "\n", na},
				{2, "table", na},
			},
		},
		{
			desc: "tokenize attrs",
			text: `<table cols={numcols}>
  {#fruits}
    <tr bgcolor="{#isx}{xthing}{/isx}">
      <td bgcolor="{bgcolor}">{name}</td>
      {#proteins}<td class='{foo} {bar} xxx'>{name}</td>{/proteins}
    </tr>
  {/fruits}
</table>`,
			expected: []Token{
				{1, "table", map[string][]Token{"cols": {{4, "numcols", na}}}},
				{3, "\n  ", na},
				{5, "fruits", na},
				{3, "\n    ", na},
				{1, "tr", map[string][]Token{"bgcolor": {{5, "isx", na}, {4, "xthing", na}, {6, "isx", na}}}},
				{3, "\n      ", na},
				{1, "td", map[string][]Token{"bgcolor": {{4, "bgcolor", na}}}},
				{4, "name", na},
				{2, "td", na},
				{3, "\n      ", na},
				{5, "proteins", na},
				{1, "td", map[string][]Token{"class": {{4, "foo", na}, {3, " ", na}, {4, "bar", na}, {3, " xxx", na}}}},
				{4, "name", na},
				{2, "td", na},
				{6, "proteins", na},
				{3, "\n    ", na},
				{2, "tr", na},
				{3, "\n  ", na},
				{6, "fruits", na},
				{3, "\n", na},
				{2, "table", na},
			},
		},
		{
			desc: "inverted sections",
			text: `<x>
  {^fruits}
    <k></k>
  {/fruits}
</x>`,
			expected: []Token{
				{1, "x", na},
				{3, "\n  ", na},
				{7, "fruits", na},
				{3, "\n    ", na},
				{1, "k", na},
				{2, "k", na},
				{3, "\n  ", na},
				{6, "fruits", na},
				{3, "\n", na},
				{2, "x", na},
			},
		},
		{
			desc: "comments",
			text: `<x>{   ! testasdsla
 sps os o s
}{  # ehlo }{  /ehlo }</x>`,
			expected: []Token{
				{1, "x", na},
				{8, " testasdsla\n sps os o s\n", na},
				{5, "ehlo", na},
				{6, "ehlo", na},
				{2, "x", na},
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			if err := Tokenize(strings.NewReader(tt.text), func(tk Token) bool {
				x := tt.expected[0]
				tt.expected = tt.expected[1:]
				assert.EqualValues(t, x, tk)
				return true
			}); err != nil {
				log.Fatal(err)
			}
		})
	}
}
