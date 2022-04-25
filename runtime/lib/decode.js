const t = require('@onur1/t')
const { Section, InvertedSection, Element, Variable, Text, Comment } = require('./tree')

const decode = u => {
  if (t.UnknownList(u)) {
    if (u.length < 1) {
      throw new TypeError('expected a non-empty array')
    }
    const [_0, _1, _2] = u
    if (t.String(_0)) {
      // is element?
      if (_0.trim().length < 1) {
        throw new TypeError('expected an element name, got empty string')
      }
      let props
      let forest = []
      const len = u.length
      if (len > 1) {
        props = {}
        if (!t.UnknownStruct(_1)) {
          throw new TypeError(`${_0}: expected a props object, got ${JSON.stringify(_1)}`)
        }
        const propNames = Object.keys(_1)
        let key
        let i = 0
        for (; i < propNames.length; i++) {
          key = propNames[i]
          if (key === 'key' || key === 'children' || key === 'ref') {
            throw new TypeError(`${_0} > ${key}: reserved property name`)
          }
          if (!t.UnknownList(_1[key])) {
            throw new TypeError(`${_0} > ${key}: expected an array of child nodes, got ${JSON.stringify(_1[key])}`)
          }
          let j = 0
          const propForest = []
          const arr = _1[key]
          for (; j < arr.length; j++) {
            propForest.push(decode(arr[j]))
          }
          props[key] = propForest
        }
        if (len > 2) {
          // has children?
          if (!t.UnknownList(_2)) {
            throw new TypeError(`${_0}: expected an array of child nodes, got ${JSON.stringify(_2)}`)
          }
          for (let i = 0; i < _2.length; i++) {
            forest.push(decode(_2[i]))
          }
        }
        // join consecutive text nodes
        i = 0
        let x, y
        for (; i < forest.length; i++) {
          if (i < 1) {
            continue
          }
          x = forest[i - 1].value
          y = forest[i].value
          if (x._type === 'Text' && y._type === 'Text') {
            forest[i - 1].value.text += y.text
            forest.splice(i, 1)
            i = i - 1
            continue
          }
        }
      }
      return Element(_0, props, forest)
    } else if (t.Integer(_0) && _0 >= 2 && _0 <= 5) {
      if (_0 === 2 || _0 === 4) {
        if (!t.NonEmptyString(_1)) {
          throw new TypeError(`expected a non-empty string as section name, got ${JSON.stringify(_1)}`)
        }
        const forest = []
        if (u.length > 2) {
          for (let i = 0; i < _2.length; i++) {
            forest.push(decode(_2[i]))
          }
        }
        if (_0 === 2) {
          return Section(_1, forest)
        }
        return InvertedSection(_1, forest)
      } else if (_0 === 3) {
        if (!t.NonEmptyString(_1)) {
          throw new TypeError(`expected a non-empty string as variable name, got ${JSON.stringify(_1)}`)
        }
        return Variable(_1)
      } else {
        if (!t.String(_1)) {
          throw new TypeError(`expected a string as comment, got ${JSON.stringify(_1)}`)
        }
        return Comment(_1)
      }
    }
    throw new TypeError(`expected a tag name or a valid id, got ${JSON.stringify(_0)}`)
  }
  if (t.String(u)) {
    return Text(u)
  }
  throw new TypeError(`expected an array or string, got ${JSON.stringify(u)}`)
}

module.exports = decode
