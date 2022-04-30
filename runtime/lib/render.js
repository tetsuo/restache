const t = require('@onur1/t')
const { Element, ELEMENT, TEXT, VARIABLE, SECTION, INVERTED_SECTION, fold } = require('./tree')
const styleToObject = require('./style')

const emptyList = []
const emptyStruct = {}
const emptyString = ''

const STYLE = 'style'
const CHILDREN = 'children'
const ROOT = 'root'
const SECTION_KEY_DELIM = '.'

const hasOwnProperty = Object.prototype.hasOwnProperty

const EmptyList = val => t.UnknownList(val) && val.length === 0
const NonEmptyList = val => t.UnknownList(val) && val.length > 0
const False = val => val === false
const Falsy = val => [t.Undefined, t.Nil, False].some(f => f(val))

const constant = a => () => a

const constantNull = constant(null)

const renderForest = forest => attrs => forest.map((f, i) => f({ ...attrs, key: i }))

  switch (value._type) {
    case ELEMENT: {
      const name = value.name
      let component
      if (hasOwnProperty.call(registry, name)) {
        if (registry[name] === true) {
          return renderForest(forest)
        } else {
          component = registry[name]
        }
      } else {
        component = (getComponent && getComponent(name)) || name
      }
      return attrs => {
        const props = {
          key: attrs.key,
          ...value.props.reduce(
            (acc, [name, forest]) => ({
              ...acc,
              ...{
                [name]:
                  forest.length === 0
                    ? true
                    : forest.length === 1
                    ? forest[0](attrs)
                    : forest
                        .map(f => f(attrs))
                        .map(String)
                        .join(emptyString),
              },
            }),
            emptyStruct
          ),
        }
        return forest.length > 0
          ? createElement(component, props, renderForest(forest)(attrs))
          : createElement(component, props)
      }
    }
    case TEXT:
      return constant(value.text.trim().length > 0 ? value.text : null)
    case VARIABLE:
      return inProperty || value.name === CHILDREN ? attrs => attrs[value.name] : attrs => String(attrs[value.name])
    case SECTION: {
      return attrs => {
        const val = attrs[value.name]
        return NonEmptyList(val)
          ? val.flatMap((b, i) =>
              forest.map((f, j) => f({ ...b, key: (i % val.length) + SECTION_KEY_DELIM + (j % forest.length) }))
            )
          : t.UnknownStruct(val)
          ? forest.map((f, i) => f({ ...attrs[value.name], key: i }))
          : Falsy(val)
          ? emptyList
          : renderForest(forest)(attrs)
      }
    }
    case INVERTED_SECTION: {
      return attrs => {
        const val = attrs[value.name]
        return [Falsy, EmptyList].some(f => f(val)) ? renderForest(forest)(attrs) : emptyList
      }
    }
  }
  return constantNull
}

const render = (root, createComponent, mapPropName) =>
  fold(
    root,
    {
      forest: emptyList,
      registry: root.forest.reduce((acc, a) => ({ ...acc, ...{ [a.value.name]: true } }), emptyStruct),
    },
    (a, acc, i, len, level) => {
      if (level === 0) {
        if (acc.forest.length > 0) {
          return acc.forest[acc.forest.length - 1]
        }
        return constantNull
      }

      let b = a

      if (b._type === ELEMENT) {
        b = {
          ...a,
          ...{
            props: Object.entries(b.props).map(([name, forest]) => [
              mapPropName ? mapPropName(name, b.name) : name,
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
                : forest.map(({ value, forest }) => createComponent(value, forest, undefined, true)),
            ]),
          },
        }
      }

      const forest =
        len > 0
          ? acc.forest.slice(0, -len).concat(createComponent(b, acc.forest.slice(-len), acc.registry, false))
          : acc.forest.concat(createComponent(b, emptyList, acc.registry, false))

      return {
        forest,
        registry: level === 1 ? { ...acc.registry, ...{ [b.name]: forest[forest.length - 1] } } : acc.registry,
      }
    }
  )

module.exports = (forest, opts) =>
  render(Element(ROOT, emptyStruct, forest), Component(opts.createElement, opts.getComponent), opts.mapPropName)
