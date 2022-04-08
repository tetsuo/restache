const { runMainTest } = require('./test-util')

const mainTestCases = [
  {
    desc: 'text',
    input: {
      roots: [
        [
          'list',
          {},
          [
            ['h2', {}, ['Fruits']],
            ['ul', {}, [[2, 'items', [['listitem', { name: [[3, 'name']] }, []]]]]],
          ],
        ],
        ['listitem', {}, [['li', {}, [[3, 'name']]]]],
      ],
      props: {
        items: [{ name: 'Apple' }, { name: 'Orange' }],
      },
    },
    expected: '<h2>Fruits</h2><ul><li>Apple</li><li>Orange</li></ul>',
  },
]

mainTestCases.map(runMainTest)
