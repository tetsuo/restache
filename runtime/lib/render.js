const t = require('@onur1/t')
const { Element, ELEMENT, TEXT, VARIABLE, SECTION, INVERTED_SECTION, fold, Section } = require('./tree')
const styleToObject = require('./style')
const { htmlBooleanAttrs } = require('./specs')

const emptyList = []
const emptyStruct = {}
const emptyString = ''

const STYLE = 'style'
const CHILDREN = 'children'
const ROOT = 'root'
const SECTION_KEY_DELIM = '.'

const hasOwnProperty = Object.prototype.hasOwnProperty

const constant = a => () => a

const constantNull = constant(null)

const renderForest = forest => attrs => forest.map((f, i) => f({ ...attrs, key: i }))

const Component = (createElement, getComponent, inProperty) => (value, forest, registry) => {
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
            (acc, [name, forest_]) => ({
              ...acc,
              ...{
                [name]:
                  forest_.length === 0
                    ? hasOwnProperty.call(htmlBooleanAttrs, name)
                      ? true
                      : ''
                    : forest_.length === 1
                    ? forest_[0](attrs)
                    : forest_
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
        return t.NonEmptyList(val)
          ? val.flatMap((b, i) =>
              forest.map((f, j) => f({ ...b, key: (i % val.length) + SECTION_KEY_DELIM + (j % forest.length) }))
            )
          : t.UnknownStruct(val)
          ? forest.map((f, i) => f({ ...attrs[value.name], key: i }))
          : t.Falsy(val)
          ? emptyList
          : renderForest(forest)(attrs)
      }
    }
    case INVERTED_SECTION: {
      return attrs => {
        const val = attrs[value.name]
        return [t.Falsy, t.EmptyList].some(f => f(val)) ? renderForest(forest)(attrs) : emptyList
      }
    }
  }
  return constantNull
}

const render = (root, createComponent, createProp, mapPropName) =>
  fold(
    root,
    {
      forest: emptyList,
      registry: root.forest.reduce((acc, a) => ({ ...acc, ...{ [a.value.name]: true } }), emptyStruct),
    },
    (a, acc, i, len, level) => {
      if (level === 0) {
        if (acc.forest.length > 0) {
          return acc.forest
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
                    constant(
                      constant(
                        styleToObject(
                          forest
                            .filter(({ value }) => value._type === TEXT)
                            .map(({ value }) => value.text)
                            .reduce((acc, x) => acc + x, emptyString)
                        )
                      )
                    )(),
                  ]
                : render(Section('', forest), createProp),
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

module.exports = (forest, opts) => {
  const components = render(
    Element(ROOT, emptyStruct, forest),
    Component(opts.createElement, opts.getComponent),
    Component(undefined, undefined, true),
    opts.mapPropName
  )
  return components[components.length - 1]
}
