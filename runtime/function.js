const isString = u => typeof u === 'string'

const isNumber = u => typeof u === 'number'

const isObject = u => typeof u === 'object' && u !== null

const isFunction = u => typeof u === 'function'

const isBoolean = u => {
  const t = typeof u
  return t === 'undefined' || t === 'boolean' || (t === 'object' && u === null)
}

const hasOwnProperty = Object.prototype.hasOwnProperty

const has = (o, k) => hasOwnProperty.call(o, k)

const constant = a => () => a

const uniq = as => {
  if (as.length === 1) {
    return as.slice()
  }
  const out = [as[0]]
  const rest = [as[as.length - 1]]
  for (const a of rest) {
    if (out.every(o => o !== a)) {
      out.push(a)
    }
  }
  return out
}

const flatten = aas => {
  const r = []
  aas.forEach(as => {
    r.push(...as)
  })
  return r
}

module.exports = {
  isString,
  isNumber,
  isObject,
  isFunction,
  isBoolean,
  constant,
  uniq,
  flatten,
  has,
}
