package restache_test

import (
	"strings"
	"testing"

	"github.com/tetsuo/restache"
)

type renderTestCase struct {
	desc string
	tmpl string
	want string
	doc  *restache.Node
	err  error
}

func TestRender(t *testing.T) {
	for _, tc := range []renderTestCase{
		{
			"empty",
			"",
			"<></>",
			nil, nil,
		},
		{
			"text no space",
			"hi",
			"<>hi</>",
			nil, nil,
		},
		{
			"text with leading and trailing space",
			"    hi    ",
			"<>hi</>", // TODO: should retain one space
			nil, nil,
		},
		{
			"text space between words",
			"hello     world",
			"<>hello world</>",
			nil, nil,
		},
		{
			"variable",
			"{somevar}",
			"$0.somevar",
			nil, nil,
		},
		{
			"two variables",
			"{somevar}{anothervar}",
			"<>{$0.somevar}{$0.anothervar}</>",
			nil, nil,
		},
		{
			"element",
			"<span></span>",
			"<span></span>",
			nil, nil,
		},
		{
			"multiple elements",
			"<span>hi</span><b>hello</b>",
			"<><span>hi</span><b>hello</b></>",
			nil, nil,
		},
		{
			"void element",
			"<img>",
			"<img />",
			nil, nil,
		},
		{
			"void element 2", // ensure non-void behaves differently
			"<span>",
			"<span></span>",
			nil, nil,
		},
		{
			"when",
			"{?x}{/x}",
			"($0.x && <></>)",
			nil, nil,
		},
		{
			"when with single text",
			"{?hi}sup{/hi}",
			"($0.hi && <>sup</>)",
			nil, nil,
		},
		{
			"when with single element",
			"{?hi}<span></span>{/hi}",
			"($0.hi && <span></span>)",
			nil, nil,
		},
		{
			"when with multi elements",
			"{?hi}<img><hr>{/hi}",
			"($0.hi && <><img /><hr /></>)",
			nil, nil,
		},
		{
			"unless",
			"{^x}{/x}",
			"(!$0.x && <></>)",
			nil, nil,
		},
		{
			"range",
			"{#x}{/x}",
			"$0.x.map($1 => <></>)",
			nil, nil,
		},
		{
			"range with single text",
			"{#hi}hi{/hi}",
			"$0.hi.map($1 => <React.Fragment key={ $1.key }>hi</React.Fragment>)",
			nil, nil,
		},
		{
			"comment",
			"{!Hi hello}",
			"<>{ /* Hi hello */ }</>",
			nil, nil,
		},
		// TODO: ComponentNode, ErrorNode
	} {
		t.Run(tc.desc, func(t *testing.T) {
			assertRender(t, tc)
		})
	}
}

func assertRender(t *testing.T, tc renderTestCase) {
	root := tc.doc
	var err error
	if root == nil {
		r := strings.NewReader(tc.tmpl)
		root, err = restache.Parse(r)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
	}

	var sb strings.Builder
	_, err = root.Render(&sb)

	if tc.err != nil {
		if err == nil || err != tc.err {
			t.Fatalf("expected render error %v, got %v", tc.err, err)
		}
		return
	}

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	tc.want = "export default function ($0) {return " + tc.want + ";}"
	got := sb.String()
	if got != tc.want {
		t.Errorf("Render() = %q, want %q", got, tc.want)
	}
}
