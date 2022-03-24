const { isObject, has, isFunction } = require('./function')

const Template = (name, roots) => ({ name, roots })

const isTemplate = u => isObject(u) && has(u, 'roots') && Array.isArray(u.roots) && has(u, 'name') && u.name.length > 0

const createElement = (_0, _1, _2) => [_0, _1, _2]

const createFragment = _2 => [_2]

const getCreateKey = () => {
  let key = 0
  return () => ++key
}

const isLayoutOptions = u =>
  isObject(u) &&
  has(u, 'externs') &&
  isObject(u.externs) &&
  has(u, 'selfClosingTags') &&
  isObject(u.selfClosingTags) &&
  has(u, 'registry') &&
  isObject(u.registry) &&
  isFunction(u.createElement) &&
  isFunction(u.createFragment) &&
  isFunction(u.createKey)

const LayoutOptions = u => {
  if (!isObject(u)) {
    u = {}
  }
  if (!isObject(u.externs)) {
    u.externs = {}
  }
  if (!isObject(u.selfClosingTags)) {
    u.selfClosingTags = {}
  }
  if (!isObject(u.registry)) {
    u.registry = {}
  }
  if (!isFunction(u.createElement)) {
    u.createElement = createElement
  }
  if (!isFunction(u.createFragment)) {
    u.createFragment = createFragment
  }
  if (!isFunction(u.createKey)) {
    u.createKey = getCreateKey()
  }

  return u
}

const Layout = (templates, opts) => ({ templates, opts: LayoutOptions(opts) })

const isLayout = u =>
  isObject(u) &&
  has(u, 'templates') &&
  Array.isArray(u.templates) &&
  u.templates.every(isTemplate) &&
  has(u, 'opts') &&
  isLayoutOptions(u.opts)

module.exports = {
  Template,
  Layout,
  LayoutOptions,
  isLayout,
  isTemplate,
  isLayoutOptions,
}
