# stache

simple templating language.

## Example

This is the file, `hello.html`:

```html
<div>Hello, {name}!</div>
```

Running `stache < hello.html` outputs a compact AST document formatted as JSON:

```json
["div",{},["Hello, ",[3,"name"],"!"]]
```

Create a layout from this object and render React elements with it:

```js
const React = require('react')
const ReactDOM = require('react-dom')
const { createLayout, createElement } = require('@onur1/stache')

const layout = createLayout(['div', {}, ['Hello, ', [3, 'name'], '!']])
const data = { name: 'Onur' }
const opts = { createElement: React.createElement }

ReactDOM.render(
  createElement(layout, data, opts),
  document.getElementById('content')
)
```

## Language

stache supports the [mustache spec](http://mustache.github.io/mustache.5.html) with the exception of lambdas and partials.

### Variables

```
<h1>{name}</h1>
```

### Sections

```
<ul>
  {#fruits}
  <li>
    {name}
    {#vitamins}
      <span>{name}</span>
    {/vitamins}
  </li>
  {/fruits}
</ul>
```

### Inverted sections

```
<div>
  {^fruits}
    No fruits :(
  {/fruits}
</div>
```

### Comments

```
<h1>Today{! ignore me }.</h1>
```

## Custom components

Provide your own components in the `options.registry` object.

For example, this html:

```html
<foo class=bla>
  <div>{n}</div>
</foo>
```

You can render a custom `foo` component like this:

```js
const foo = ({ className, children }) =>
  <p className={ className }>
      <b>bold</b>
      { children }
  </p>

createElement(layout, { n: 1 }, { registry: { foo: foo } })
```
