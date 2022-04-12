const hasOwnProperty = Object.prototype.hasOwnProperty

const t = require('@onur1/t')

const constant = a => () => a

const constantNull = constant(null)

const getRenderChildren = componentTree => (val, i) => {
  let j = 0
  let c = []
  let o
  for (; j < componentTree.length; ++j) {
    o = componentTree[j](val, t.Number(i) ? i : j)
    if (!t.Nil(o)) {
      if (t.String(o) && c.length > 0 && t.String(c[c.length - 1])) {
        // combine with the previous if both components output text
        c[c.length - 1] += o
      } else {
        c.push(o)
      }
    }
  }
  return c
}

const getRenderProps = props => {
  let i, o, c, p, n
  return s => {
    p = {}
    i = 0
    for (; i < props.length; i++) {
      o = props[i]
      n = o[0]
      c = o[1].map(f => f(s))
      if (c.length > 0) {
        if (c.length > 1) {
          p[n] = [c.map(String).join('')]
        } else {
          p[n] = c[0]
        }
      } else {
        p[n] = true
      }
    }
    return p
  }
}

const createVariableComponent = (v, inProp) => s => inProp || v.name === 'children' ? s[v.name] : String(s[v.name])

const createSectionComponent = (name, renderChildren) => {
  let val
  return s => {
    val = s[name]
    if (t.UnknownList(val) && val.length > 0) {
      return val.flatMap(renderChildren)
    } else if (t.UnknownStruct(val)) {
      return renderChildren(val)
    }
    if ((t.Boolean(val) && val === false) || t.Undefined(val) || t.Nil(val)) {
      return []
    }
    return renderChildren(s)
  }
}

const createInvertedSectionComponent = (name, renderChildren) => {
  let val
  return s => {
    if (!hasOwnProperty.call(s, name)) {
      return renderChildren(s)
    }
    val = s[name]
    if (
      t.Undefined(val) ||
      t.Nil(val) ||
      (t.Boolean(val) && val === false) ||
      (t.UnknownList(val) && val.length === 0)
    ) {
      return renderChildren(s)
    }
    return []
  }
}

const createTextComponent = t => constant(String(t.text))

const createPropertyComponent = (propName, propChildren, opts, index, inputTag) => {
  switch (propName) {
    case 'class':
      propName = 'className'
      break
    case 'for':
      propName = 'htmlFor'
      break
    case 'checked':
      propName = inputTag ? 'defaultChecked' : propName
      break
    case 'value':
      propName = inputTag ? 'defaultValue' : propName
      break
    default:
      if (hasOwnProperty.call(opts.syntheticEvents, propName)) {
        propName = opts.syntheticEvents[propName]
      }
  }
  return [propName, propChildren.map(c => createComponent(c, opts, index, true))]
}

const createElementComponent = (e, opts, index) => {
  const external = hasOwnProperty.call(opts.externs, e.name)
  const renderProps = getRenderProps(
    Object.entries(e.props).map(([propName, propChildren]) =>
      createPropertyComponent(
        external && hasOwnProperty.call(opts.externProps, propName) ? opts.externProps[propName] : propName,
        propChildren,
        opts,
        index,
        e.name === 'input' || e.name === 'select' || e.name === 'textarea' // has default value?
      )
    )
  )
  let renderChildren
  if (!hasOwnProperty.call(opts.selfClosingTags, e.name)) {
    renderChildren = getRenderChildren(e.children.map(c => createComponent(c, opts, index, false)))
  } else {
    renderChildren = constantNull
  }
  let p, c
  if (external) {
    return (s, key) => opts.createElement(e.name, { ...renderProps(s), ...{ key } }, renderChildren(s))
  }
  // provided by user?
  if (hasOwnProperty.call(opts.registry, e.name)) {
    return (s, key) => {
      p = renderProps(s)
      p.key = key
      c = renderChildren(p)
      if (c.length > 0) {
        p.children = c
      }
      return opts.registry[e.name](p)
    }
  }
  // own component?
  if (hasOwnProperty.call(index, e.name)) {
    return (s, key) => opts.createElement(index[e.name], { ...renderProps(s), ...{ key } }, renderChildren(s))
  }
  return renderChildren
}

const createComponent = (node, opts, index, inProp) => {
  switch (node._tag) {
    case 'Element':
      if (inProp) {
        throw new TypeError(`${node.name}: elements are not valid as prop children`)
      }
      return createElementComponent(node, opts, index)
    case 'Variable':
      return createVariableComponent(node, inProp)
    case 'Text':
      return createTextComponent(node)
    case 'Section':
      return createSectionComponent(
        node.name,
        getRenderChildren(node.children.map(c => createComponent(c, opts, index, inProp)))
      )
    case 'InvertedSection':
      return createInvertedSectionComponent(
        node.name,
        getRenderChildren(node.children.map(c => createComponent(c, opts, index, inProp)))
      )
  }
  return constantNull
}

const render = (sorted, opts) =>
  sorted.reduce((index, node, i) => {
    const c = createComponent(node, opts, index, false)
    return i === sorted.length - 1
      ? c
      : {
          ...index,
          ...{
            [node.name]: c,
          },
        }
  }, {})

module.exports = render
