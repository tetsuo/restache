const { Element, Text, Section, InvertedSection, Variable, Comment } = require('./lib/domain')
const { externs, selfClosingTags, syntheticEvents } = require('./lib/spec')
const { runRenderTest } = require('./test-util')
const { createElement } = require('react')

const getopts = (externs = {}, registry = {}, selfClosingTags = {}, syntheticEvents = {}, externProps = {}) => ({
  externs,
  registry,
  selfClosingTags,
  syntheticEvents,
  externProps,
})

const customComponent = props => createElement('a', { href: '#' }, props.children)

const renderTestCases = [
  {
    desc: 'text',
    input: {
      trees: [Element('foo', {}, [Text('qux')])],
      opts: getopts(),
    },
    expected: 'qux',
  },
  {
    desc: 'text join',
    input: {
      trees: [Element('foo', {}, [Text('qux'), Text(', quux!')])],
      opts: getopts(),
    },
    expected: 'qux, quux!',
  },
  {
    desc: 'comment',
    input: {
      trees: [Element('foo', {}, [Text('qux'), Comment('test'), Text(', quux!')])],
      opts: getopts(),
    },
    expected: 'qux, quux!',
  },
  {
    desc: 'variable',
    input: {
      trees: [Element('foo', {}, [Variable('a'), Variable('b')])],
      opts: getopts(),
      props: {
        a: 'qu',
        b: ', qux!',
      },
    },
    expected: 'qu, qux!',
  },
  {
    desc: 'self closing tag skip children',
    input: {
      trees: [Element('foo', {}, [Element('div', {}, [Text('hello')])])],
      opts: getopts({ div: true }, {}, { div: true }),
    },
    expected: '<div></div>',
  },
  {
    desc: 'render children',
    input: {
      trees: [Element('foo', {}, [Element('div', {}, [Text('hello'), Element('span', {}, [Text('world')])])])],
      opts: getopts({ div: true, span: true }),
    },
    expected: '<div>hello<span>world</span></div>',
  },
  {
    desc: 'inverted section',
    input: {
      trees: [
        Element('foo', {}, [
          InvertedSection('c', [Text('zzz')]),
          InvertedSection('d', [Text('xxx')]),
          InvertedSection('e', [Text('abc')]),
          InvertedSection('f', [Text('def')]),
          InvertedSection('g', [Text('yyy')]),
          InvertedSection('h', [Text('ghi')]),
          InvertedSection('a', [Text('beep')]),
          InvertedSection('b', [Text('bip')]),
          InvertedSection('boop', [Text('boop')]),
          Text('bar'),
        ]),
      ],
      opts: getopts(),
      props: {
        c: [1],
        f: false,
        d: { a: 4 },
        g: {},
        h: [],
        a: undefined,
        b: null,
      },
    },
    expected: 'abcdefghibeepbipboopbar',
  },
  {
    desc: 'section',
    input: {
      trees: [
        Element('foo', {}, [
          Section('c', [Text('zzz')]),
          Section('d', [Text('xxx')]),
          Section('e', [Text('abc')]),
          Section('f', [Text('def')]),
          Section('g', [Text('yyy')]),
          Section('h', [Text('ghi')]),
          Section('i', [Text('ba')]),
          Section('j', [Text('sf')]),
          Section('k', [Text('rf')]),
          Section('a', [Text('beep')]),
          Section('b', [Text('bip')]),
          Section('boop', [Text('boop')]),
          Text('bar'),
        ]),
      ],
      opts: getopts(),
      props: {
        c: [1, 2],
        f: true,
        d: { a: 4, b: 2 },
        g: {},
        h: [],
        a: undefined,
        b: null,
        i: false,
        j: Symbol(''),
      },
    },
    expected: 'zzzzzzxxxdefyyyghisfbar',
  },
  {
    desc: 'section variable scope',
    input: {
      trees: [
        Element('foo', {}, [
          Section('a', [Variable('x'), Text('hi'), Variable('y')]),
          Section('b', [Variable('x'), Text('hi')]),
          Section('c', [Variable('x'), Text('hej')]),
          Section('d', [Variable('x')]),
          Section('e', [Variable('x')]),
          Section('f', [Variable('x')]),
          Section('g', [Variable('x')]),
          Section('h', [Variable('x')]),
        ]),
      ],
      opts: getopts(),
      props: {
        x: 'ab',
        a: {
          x: 'a1',
          y: 'a2',
        },
        b: ['b', 3, {}, { x: 'b4' }, { x: 'b5' }],
        c: {},
        d: { x: true }, // TODO:
        e: { x: false }, // TODO:
        f: { x: undefined }, // TODO:
        g: { x: null }, // TODO:
        h: { x: 0 }, // TODO:
      },
    },
    expected: 'a1hia2undefinedhiundefinedhiundefinedhib4hib5hiundefinedhejtruefalseundefinednull0',
  },
  {
    desc: 'section variable scope',
    input: {
      trees: [
        Element('foo', {}, [
          Section('a', [Variable('x'), Text('hi'), Variable('y')]),
          Section('b', [Variable('x'), Text('hi')]),
          Section('c', [Variable('x'), Text('hej')]),
          Section('d', [Variable('x')]),
          Section('e', [Variable('x')]),
          Section('f', [Variable('x')]),
          Section('g', [Variable('x')]),
          Section('h', [Variable('x')]),
          Section('i', [Section('j', [Variable('x')])]),
        ]),
      ],
      opts: getopts(),
      props: {
        x: 'ab',
        a: {
          x: 'a1',
          y: 'a2',
        },
        b: ['b', 3, {}, { x: 'b4' }, { x: 'b5' }],
        c: {},
        d: { x: true },
        e: { x: false },
        f: { x: undefined },
        g: { x: null },
        h: { x: 0 },
        i: {
          j: {
            x: 'hoi',
          },
        },
      },
    },
    expected: 'a1hia2undefinedhiundefinedhiundefinedhib4hib5hiundefinedhejtruefalseundefinednull0hoi',
  },
  {
    desc: 'dependent components',
    input: {
      trees: [
        Element('qux', {}, [Text('!')]),
        Element('baz', {}, [Text('world')]),
        Element('bar', {}, [Text(', '), Element('baz'), Element('qux')]),
        Element('foo', {}, [Text('Hello'), Element('bar', {}, [])]),
      ],
      opts: getopts(),
      props: {},
    },
    expected: 'Hello, world!',
  },
  {
    desc: 'custom component',
    input: {
      trees: [
        Element('baz', {}, [Element('qux')]),
        Element('bar', {}, [Element('baz')]),
        Element('foo', {}, [Text('anchor'), Element('bar')]),
      ],
      opts: getopts(
        {},
        {
          qux: customComponent,
        }
      ),
      props: {},
    },
    expected: 'anchor<a href="#"></a>',
  },
  {
    desc: 'custom component with children',
    input: {
      trees: [
        Element('baz', {}, [Element('qux', {}, [Text('hi')])]),
        Element('bar', {}, [Element('baz')]),
        Element('foo', {}, [Text('anchor'), Element('bar')]),
      ],
      opts: getopts(
        {},
        {
          qux: customComponent,
        }
      ),
      props: {},
    },
    expected: 'anchor<a href="#">hi</a>',
  },
  {
    desc: 'own component as custom component children',
    input: {
      trees: [
        Element('qux', {}, [Text('hej'), Element('baz', {}, [Text('hi')])]),
        Element('bar', {}, [Element('baz', {}, [Element('qux')])]),
        Element('foo', {}, [Element('bar')]),
      ],
      opts: getopts(
        {},
        {
          baz: customComponent,
        }
      ),
      props: {},
    },
    expected: '<a href="#">hej<a href="#">hi</a></a>',
  },
  {
    desc: 'pass text property to own component',
    input: {
      trees: [
        Element('bar', {}, [Variable('x')]),
        Element('foo', {}, [
          Element('bar', {
            x: [Text('hi')],
          }),
        ]),
      ],
      opts: getopts(),
      props: {},
    },
    expected: 'hi',
  },
  {
    desc: 'pass multiple text property to own component',
    input: {
      trees: [
        Element('bar', {}, [Variable('x')]),
        Element('foo', {}, [
          Element('bar', {
            x: [Text('hi'), Text('hello')],
          }),
        ]),
      ],
      opts: getopts(),
      props: {},
    },
    expected: 'hihello',
  },
  {
    desc: 'pass props to own component 1',
    input: {
      trees: [
        Element('bar', {}, [Variable('a'), Variable('b'), Variable('c'), Variable('d'), Variable('e')]),
        Element('foo', {}, [
          Element('bar', {
            a: [Variable('y')],
            b: [Variable('z')],
            c: [Variable('bar'), Variable('baz'), Text('nnn')],
            d: [Variable('qux'), Variable('quux')],
            e: [Variable('faz')],
          }),
        ]),
      ],
      opts: getopts(),
      props: {
        y: 'hej',
        z: 42,
        bar: 555,
        baz: 100,
        qux: {},
        quux: [1, 2, 3],
        faz: {},
      },
    },
    expected: 'hej42555100nnn[object Object]1,2,3[object Object]',
  },
  {
    desc: 'pass props to own component 2',
    input: {
      trees: [
        Element('bar', {}, [
          Variable('e'),
          Element(
            'div',
            {
              value: [Variable('h'), Text('hi')],
              checked: [Variable('e'), Text('hi')],
              id: [Variable('e'), Text('hi')],
            },
            [Text('letsgo')]
          ),
          Element('input', {
            type: [Text('check'), Text('box')],
            id: [Variable('f')],
            for: [Text('label'), Variable('f')],
            checked: [Variable('e')],
            style: [Variable('g')],
            class: [],
            email: [Variable('h')],
            value: [Variable('e'), Variable('f')],
            onclick: [Variable('fn')],
          }),
          Element('div', {
            class: [Section('s1', [Section('s2', [Text('xu')])]), InvertedSection('s3', [Text('bu')])],
          }),
        ]),
        Element('foo', {}, [
          Element('bar', {
            e: [Variable('faz')],
            f: [Variable('qux')],
            g: [Variable('style')],
            h: [Variable('quux')],
            fn: [Variable('xx')],
            s1: [Variable('section1')],
            s3: [Variable('section2')],
          }),
        ]),
      ],
      opts: getopts(
        {
          input: true,
          div: true,
        },
        {},
        selfClosingTags,
        syntheticEvents
      ),
      props: {
        faz: true,
        qux: {},
        quux: undefined,
        style: {
          marginRight: '3em',
        },
        xx: () => null,
        section1: {
          s2: [1, 2, 3],
        },
        section2: false,
      },
    },
    expected:
      'true<div value="undefinedhi" checked="" id="truehi">letsgo</div><input type="checkbox" id="[object Object]" for="label[object Object]" style="margin-right:3em" value="true[object Object]" checked=""/><div class="xu,xu,xubu"></div>',
  },
  {
    desc: 'pass props to own component 2',
    input: {
      trees: [
        Element('bar', {}, [
          Element(
            'div',
            {
              id: [Element('span', {}, [Text('hi')])],
            },
            [Text('letsgo')]
          ),
        ]),
        Element('foo', {}, [
          Element('bar', {
            foo: [Variable('foo')],
          }),
        ]),
      ],
      opts: getopts(externs, {}, selfClosingTags, syntheticEvents),
      props: {
        foo: 'test',
      },
    },
    expectedErr: 'span: elements are not valid as prop children',
  },
]

renderTestCases.map(runRenderTest)
