package parser_test

import (
	"log"
	"strings"
	"testing"

	. "github.com/onur1/stache/parser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		text     string
		expected *TagNode
	}{
		{
			text: `<x class={foo} for="{foo} {bar}">{bar}</x>`,
			expected: &TagNode{
				Name: "x",
				Children: []Node{
					&VariableNode{
						Name: "bar",
					},
				},
				Attrs: map[string][]Node{
					"class": {
						&VariableNode{
							Name: "foo",
						},
					},
					"for": {
						&VariableNode{
							Name: "foo",
						},
						&TextNode{
							Text: " ",
						},
						&VariableNode{
							Name: "bar",
						},
					},
				},
			},
		},
		{
			text: `<x><d>bla</d> {#foo} <y> <z>bom <c>{qux} {sox}</c><k></k></z>{quux} </y>bum {#bar}555<r></r><u></u>{/bar} 888{/foo} <k></k>333<d>{#qux}999{/qux}<h x="{nope}"></h></d> {jj}{dd}uu</x>`,
			expected: &TagNode{
				Name: "x",
				Children: []Node{
					&TagNode{
						Name: "d",
						Children: []Node{
							&TextNode{
								Text: "bla",
							},
						},
					},
					&TextNode{
						Text: " ",
					},
					&SectionNode{
						Name: "foo",
						Children: []Node{
							&TextNode{
								Text: " ",
							},
							&TagNode{
								Name: "y",
								Children: []Node{
									&TextNode{
										Text: " ",
									},
									&TagNode{
										Name: "z",
										Children: []Node{
											&TextNode{
												Text: "bom ",
											},
											&TagNode{
												Name: "c",
												Children: []Node{
													&VariableNode{
														Name: "qux",
													},
													&TextNode{
														Text: " ",
													},
													&VariableNode{
														Name: "sox",
													},
												},
											},
											&TagNode{
												Name: "k",
											},
										},
									},
									&VariableNode{
										Name: "quux",
									},
									&TextNode{
										Text: " ",
									},
								},
							},
							&TextNode{
								Text: "bum ",
							},
							&SectionNode{
								Name: "bar",
								Children: []Node{
									&TextNode{
										Text: "555",
									},
									&TagNode{
										Name: "r",
									},
									&TagNode{
										Name: "u",
									},
								},
							},
							&TextNode{
								Text: " 888",
							},
						},
					},
					&TextNode{
						Text: " ",
					},
					&TagNode{
						Name: "k",
					},
					&TextNode{
						Text: "333",
					},
					&TagNode{
						Name: "d",
						Children: []Node{
							&SectionNode{
								Name: "qux",
								Children: []Node{
									&TextNode{
										Text: "999",
									},
								},
							},
							&TagNode{
								Name: "h",
								Attrs: map[string][]Node{
									"x": {
										&VariableNode{
											Name: "nope",
										},
									},
								},
							},
						},
					},
					&TextNode{
						Text: " ",
					},
					&VariableNode{
						Name: "jj",
					},
					&VariableNode{
						Name: "dd",
					},
					&TextNode{
						Text: "uu",
					},
				},
			},
		},
		{
			text: `<tr bg="{#hasx}{x}{/hasx}"><td bg={name} ag="{name} is"></td></tr>`,
			expected: &TagNode{
				Name: "tr",
				Attrs: map[string][]Node{
					"bg": {
						&SectionNode{
							Name: "hasx",
							Children: []Node{
								&VariableNode{
									Name: "x",
								},
							},
						},
					},
				},
				Children: []Node{
					&TagNode{
						Name: "td",
						Attrs: map[string][]Node{
							"ag": {
								&VariableNode{
									Name: "name",
								},
								&TextNode{
									Text: " is",
								},
							},
							"bg": {
								&VariableNode{
									Name: "name",
								},
							},
						},
					},
				},
			},
		},
		{
			text: `<tr bg="{#hasx}{x}{/hasx}{^noty}{x}{/noty}"><td bg={name} ag="{name} is"></td></tr>`,
			expected: &TagNode{
				Name: "tr",
				Attrs: map[string][]Node{
					"bg": {
						&SectionNode{
							Name:     "hasx",
							Inverted: false,
							Children: []Node{
								&VariableNode{
									Name: "x",
								},
							},
						},
						&SectionNode{
							Name:     "noty",
							Inverted: true,
							Children: []Node{
								&VariableNode{
									Name: "x",
								},
							},
						},
					},
				},
				Children: []Node{
					&TagNode{
						Name: "td",
						Attrs: map[string][]Node{
							"ag": {
								&VariableNode{
									Name: "name",
								},
								&TextNode{
									Text: " is",
								},
							},
							"bg": {
								&VariableNode{
									Name: "name",
								},
							},
						},
					},
				},
			},
		},
		{
			text: `<tr>{! ehlo this is comment }</tr>`,
			expected: &TagNode{
				Name: "tr",
				Children: []Node{
					&CommentNode{
						Comment: " ehlo this is comment ",
					},
				},
			},
		},
		{
			text: `<tr>{!
    testing multiline comment
    ehlo this is comment
}</tr>`,
			expected: &TagNode{
				Name: "tr",
				Children: []Node{
					&CommentNode{
						Comment: "\n    testing multiline comment\n    ehlo this is comment\n",
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.text, func(t *testing.T) {
			if err := Parse(strings.NewReader(tt.text), func(tree *TagNode) bool {
				assert.EqualValues(t, tt.expected, tree)
				return true
			}); err != nil {
				log.Fatal(err)
			}
		})
	}
}
