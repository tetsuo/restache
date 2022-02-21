const hasOwnProperty = Object.prototype.hasOwnProperty

const selfClosingTags = [
  'area',
  'base',
  'br',
  'col',
  'embed',
  'hr',
  'img',
  'input',
  'keygen',
  'link',
  'menuitem', // see: https://github.com/facebook/react/blob/85dcbf83/src/renderers/dom/shared/ReactDOMComponent.js#L437
  'meta',
  'param',
  'source',
  'track',
  'wbr',
]

const syntheticEvents = [
  'onTransitionEnd',
  'onAnimationStart',
  'onAnimationEnd',
  'onAnimationIteration',
  'onLoad',
  'onError',
  'onAbort',
  'onCanPlay',
  'onCanPlayThrough',
  'onDurationChange',
  'onEmptied',
  'onEncrypted',
  'onEnded',
  'onError',
  'onLoadedData',
  'onLoadedMetadata',
  'onLoadStart',
  'onPause',
  'onPlay',
  'onPlaying',
  'onProgress',
  'onRateChange',
  'onSeeked',
  'onSeeking',
  'onStalled',
  'onSuspend',
  'onTimeUpdate',
  'onVolumeChange',
  'onWaiting',
  'onWheel',
  'onScroll',
  'onTouchCancel',
  'onTouchEnd',
  'onTouchMove',
  'onTouchStart',
  'onSelect',
  'onClick',
  'onContextMenu',
  'onDoubleClick',
  'onDrag',
  'onDragEnd',
  'onDragEnter',
  'onDragExit',
  'onDragLeave',
  'onDragOver',
  'onDragStart',
  'onDrop',
  'onMouseDown',
  'onMouseEnter',
  'onMouseLeave',
  'onMouseMove',
  'onMouseOut',
  'onMouseOver',
  'onMouseUp',
  'onChange',
  'onInput',
  'onSubmit',
  'onFocus',
  'onBlur',
  'onKeyDown',
  'onKeyPress',
  'onKeyUp',
  'onCompositionEnd',
  'onCompositionStart',
  'onCompositionUpdate',
  'onCopy',
  'onCut',
  'onPaste',
]

const lowerCaseSyntheticEvents = syntheticEvents.map(d => d.toLowerCase())

/**
 * A `Layout` will either instantiate a new `Template` or use the provided
 * one in `ctor` parameter, and call the `Template.render` method with this value.
 */
export type Layout<U> = (state: any, options: RenderOptions<U>, ctor?: TemplateConstructor<U>) => U[]

export enum ParseTreeKind {
  Section = 2,
  Variable = 3,
  InvertedSection = 4,
  Comment = 5,
}

export enum ParseTreeIndex {
  Kind = 0,
  Tag = 0,
  Attrs = 1,
  Variable = 1,
  Children = 2,
  Parent = 3,
}

export type ParseTree =
  | {
      0: ParseTreeKind | string /* tag name */
      1: ({ [s: string]: ParseTree[] } | null) /* props */ | string /* node ref (section/variable names) */
      2?: ParseTree[] /* children */
      3?: ParseTree /* parent */
    }
  | string

export interface VisitorOptions<U> {
  attributes?: { [s: string]: any } | null
  parseTree?: ParseTree
  traverseChildren?: (children: (U | string)[] | null, node: ParseTree | string, stack?: any[]) => U[]
  state?: any
}

export interface RenderOptions<P> {
  createElement: React.Factory<P>
  registry?: { [s: string]: any }
}

export interface TemplateOptions<U> {
  visitNode: (type: string, options: VisitorOptions<U>, children: U[]) => U
}

export interface TemplateInterface<U> {
  root: ParseTree
  options?: TemplateOptions<U>
  render: (state: any, options: TemplateOptions<U>) => U[]
}

/**
 * @hidden
 */
interface TemplateConstructor<U> {
  new (options: TemplateOptions<U>): TemplateInterface<U>
}

/**
 * Template is the runtime interpreter for `ParseTree`s.
 *
 * It exposes a single `Template.render` method which, when called, will call the provided
 * `createElement` method for the each tag it sees as it traverses the parse-tree.
 *
 * @typeparam U  A generic type for the returned values from `options.createElement`.
 */
export class Template<U> implements TemplateInterface<U> {
  root: ParseTree

  options?: TemplateOptions<U>

  constructor(options?: TemplateOptions<U>) {
    this.options = options
  }

  render(state: any, options?: TemplateOptions<U>): U[] {
    this.options = { ...this.options, ...options }
    return this.traverse(null, this.root, [state])
  }

  protected traverse: (children: (U | string)[] | null, node: ParseTree | string, stack?: any[]) => U[] = (
    children,
    node,
    stack?
  ) => {
    const { visitNode } = this.options

    if (!Array.isArray(node)) {
      if (!children) {
        throw new Error('top-level text :' + node)
      }
      children.push(node as string)
    } else if ('string' === typeof node[ParseTreeIndex.Tag]) {
      let right: any[] = []
      this._reduceTree(node, stack, right)

      let propAttrs = node[ParseTreeIndex.Attrs]

      let element: any
      const tag = node[ParseTreeIndex.Tag]
      const visitorOptions = {
        ...this._formatAttributes(propAttrs as { [s: string]: any }),
        ...{
          parseTree: node,
          state: stack,
          traverseChildren: this.traverse,
        },
      }

      if ('function' === typeof visitNode) {
        element = visitNode(tag as string, visitorOptions, right)
      } else {
        element = { tag, props: propAttrs, children: right } // XXX:
      }

      if (!children) {
        return element
      } else {
        children.push(element)
      }
    } else if (
      [ParseTreeKind.Section, ParseTreeKind.InvertedSection].indexOf(node[ParseTreeIndex.Kind] as number) !== -1
    ) {
      const tail = stack[stack.length - 1]
      const value = tail[node[ParseTreeIndex.Variable] as string]

      const notInverted = (node[ParseTreeIndex.Kind] as number) !== ParseTreeKind.InvertedSection

      if (isBoolean(value)) {
        if ((notInverted && !value) || (!notInverted && value)) {
          return children
        }
        stack.push(tail)
        this._reduceTree(node, stack, children)
        stack.pop()
      } else if (Array.isArray(value)) {
        if ((notInverted && value.length === 0) || (!notInverted && value.length > 0)) {
          return children
        }
        stack.push(value)
        if (notInverted) {
          value.forEach(item => {
            stack.push(item)
            this._reduceTree(node, stack, children)
            stack.pop()
          })
        } else {
          this._reduceTree(node, stack, children)
        }

        stack.pop()
      } else if (isObject(value)) {
        if (!notInverted) {
          return children
        }
        stack.push(value)
        this._reduceTree(node, stack, children)
        stack.pop()
      } else if (typeof value === 'string' || typeof value === 'number') {
        throw new Error('invalid section type: string')
      } else {
        throw new Error('could not determine section type')
      }
    } else if (node[ParseTreeIndex.Kind] === ParseTreeKind.Variable) {
      children.push(stack[stack.length - 1][node[ParseTreeIndex.Variable] as string])
    }

    return children
  }

  private _reduceTree = (node: ParseTree | string, stack: any[], children: (U | string)[]) => {
    ;(node[ParseTreeIndex.Children] as ParseTree[]).reduce((acc, child) => {
      return this.traverse(acc, child, stack)
    }, children)
  }

  private _formatAttributes(attrs: { [s: string]: any }): {
    attributes?: { [s: string]: any }
  } {
    if (typeof attrs === 'object' && attrs !== null && Object.keys(attrs).length) {
      return { attributes: attrs }
    }
    return {}
  }
}

/**
 * Builds a `ReactElement`.
 *
 * Normalizes property names for React and ensures that we have a valid `ReactElement`.
 *
 * @param options
 * @param visitorOptions
 */
function visitObserver<T>(
  options: RenderOptions<T>,
  visitorOptions: VisitorOptions<React.ReactElement<T>>
): React.ReactElement<T> {
  const { createElement } = options
  const { parseTree, traverseChildren, state } = visitorOptions

  let type = parseTree[ParseTreeIndex.Tag] as any

  if (options.registry && hasOwnProperty.call(options.registry, type)) {
    type = options.registry[type]
  }

  let children: React.ReactNode[] = null
  let visitChildren = selfClosingTags.indexOf(type) === -1

  let attrs = visitorOptions.attributes
  if (attrs === null || typeof attrs !== 'object') {
    attrs = {}
  }

  let newAttrs = attrs as any

  Object.keys(attrs).forEach(propKey => {
    const tpl = new Template<{ [s: string]: any }>({
      visitNode: (_propType, propOptions, _propChildren) => {
        const propParseTree = propOptions.parseTree
        const traverseFn = propOptions.traverseChildren
        const newPropChildren = (propParseTree[ParseTreeIndex.Children] as (ParseTree | string)[]).reduce<
          { [s: string]: any }[]
        >((acc: ({ [s: string]: any } | string)[], childTree: ParseTree | string) => {
          return traverseFn(acc, childTree, [state])
        }, [])

        let propValue = newPropChildren.join('') as any

        if (propValue.length) {
          if (hasOwnProperty.call(state, propValue) && typeof state[propValue] === 'function') {
            propValue = state[propValue]
          }
          return {
            [propKey]: propValue,
          }
        } else {
          return {
            [propKey]: null,
          }
        }
      },
    })
    tpl.root = ['root', null, attrs[propKey]]
    newAttrs = { ...newAttrs, ...tpl.render(state) }
  })

  const normalizedProps = Object.keys(newAttrs).reduce((acc: any, key) => {
    let value: any = newAttrs[key]

    let isInput = false

    if (typeof type === 'string' && (type === 'input' || type === 'textarea')) {
      if (hasOwnProperty.call(newAttrs, 'checked')) {
        acc.defaultChecked = newAttrs.checked && newAttrs.checked === 'true'
      }
      if (hasOwnProperty.call(newAttrs, 'value')) {
        acc.defaultValue = newAttrs.value
        if (type === 'textarea') {
          visitChildren = false
        }
      }

      isInput = true
    }

    if (key === 'class') {
      acc.className = value
    } else if (key === 'for') {
      acc.htmlFor = value
    } else {
      const eventNameIndex = lowerCaseSyntheticEvents.indexOf(key)
      if (eventNameIndex !== -1) {
        acc[syntheticEvents[eventNameIndex]] = value
      } else {
        if (isInput) {
          if (['checked', 'value'].indexOf(key) === -1) {
            acc[key] = value
          }
        } else {
          const asInt = parseInt(value, 10)
          acc[key] = isNaN(asInt) ? value : asInt
        }
      }
    }

    return acc
  }, {}) as { [s: string]: any }

  if (visitChildren) {
    children = (parseTree[ParseTreeIndex.Children] as (ParseTree | string)[]).reduce<React.ReactElement<any>[]>(
      (acc: React.ReactElement<any>[], childTree: ParseTree | string) => traverseChildren(acc, childTree, [state]),
      []
    )
  }

  return createElement(
    type,
    typeof type === 'function'
      ? {
          state,
          ...normalizedProps,
        } /* state is passed as 'state' prop to custom components only */
      : normalizedProps,
    children
  )
}

/**
 * Renders a `Layout` into a `ReactElement`.
 *
 * If a tag name doesn't resolve to a `ComponentClass` in the provided `options.registry`,
 * assumes that it is an ordinary HTML tag and wrap it with `observer`. If a `ComponentClass`
 * is found instead, then it's up to the provider to make this an observer, or not.
 *
 * @param ctx  A `Layout` created with `createLayout`.
 * @param state
 * @param options
 */
export function createElement<P>(
  layout: Layout<P>,
  state: { [s: string]: any },
  options: RenderOptions<P>
): React.ReactElement<P> {
  if (!options || (typeof options === 'object' && !hasOwnProperty.call(options, 'createElement'))) {
    throw new Error("you must provide 'createElement'")
  }

  const { createElement } = options

  const registry = {
    ObserverComponent: visitObserver.bind(null, options),
  }

  if (options.registry) {
    Object.keys(options.registry).forEach(type => {
      registry[type] = visitObserver.bind(null, options)
    })
  }

  const element = layout.call(new Template(), state, {
    ...({
      registry: registry,
      visitNode: _visitNode,
    } as TemplateOptions<any>),
    ...options,
  })

  function _visitNode(
    type: string,
    options: VisitorOptions<React.ReactElement<any>>,
    children: (React.ReactNode | string)[]
  ) {
    return createElement(
      hasOwnProperty.call(registry, type) ? registry[type] : registry.ObserverComponent,
      {
        ...options,
        ...{
          state: options.state[options.state.length - 1],
          key: Math.random().toString(16).split('0.')[1],
        },
      },
      children
    )
  }

  return element
}

export const createLayout = <U>(tree: ParseTree) =>
  new Function(
    'd' /* state */,
    'm' /* options */,
    't' /* opt template inst */,
    't=t?new t:this;t.root=' + JSON.stringify(tree) + ';return t.render(d,m)'
  ) as Layout<U>

function isObject(value: unknown): value is object {
  return value !== null && typeof value === 'object'
}

function isBoolean(value: unknown): value is boolean {
  const valueType = typeof value
  return valueType === 'undefined' || valueType === 'boolean' || (valueType === 'object' && value === null)
}
