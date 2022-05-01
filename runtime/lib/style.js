const upperCaseSecond = (_, char) => char.toUpperCase()

const styleToObject = require('style-to-object')

const keyPattern = /-(\w|$)/g

module.exports = val => {
  const o = styleToObject(val)
  if (o === null) {
    return {}
  }
  return Object.entries(o).reduce((acc, [k, v]) => ({ ...acc, ...{ [k.replace(keyPattern, upperCaseSecond)]: v } }), {})
}
