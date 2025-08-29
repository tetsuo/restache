### Integration with React hooks

Currently, most logic must live in a `.jsx` file next to the corresponding `.html` file.

The long-term plan is to introduce a minimal set of expressions into the language itself so that common hooks like `useState`, `useContext`, and `useSelector` from Redux can be inferred from markup and compiled automatically.

### Relational and logical expressions

Support for predicates inside range blocks is planned for v1.

```
{#products: price < 100 && inStock}
  <product-card />
{/products}
```

This could be compiled as a filter, and potentially mapped to things like backend queries (e.g. MongoDB, CouchDB, Algolia, ...) as well.

### Smarter code generation

Before adding expressions, there's still a lot that can be optimized with the syntax that's already in place.

Consider pattern matching. Since restache doesn't support expressions beyond dot notation, patterns have to be represented structurally. For example, using an object with a mutually exclusive key set:

```js
{
  home: {...},
  settings: undefined,
  products: undefined
}
```

Then in the template:

```
{?home}    <home />    {/home}
{?settings}<settings />{/settings}
{?products}<products />{/products}
```

Compiles to:

```js
if (props.home) <Home />
if (props.settings) <Settings />
```

This results in an `O(n)` operation instead of an `O(1)` equality check like `switch(route)`, but the difference is negligible unless you're dealing with many conditions.

Future versions of restache will generate `if/else` or `switch` statements when keyed unions are used, along with other optimizations such as merging adjacent `{?x}` and `{^x}` blocks into a single conditional.

### More codegen targets

Plans include:

* Emitting a real JavaScript AST instead of raw JSX strings
* Supporting things other than React
* Option to emit TSX

When TypeScript is used and type information is available at build time, more advanced optimizations may be possible. But even without that, structural inference allows detecting optionals and iterables, which can be used to emit a generic type for the component.
