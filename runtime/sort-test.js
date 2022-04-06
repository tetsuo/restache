const { Element, Text, Section, InvertedSection, Variable, Comment } = require('./lib/domain')
const { runSortTest } = require('./test-util')

const sortTestCases = [
  {
    desc: 'dependency graph',
    input: [
      {
        name: 'bar',
        children: [Element('bar-1'), Element('qux', {}, [Element('foo')])],
      },
      {
        name: 'foo',
        children: [
          Element('foo-1', {}, [
            Element('foo-11'),
            Comment('hi'),
            Element('foo-12'),
            InvertedSection('foo-13', [
              Text('x'),
              Element('foo-131', {}, [Section('y', [Variable('test'), Element('foo-1311')])]),
            ]),
          ]),
        ],
      },
      {
        name: 'qux',
        children: [],
      },
    ],
    expected: [
      {
        name: 'qux',
        children: [],
      },
      {
        name: 'foo',
        children: [
          Element('foo-1', {}, [
            Element('foo-11'),
            Comment('hi'),
            Element('foo-12'),
            InvertedSection('foo-13', [
              Text('x'),
              Element('foo-131', {}, [Section('y', [Variable('test'), Element('foo-1311')])]),
            ]),
          ]),
        ],
      },
      {
        name: 'bar',
        children: [Element('bar-1'), Element('qux', {}, [Element('foo')])],
      },
    ],
  },
]

sortTestCases.map(runSortTest)
