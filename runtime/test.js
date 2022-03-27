const ReactDOMServer = require('react-dom/server')
const { Layout, Template } = require('./layout')
const createComponent = require('./react')
const t = require('tap')

const sectionTests = [
  {
    desc: 'object is accessible within the scope',
    html: '<span>{yar}{#bar}{nor}<a>{ssw}</a>{/bar}</span>',
    layout: Layout([
      Template('root', [
        [
          'span',
          {},
          [
            [3, 'yar'],
            [
              2,
              'bar',
              [
                [3, 'nor'],
                ['a', {}, [[3, 'ssw']]],
              ],
            ],
          ],
        ],
      ]),
    ]),
    props: {
      yar: '555',
      bar: [{ nor: 'xyz', ssw: 'rre' }],
    },
    expected: '<span>555xyz<a>rre</a></span>',
  },
  {
    desc: 'nested array sections',
    html: '<div>{#bar}<article>{#da}<span>{m}{#ow}8<b></b>{/ow}</span>{/da}</article>{#dor}{kw}{#mor}{sd}55{/mor}df{ke}1s{/dor}63{jh}{/bar}{kq}</div>',
    layout: Layout([
      Template('root', [
        [
          'div',
          {},
          [
            [
              2,
              'bar',
              [
                [
                  'article',
                  {},
                  [
                    [
                      2,
                      'da',
                      [
                        [
                          'span',
                          {},
                          [
                            [3, 'm'],
                            [2, 'ow', ['8', ['b', {}, []]]],
                          ],
                        ],
                      ],
                    ],
                  ],
                ],
                [2, 'dor', [[3, 'kw'], [2, 'mor', [[3, 'sd'], '55']], 'df', [3, 'ke'], '1s']],
                '63',
                [3, 'jh'],
              ],
            ],
            [3, 'kq'],
          ],
        ],
      ]),
    ]),
    props: {
      kq: 'qkw',
      bar: {
        da: {
          m: 'ram',
          ow: [1, 1, 1],
        },
        dor: {
          kw: 'ghj',
          mor: [{ sd: '77' }, { sd: '55' }],
          ke: 's1',
        },
        jh: '6734',
      },
    },
    expected: '<div><article><span>ram8<b></b>8<b></b>8<b></b></span></article>ghj77555555dfs11s636734qkw</div>',
  },
  {
    desc: 'sections 2',
    html: '<span>{#s}{s}{#z}{c}{/z}{#u}{k}{#r}{u}{#y}55{/y}{/r}{/u}{/s}</span>',
    layout: Layout([
      Template('root', [
        [
          'span',
          {},
          [
            [
              2,
              's',
              [
                [3, 's'],
                [2, 'z', [[3, 'c']]],
                [
                  2,
                  'u',
                  [
                    [3, 'k'],
                    [
                      2,
                      'r',
                      [
                        [3, 'u'],
                        [2, 'y', ['55']],
                      ],
                    ],
                  ],
                ],
              ],
            ],
          ],
        ],
      ]),
    ]),
    props: {
      s: {
        s: 'sd',
        z: false,
        u: [
          {
            k: '12',
            r: [{ u: '34' }, { u: '23', y: true }],
          },
          {
            k: '65',
            r: [{ u: '89' }],
          },
        ],
      },
    },
    expected: '<span>sd123423556589</span>',
  },
  {
    desc: 'interted section 1',
    html: '<span>{^s}foo{/s}</span>',
    layout: Layout([Template('root', [['span', {}, [[4, 's', ['foo']]]]])]),
    props: {
      s: true,
    },
    expected: '<span></span>',
  },
  {
    desc: 'interted section 2',
    html: '<span>{^s}foo{/s}</span>',
    layout: Layout([Template('root', [['span', {}, [[4, 's', ['foo']]]]])]),
    props: {
      s: false,
    },
    expected: '<span>foo</span>',
  },
  {
    desc: 'interted section 3',
    html: '<span>{^s}foo{/s}</span>',
    layout: Layout([Template('root', [['span', {}, [[4, 's', ['foo']]]]])]),
    props: {
      s: [],
    },
    expected: '<span>foo</span>',
  },
  {
    desc: 'interted section 4',
    html: '<span>{^s}{span}{/s}</span>',
    layout: Layout([Template('root', [['span', {}, [[4, 's', [[3, 'span']]]]]])]),
    props: {
      s: { span: 1 },
    },
    expected: '<span></span>',
  },
]

const errorTests = [
  {
    desc: 'invalid template',
    layout: Layout([Template('root', 42)]),
    message: 'expected a Layout object as its first parameter, got {"templates"',
  },
  {
    desc: 'empty element name',
    layout: Layout([Template('root', [['', {}, []]])]),
    message: 'expected an element name, got empty string',
  },
  {
    desc: 'invalid element props',
    layout: Layout([Template('root', [['div', 0, []]])]),
    message: 'expected an object, got 0',
  },
  {
    desc: 'invalid top-level prop children',
    layout: Layout([Template('root3', [['div', { k: 0 }, []]])]),
    message: 'expected an array, got 0',
  },
  {
    desc: 'invalid prop value',
    layout: Layout([Template('root3', [['div', { k: [['div', {}, []]] }, []]])]),
    message: 'expected a valid property value for key "k"',
  },
  // {
  //   desc: 'invalid children',
  //   layout: Layout([Template('root3', [['div', {}, 24]])]),
  //   message: 'expected an array with 2 items at least',
  // },
]

const componentDependencyTests = [
  {
    props: {
      itemsTop: [
        {
          label: 'Foo',
          href: '/foo',
          icon: 'ShowChartIcon',
        },
        {
          label: 'Bar',
          href: '/bar',
          icon: 'BarIcon',
        },
        {
          label: 'Baz',
          href: '/baz',
          icon: 'RouterIcon',
        },
      ],
    },
    expected:
      '<article class="permanent"><div>bla<div>here</div><div>hihello, Fooworld!<a href="/foo"><span><b>ShowChartIcon</b><i primary="Foo"></i></span></a>hihello, Barworld!<a href="/bar"><span><b>BarIcon</b><i primary="Bar"></i></span></a>hihello, Bazworld!<a href="/baz"><span><b>RouterIcon</b><i primary="Baz"></i></span></a></div></div></article>',
    layout: Layout([
      Template('sidebarlistitem', [
        'hi',
        'hello, ',
        [3, 'label'],
        'world!',
        [
          'a',
          { href: [[3, 'href']] },

          [
            [
              'span',
              {},
              [
                ['b', {}, [[3, 'icon']]],
                ['i', { primary: [[3, 'label']] }, []],
              ],
            ],
          ],
        ],
      ]),
      Template('sidebarlist', [
        [
          'div',
          {},
          [
            [
              2,
              'items',
              [
                ['sidebarlistitem', { label: [[3, 'label']], icon: [[3, 'icon']], href: [[3, 'href']] }, []],
                // 'hi',
                // 'hello, ',
                // [3, 'label'],
                // 'world!',
                // [
                //   'a',
                //   { href: [[3, 'href']] },

                //   [
                //     [
                //       'span',
                //       {},
                //       [
                //         ['b', {}, [[3, 'icon']]],
                //         ['i', { primary: [[3, 'label']] }, []],
                //       ],
                //     ],
                //   ],
                // ],
              ],
            ],
          ],
        ],
      ]),
      Template('sidebar', [
        [
          'article',
          { class: ['permanent'] },
          [
            ['div', {}, ['bla', ['div', {}, ['here']], ['sidebarlist', { items: [[3, 'itemsTop']] }, []]]],
          ],
        ],
      ]),
    ]),
  },
]

t.plan(sectionTests.length + componentDependencyTests.length + errorTests.length)

sectionTests.forEach(d => {
  const render = createComponent(d.layout)
  const fragment = render(d.props)
  const markup = ReactDOMServer.renderToStaticMarkup(fragment)
  if (markup === d.expected) {
    t.pass(d.html)
  } else {
    t.fail(`expected "${d.expected}", got "${markup}"`)
  }
})

componentDependencyTests.forEach(d => {
  const render = createComponent(d.layout)
  const fragment = render(d.props)
  const markup = ReactDOMServer.renderToStaticMarkup(fragment)
  if (markup === d.expected) {
    t.pass(d.html)
  } else {
    t.fail(`expected "${d.expected}", got "${markup}"`)
  }
})

errorTests.forEach(d => {
  try {
    createComponent(d.layout)
  } catch (e) {
    if (!e.message.startsWith(d.message)) {
      t.fail(`expected "${d.message}", got "${e.message}"`)
      return
    }
    t.pass(d.desc)
    return
  }
  t.fail(`expected "${d.expected}", got no error`)
})
