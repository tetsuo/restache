const sort = require('./sort')
const { decode } = require('./stache')
const { isLayout } = require('./layout')
const { constant, isBoolean, isObject, isString, flatten, has } = require('./function')
const { syntheticEventNamesLowercase, syntheticEventNames } = require('./constants')

const variableComponent = e => s => isObject(s) ? s[e.name] : null

const getRenderChildren = (e, opts) => {
  const children = e.children.map(c => toComponent(c, opts))
  let c, j, o
  return val => {
    c = []
    j = 0
    for (; j < children.length; ++j) {
      o = children[j](val)
      if (o !== null) {
        if (isString(o) && c.length > 0 && isString(c[c.length - 1])) {
          c[c.length - 1] += o
        } else {
          c.push(o)
        }
      }
    }
    return c
  }
}

const invertedSectionComponent = (e, opts) => {
  const render = getRenderChildren(e, opts)
  let val
  return s => {
    val = s[e.name]
    if (isBoolean(val)) {
      if (val) {
        return []
      }
      return render(s)
    } else if (Array.isArray(val)) {
      if (val.length > 0) {
        return []
      }
      return render(val)
    } else if (isObject(val)) {
      return []
    }
    throw new TypeError(
      `inverted section: ${e.name}: expected a boolean, array, or an object value, got ${JSON.stringify(val)}`
    )
  }
}

const sectionComponent = (e, opts) => {
  const render = getRenderChildren(e, opts)
  let val
  return s => {
    val = s[e.name]
    if (isBoolean(val)) {
      if (!val) {
        return []
      }
      return render(s)
    } else if (Array.isArray(val)) {
      if (val.length === 0) {
        return []
      }
      return flatten(val.map(render))
    } else if (isObject(val)) {
      return render(val)
    }
    throw new TypeError(`section: ${e.name}: expected a boolean, array, or an object value, got ${JSON.stringify(val)}`)
  }
}

const isElementInput = e => e.name === 'input' || e.name === 'textarea'

const getRenderProps = (e, opts) => {
  let i, o, c, p, n

  const props = Object.entries(e.props).map(([key, children]) => {
    switch (key) {
      case 'class':
        key = 'className'
        break
      case 'for':
        key = 'htmlFor'
        break
      case 'checked':
        key = isElementInput(e) ? 'defaultChecked' : key
        break
      case 'value':
        key = isElementInput(e) ? 'defaultValue' : key
        break
      default:
        if (has(syntheticEventNamesLowercase, key)) {
          key = syntheticEventNames[key]
        }
    }
    return [key, children.map(c => toComponent(c, opts))]
  })

  return s => {
    p = {}
    i = 0
    let maybeInt
    for (; i < props.length; i++) {
      o = props[i]
      n = o[0]
      c = o[1].map(f => f(s))
      if (c.length > 1) {
        p[n] = [c.filter(Boolean).join('')]
      } else if (c.length > 0) {
        // take the first node only
        maybeInt = parseInt(c[0], 10)
        p[n] = isNaN(maybeInt) ? c[0] : maybeInt
      }
    }
    return p
  }
}

const constantEmptyArray = constant([])

const constantEmptyObject = constant({})

const elementComponent = (e, opts) => {
  let renderProps, renderChildren, component
  if (!has(opts.selfClosingTags, e.name)) {
    renderChildren = getRenderChildren(e, opts)
    renderProps = getRenderProps(e, opts)
  } else {
    renderChildren = constantEmptyArray
    renderProps = constantEmptyObject
  }
  let p, c
  if (e.external && has(opts.components, e.name)) {
    return s => {
      p = renderProps(s)
      c = renderChildren(p)
      if (c.length > 0) {
        p.children = c
      }
      return opts.components[e.name](p)
    }
  }
  if (!has(opts.externs, e.name) && !has(opts.components, e.name)) {
    return s => opts.createFragment(renderChildren(s))
  } else if (!e.external && has(opts.registry, e.name)) {
    component = opts.registry[e.name]
  } else {
    component = e.name
  }
  return s => opts.createElement(component, { ...renderProps(s), ...{ key: opts.createKey() } }, renderChildren(s))
}

const textComponent = e => constant(e.text)

const commentComponent = constant(null)

const toComponent = (e, opts) => {
  switch (e._tag) {
    case 'Element':
      return elementComponent(e, opts)
    case 'Text':
      return textComponent(e)
    case 'InvertedSection':
      return invertedSectionComponent(e, opts)
    case 'Section':
      return sectionComponent(e, opts)
    case 'Variable':
      return variableComponent(e)
    case 'Comment':
      return commentComponent
    default:
      throw new TypeError(
        `expected a Text, Variable, Comment, Section, InvertedSection, or an Element, got ${JSON.stringify(e)}`
      )
  }
}

const findTopLevelProps = u => {
  let props = {}
  if (Array.isArray(u) && u.length > 0) {
    if (u[0] === 3 || u[0] === 2 || u[0] === 4) {
      if (u.length > 1) {
        props[u[1]] = [[3, u[1]]] // bound to the same variable name
      } else {
        throw new TypeError(`expected a property name in the 2nd index, got ${JSON.stringify(u)}`)
      }
    } else if (isString(u[0])) {
      if (isObject(u[1])) {
        props = flatten(
          Object.values(u[1]).map(children => {
            if (!Array.isArray(children)) {
              throw new Error(`expected an array, got ${JSON.stringify(children)}`)
            }
            return children.map(findTopLevelProps)
          })
        ).reduce((acc, x) => ({ ...acc, ...x }), props)
      }
      if (Array.isArray(u[2]) && u[2].length > 0) {
        props = u[2].reduce((acc, x) => ({ ...acc, ...findTopLevelProps(x) }), props)
      }
    }
  }
  return props
}

const templateToElement = t => {
  let i = 0
  let root
  let props = {}
  for (; i < t.roots.length; ++i) {
    root = t.roots[i]
    props = { ...props, ...findTopLevelProps(root) }
  }
  return [t.name, props, t.roots]
}

const Component = layout => {
  if (!isLayout(layout)) {
    throw new TypeError(`expected a Layout object as its first parameter, got ${JSON.stringify(layout)}`)
  }
  const [sorted, recursive] = sort(layout.templates.map(t => decode(templateToElement(t), layout.opts)))
  if (Object.keys(recursive).length) {
    throw new TypeError(`got recursive keys: ${Object.keys(recursive).join(', ')}`)
  }
  return sorted.reduce((components, e, i) => {
    const c = toComponent(e, { ...layout.opts, ...{ components } })
    return i === sorted.length - 1
      ? c
      : {
          ...components,
          ...{
            [e.name]: c,
          },
        }
  }, {})
}

module.exports = Component
