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

const UnknownList = Array.isArray

module.exports = {
  String: string,
  Number,
  Integer,
  UnknownStruct,
  NonEmptyString,
  Nil,
  Undefined,
  Boolean,
  UnknownList,
}
