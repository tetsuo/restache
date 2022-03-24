const React = require('react')
const createComponent = require('.')

const createFragment = _2 => [_2]

const createReactComponent = layout =>
  createComponent({
    ...layout,
    ...{
      opts: {
        ...layout.opts,
        ...{
          createElement: React.createElement,
          createFragment: createFragment,
        },
      },
    },
  })

module.exports = createReactComponent
