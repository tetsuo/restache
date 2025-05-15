package restache_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/tetsuo/restache"
	"golang.org/x/net/html/atom"
)

func TestRender(t *testing.T) {
	const file = "testdata/render_jsx.txt"
	for _, tc := range buildTestcases(t, file) {
		t.Run(fmt.Sprintf("%s L%d", file, tc.line), func(t *testing.T) {
			root, err := restache.Parse(strings.NewReader(tc.data))
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			var sb strings.Builder
			_, err = root.Render(&sb)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			got := sb.String()
			want := "export default function ($0) {return " + tc.expected + ";}"

			if got != want {
				t.Errorf("Render mismatch at line %d:\nwant:\n%s\ngot:\n%s\n", tc.line, want, got)
			}
		})
	}
}

type renderErrorCase struct {
	desc    string
	node    *restache.Node
	wantErr error
}

func TestRenderErrors(t *testing.T) {
	for _, tc := range []renderErrorCase{
		{
			desc:    "unknown node type",
			node:    &restache.Node{Type: restache.NodeType(999)},
			wantErr: restache.ErrUnknownNode,
		},
		{
			desc:    "error node",
			node:    &restache.Node{Type: restache.ErrorNode},
			wantErr: restache.ErrErrorNode,
		},
		{
			desc: "void element with children",
			node: func() *restache.Node {
				n := &restache.Node{
					Type:     restache.ElementNode,
					DataAtom: atom.Img,
				}
				n.AppendChild(&restache.Node{
					Type:     restache.ElementNode,
					DataAtom: atom.A,
				})
				return n
			}(),
			wantErr: restache.ErrVoidChildren,
		},
		{
			desc:    "when missing body",
			node:    &restache.Node{Type: restache.WhenNode, Data: "x"},
			wantErr: restache.ErrMissingBody,
		},
		{
			desc: "when too many children",
			node: func() *restache.Node {
				n := &restache.Node{Type: restache.WhenNode, Data: "x"}
				n.AppendChild(&restache.Node{Type: restache.TextNode, Data: "a"})
				n.AppendChild(&restache.Node{Type: restache.TextNode, Data: "b"})
				return n
			}(),
			wantErr: restache.ErrTooManyChildren,
		},
		{
			desc:    "range missing body",
			node:    &restache.Node{Type: restache.RangeNode, Data: "items"},
			wantErr: restache.ErrMissingBody,
		},
		{
			desc: "range too many children",
			node: func() *restache.Node {
				n := &restache.Node{Type: restache.RangeNode, Data: "items"}
				n.AppendChild(&restache.Node{Type: restache.TextNode, Data: "a"})
				n.AppendChild(&restache.Node{Type: restache.TextNode, Data: "b"})
				return n
			}(),
			wantErr: restache.ErrTooManyChildren,
		},
		{
			desc: "component with parent",
			node: func() *restache.Node {
				parent := &restache.Node{Type: restache.ElementNode}
				c := &restache.Node{Type: restache.ComponentNode, Data: "C"}
				parent.AppendChild(c)
				return c
			}(),
			wantErr: restache.ErrTopLevelOnly,
		},
		{
			desc: "text node at top level",
			node: func() *restache.Node {
				n := &restache.Node{Type: restache.TextNode, Data: "hello"}
				n2 := &restache.Node{Type: restache.TextNode, Data: "world"}
				n.AppendChild(n2)
				return n2
			}(),
			wantErr: restache.ErrChildOnly,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := restache.Render(io.Discard, tc.node)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}
