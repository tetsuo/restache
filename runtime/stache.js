const { isString, isNumber, isObject, has } = require('./function')

const Section = (name, children) => ({
  _tag: 'Section',
  name,
  children,
})

const InvertedSection = (name, children) => ({
  _tag: 'InvertedSection',
  name,
  children,
})

const Variable = name => ({ _tag: 'Variable', name })

const Comment = comment => ({ _tag: 'Comment', comment })

const Text = text => ({ _tag: 'Text', text })

const Element = (name, props, children, external) => ({
  _tag: 'Element',
  name,
  props,
  children,
  external,
})

const decode = (u, opts) => {
  if (Array.isArray(u)) {
    if (u.length < 1) {
      throw new TypeError(`expected a non-empty array, got ${JSON.stringify(u)}`)
    }
    const [_0, _1, _2] = u
    if (isString(_0)) {
      // element?
      if (_0.trim().length < 1) {
        throw new TypeError(`expected an element name, got empty string`)
      }
      const props = {}
      if (u.length > 1) {
        // has props?
        if (!isObject(_1)) {
          throw new TypeError(`expected an object, got ${JSON.stringify(_1)}`)
        }
        const keys = Object.keys(_1)
        let key
        let e
        let i = 0
        for (; i < keys.length; i++) {
          key = keys[i]
          if (!Array.isArray(_1[key])) {
            throw new TypeError(`expected an array, got ${JSON.stringify(_1[key])}`)
          }
          let j = 0
          const propChildren = []
          const arr = _1[key]
          for (; j < arr.length; j++) {
            e = decode(arr[j], opts)
            // validate children
            if (e._tag === 'Element') {
              throw new TypeError(`expected a valid property value for key "${key}", got ${JSON.stringify(e)}`)
            }
            propChildren.push(e)
          }
          props[key] = propChildren
        }
      }
      const children = []
      if (u.length > 2) {
        // has children?
        if (!Array.isArray(_2)) {
          throw `expected an array, got ${JSON.stringify(_2)}`
        }
        for (let i = 0; i < _2.length; i++) {
          children.push(decode(_2[i], opts))
        }
      }
      return Element(_0, props, children, has(opts.externs, _0) === false)
    } else if (isNumber(_0) && _0 >= 2 && _0 <= 5) {
      // stache element?
      if (u.length < 2) {
        throw new TypeError(`expected an array with 2 items at least, got ${JSON.stringify(u)}`)
      }
      if (!(isString(_1) && _1.trim().length > 0)) {
        throw new TypeError(`expected a non-empty string, got ${JSON.stringify(_1)}`)
      }
      if (_0 === 2 || _0 === 4) {
        const children = []
        let e
        if (u.length > 2) {
          for (let i = 0; i < _2.length; i++) {
            e = decode(_2[i], opts)
            children.push(e)
          }
        }
        if (_0 === 2) {
          return Section(_1, children)
        }
        return InvertedSection(_1, children)
      } else if (_0 === 3) {
        return Variable(_1)
      } else if (_0 === 5) {
        return Comment(_1)
      }
    }
    throw new TypeError(`expected a string or a number between 2-5, got ${JSON.stringify(_0)}`)
  }
  if (isString(u)) {
    return Text(u)
  }
  throw new Error(`expected an array or string, got ${JSON.stringify(u)}`)
}

module.exports = {
  Section,
  InvertedSection,
  Variable,
  Comment,
  Text,
  Element,
  decode,
}
