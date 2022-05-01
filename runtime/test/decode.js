const { Element, Text, Section, InvertedSection, Variable, Comment } = require('../lib/tree')
const { runDecodeTest } = require('../test-util')

const decodeTestCases = [
  {
    desc: 'invalid type',
    expectedErr: `expected an array or string, got 42`,
    input: 42,
  },
  {
    desc: 'invalid type',
    expectedErr: `expected an array or string, got null`,
    input: null,
  },
  {
    desc: 'invalid type',
    expectedErr: `expected an array or string, got undefined`,
  },
  {
    desc: 'text',
    input: 'hello',
    expected: Text('hello'),
  },
  {
    desc: 'text preserve whitespace',
    input: '  hello   \n\t   ',
    expected: Text('  hello   \n\t   '),
  },
  {
    desc: 'element name empty',
    input: [''],
    expectedErr: 'expected an element name, got empty string',
  },
  {
    desc: 'empty array',
    input: [],
    expectedErr: 'expected a non-empty array',
  },
  {
    desc: 'element invalid props',
    input: ['a', null],
    expectedErr: 'a: expected a props object, got null',
  },
  {
    desc: 'element invalid prop children',
    input: ['a', { b: null }],
    expectedErr: 'a > b: expected an array of child nodes, got null',
  },
  {
    desc: 'element basic',
    input: ['a'],
    expected: Element('a'),
  },
  {
    desc: 'element with empty props',
    input: ['a', { b: [], c: [] }],
    expected: Element('a', { b: [], c: [] }),
  },
  {
    desc: 'element with a non-empty prop',
    input: ['a', { b: ['hello', 'world'], c: [] }],
    expected: Element('a', { b: [Text('hello'), Text('world')], c: [] }),
  },
  {
    desc: 'element with invalid children',
    input: ['a', {}, 42],
    expectedErr: 'a: expected an array of child nodes, got 42',
  },
  {
    desc: 'element with children',
    input: ['a', {}, ['hello', 'world', 'abc', 'def', [3, 'x'], 'ghi', 'jkl']],
    expected: Element('a', {}, [Text('helloworldabcdef'), Variable('x'), Text('ghijkl')]),
  },
  {
    desc: 'invalid tag id',
    input: [42],
    expectedErr: 'expected a tag name or a valid id, got 42',
  },
  {
    desc: 'invalid section name',
    input: [2],
    expectedErr: 'expected a non-empty string as section name, got undefined',
  },
  {
    desc: 'invalid inverted section name',
    input: [4, 42],
    expectedErr: 'expected a non-empty string as section name, got 42',
  },
  {
    desc: 'invalid variable name',
    input: [3, []],
    expectedErr: 'expected a non-empty string as variable name, got []',
  },
  {
    desc: 'invalid variable name',
    input: [3],
    expectedErr: 'expected a non-empty string as variable name, got undefined',
  },
  {
    desc: 'section with no children',
    input: [2, 'a'],
    expected: Section('a'),
  },
  {
    desc: 'inverted section with no children',
    input: [4, 'a'],
    expected: InvertedSection('a'),
  },
  {
    desc: 'section with children',
    input: [2, 'a', ['hello']],
    expected: Section('a', [Text('hello')]),
  },
  {
    desc: 'inverted section with children',
    input: [4, 'a', ['hello']],
    expected: InvertedSection('a', [Text('hello')]),
  },
  {
    desc: 'variable',
    input: [3, 'a'],
    expected: Variable('a'),
  },
  {
    desc: 'invalid comment',
    input: [5, null],
    expectedErr: 'expected a string as comment, got null',
  },
  {
    desc: 'comment',
    input: [5, 'hi'],
    expected: Comment('hi'),
  },
  {
    desc: 'reserved prop name',
    input: ['a', { children: [] }, 42],
    expectedErr: 'a > children: reserved property name',
  },
]

decodeTestCases.map(runDecodeTest)
