const { htmlTags, selfClosingTags } = require('./defaults')
const Component = require('./component')

const createComponent = layout =>
  Component({
    ...layout,
    ...{
      opts: {
        ...layout.opts,
        ...{
          externs: { ...htmlTags, ...layout.opts.externs },
          selfClosingTags: { ...selfClosingTags, ...layout.opts.selfClosingTags },
        },
      },
    },
  })

module.exports = createComponent
