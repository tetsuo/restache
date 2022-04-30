const { createElement } = require('react')
const { connect } = require('react-redux')
const decode = require('./decode')
const sort = require('./sort')
const render = require('./render')
const { reactProps, reactSyntheticEvents } = require('./specs')

const hasOwnProperty = Object.prototype.hasOwnProperty

const mapPropName = (propName, tagName) =>
  hasOwnProperty.call(reactProps, propName) // react prop name?
    ? reactProps[propName]
    : hasOwnProperty.call(reactSyntheticEvents, propName) // react event name?
    ? reactSyntheticEvents[propName]
    : propName

const main = (roots, opts = {}) =>
  render(sort(roots.map(decode)), {
    ...opts,
    ...{
      getComponent: opts.getComponent || (() => null),
      mapPropName: opts.mapPropName || mapPropName,
      createElement: opts.createElement || createElement,
      connect: opts.connect || connect,
    },
  })

module.exports = main
