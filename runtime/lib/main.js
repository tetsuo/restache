const decode = require('./decode')
const sort = require('./sort')
const render = require('./render')
const { externs, selfClosingTags, syntheticEvents, externProps } = require('./spec')

const main = (roots, opts = {}) =>
  render(sort(roots.map(decode)), {
    ...opts,
    ...{
      registry: opts.registry || {},
      externs: { ...externs, ...opts.externs },
      selfClosingTags: { ...selfClosingTags, ...opts.selfClosingTags },
      syntheticEvents: { ...syntheticEvents, ...opts.syntheticEvents },
      externProps: { ...externProps, ...opts.externProps },
    },
  })

module.exports = main
