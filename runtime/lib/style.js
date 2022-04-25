const upperCaseSecond = (_, char) => char.toUpperCase()

const styleToObject = require('style-to-object')

const keyPattern = /-(\w|$)/g

module.exports = val =>
  Object.entries(styleToObject(val)).reduce(
    (acc, [k, v]) => ({ ...acc, ...{ [k.replace(keyPattern, upperCaseSecond)]: v } }),
    {}
  )
