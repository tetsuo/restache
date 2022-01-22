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

Create a layout from the parse tree that you can render React components with:

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
