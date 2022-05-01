const { Element, Text, Section, InvertedSection, Variable, Comment } = require('../lib/tree')
const { runSortTest } = require('../test-util')

const sortTestCases = [
  {
    desc: 'dependency graph',
    input: [
      Element('bar', {}, [Element('bar-1'), Element('qux', {}, [Element('foo')])]),
      Element('foo', {}, [
        Element('foo-1', {}, [
          Element('foo-11'),
          Comment('hi'),
          Element('foo-12'),
          InvertedSection('foo-13', [
            Text('x'),
            Element('foo-131', {}, [Section('y', [Variable('test'), Element('foo-1311')])]),
          ]),
        ]),
      ]),
      Element('qux'),
    ],
    expected: [
      Element('qux'),
      Element('foo', {}, [
        Element('foo-1', {}, [
          Element('foo-11'),
          Comment('hi'),
          Element('foo-12'),
          InvertedSection('foo-13', [
            Text('x'),
            Element('foo-131', {}, [Section('y', [Variable('test'), Element('foo-1311')])]),
          ]),
        ]),
      ]),
      Element('bar', {}, [Element('bar-1'), Element('qux', {}, [Element('foo')])]),
    ],
  },
]

sortTestCases.map(runSortTest)
