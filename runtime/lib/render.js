const t = require('@onur1/t')
const { Element, ELEMENT, TEXT, VARIABLE, SECTION, INVERTED_SECTION } = require('./tree')
const styleToObject = require('./style')

const EmptyList = val => t.UnknownList(val) && val.length === 0
const NonEmptyList = val => t.UnknownList(val) && val.length > 0
const False = val => val === false
const Falsy = val => [t.Undefined, t.Nil, False].some(f => f(val))

const emptyList = []
const emptyString = ''

const STYLE = 'style'
const ROOT = 'root'
const SECTION_KEY_DELIM = '.'

const data = Symbol()

const hasOwnProperty = Object.prototype.hasOwnProperty

const constant = a => () => a
const constantNull = constant(null)

const fold = (tree, b, f, i = 0, level = 0) => {
  let r = b
  const len = tree.forest.length
  for (let j = 0; j < len; j++) {
    r = fold(tree.forest[j], r, f, j, level + 1)
  }
  return f(tree.value, r, i, len, level)
}

const Component = (value, forest, createElement, registry, getComponent, mapPropName) => {
  switch (value._type) {
    case ELEMENT: {
      const name = value.name
      let component
      if (hasOwnProperty.call(registry, name)) {
        if (registry[name] === true) {
          return attrs => forest.map((f, i) => f({ key: i, [data]: attrs[data] }))
        } else {
          component = registry[name]
        }
      } else {
        component = (getComponent && getComponent(name)) || name
      }
      const propTrees = Object.entries(value.props).map(([name, forest]) => [
        mapPropName ? mapPropName(name, value.name) : name,
        name === STYLE
          ? [
              (() => {
                const s = styleToObject(
                  forest
                    .filter(({ value }) => value._type === TEXT)
                    .map(({ value }) => value.text)
                    .reduce((acc, x) => acc + x, emptyString)
                )
                return () => s
              })(),
            ]
          : forest.map(({ value, forest }) => Component(value, forest)),
      ])
      return attrs => {
        const props = {
          key: attrs.key,
          ...propTrees.reduce(
            (acc, [name, forest]) => ({
              ...acc,
              ...{
                [name]:
                  forest.length === 0
                    ? true
                    : forest.length === 1
                    ? forest.map(f => f({ [data]: attrs[data] }))[0]
                    : forest
                        .map(f => f({ [data]: attrs[data] }))
                        .map(String)
                        .join(emptyString),
              },
            }),
            {}
          ),
        }
        return forest.length > 0
          ? createElement(
              component,
              props,
              forest.map((f, i) => f({ key: i, [data]: attrs[data] }))
            )
          : createElement(component, props)
      }
    }
    case TEXT:
      return constant(value.text.trim().length > 0 ? value.text : null)
    case VARIABLE:
      return attrs => String(attrs[data][value.name])
    case SECTION: {
      return attrs => {
        const val = attrs[data][value.name]
        return NonEmptyList(val)
          ? val.flatMap((b, i) =>
              forest.map((f, j) => f({ key: (i % val.length) + SECTION_KEY_DELIM + (j % forest.length), [data]: b }))
            )
          : t.UnknownStruct(val)
          ? forest.map((f, i) => f({ key: i, [data]: attrs[data][value.name] }))
          : Falsy(val)
          ? emptyList
          : forest.map((f, i) => f({ key: i, [data]: attrs[data] }))
      }
    }
    case INVERTED_SECTION: {
      return attrs => {
        const val = attrs[data][value.name]
        return [Falsy, EmptyList].some(f => f(val))
          ? forest.map((f, i) => f({ key: i, [data]: attrs[data] }))
          : emptyList
      }
    }
  }

  return constantNull
}

const init = root => ({
  forest: emptyList,
  registry: root.forest.reduce((acc, a) => ({ ...acc, ...{ [a.value.name]: true } }), {}),
})

const render = (root, opts) =>
  fold(root, init(root), (a, acc, i, len, level) => {
    if (level === 0) {
      if (acc.forest.length > 0) {
        return attrs => acc.forest[acc.forest.length - 1]({ key: i, [data]: attrs })
      }
      return constantNull
    }

    const forest =
      len > 0
        ? acc.forest
            .slice(0, -len)
            .concat(
              Component(
                a,
                acc.forest.slice(-len),
                opts.createElement,
                acc.registry,
                opts.getComponent,
                opts.mapPropName
              )
            )
        : acc.forest.concat(
            Component(a, emptyList, opts.createElement, acc.registry, opts.getComponent, opts.mapPropName)
          )

    return {
      forest,
      registry: level === 1 ? { ...acc.registry, ...{ [a.name]: forest[forest.length - 1] } } : acc.registry,
    }
  })

module.exports = (forest, opts) => render(Element(ROOT, {}, forest), opts)
