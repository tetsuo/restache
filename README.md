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
