const hasOwnProperty = Object.prototype.hasOwnProperty

const MIN_SAFE_INTEGER = -9007199254740991

const MAX_SAFE_INTEGER = 9007199254740991

const Number = x => typeof x === 'number' && x <= MAX_SAFE_INTEGER && x >= MIN_SAFE_INTEGER

const Integer = x => Number(x) && x === Math.floor(x)

const UnknownStruct = u => typeof u === 'object' && !Nil(u) && !UnknownList(u)

const string = x => typeof x === 'string'

const NonEmptyString = x => string(x) && x.length > 0

const Nil = x => x === null

const Undefined = u => u === void 0

const Boolean = x => x === true || x === false

const Struct = props => {
  const keys = Object.keys(props)
  const types = keys.map(key => props[key])
  const len = keys.length
  return u => {
    if (UnknownStruct(u)) {
      for (let i = 0; i < len; i++) {
        const k = keys[i]
        const uk = u[k]
        if ((uk === undefined && !hasOwnProperty.call(u, k)) || !types[i](uk)) {
          return false
        }
      }
      return true
    }
    return false
  }
}

const UnknownList = Array.isArray

const List = item => u => UnknownList(u) && u.every(item)

const Literal = value => u => u === value

const Keyof = keys => u => string(u) && hasOwnProperty.call(keys, u)

const Partial = props => {
  const keys = Object.keys(props)
  const len = keys.length
  return u => {
    if (UnknownStruct(u)) {
      for (let i = 0; i < len; i++) {
        const k = keys[i]
        const uk = u[k]
        if (uk !== undefined && !props[k](uk)) {
          return false
        }
      }
      return true
    }
    return false
  }
}

const Intersection = codecs => u => codecs.every(type => type(u))

const Union = codecs => u => codecs.some(type => type(u))

const Tuple = codecs => {
  const len = codecs.length
  return u => UnknownList(u) && u.length === len && codecs.every((type, i) => type(u[i]))
}

module.exports = {
  MIN_SAFE_INTEGER,
  MAX_SAFE_INTEGER,
  String: string,
  Number,
  Integer,
  UnknownStruct,
  NonEmptyString,
  Nil,
  Undefined,
  Boolean,
  Struct,
  UnknownList,
  List,
  Literal,
  Keyof,
  Partial,
  Intersection,
  Union,
  Tuple,
}
