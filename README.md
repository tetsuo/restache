# restache

restache extends HTML5 syntax with [Mustache](https://mustache.github.io/)-like primitives and compiles to modern [JSX](https://facebook.github.io/jsx/), so you can write React components like it's 2013.

### Example

```html
<ul>
  {#fruits}
    <li>{name}</li>
  {/fruits}
</ul>
```

This becomes:

```jsx
<ul>
  {$0.fruits.map(($1) => (
    <li key={$1.key}>{$1.name}</li>
  ))}
</ul>
```

## Usage

Ships with an [ESBuild](https://esbuild.github.io/) plugin that you can drop straight into the `Plugins` field of your configuration:

```go
package main

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/tetsuo/restache"
)

func main() {
	options := api.BuildOptions{
		ResolveExtensions: []string{".jsx", ".tsx", ".html", ".js", ".mjs", ".ts"},
		Plugins: []api.Plugin{
			restache.Plugin(
				restache.WithExtensionName(".html"),
				restache.WithTagPrefixes(map[string]string{
					"mui":   "@mui/material",
					"icons": "@mui/icons-material",
				}),
			),
		},
		Format: api.FormatESModule,
		JSX:    api.JSXTransform,
	}

	// run build with options...
	api.Build(options)
}
```

Check out [`tetsuo/dashboard`](https://github.com/tetsuo/dashboard) for a complete working setup and a few basic components.

There's currently no support for Node.JS environment, but planned.

## Language syntax

### Variables

Variables provide access to data in the current scope using dot notation.

Accessing component props:

```html
<article class="fruit-card" data-fruit-id={id}>
  <h3>{name}, eaten at {ateAt}</h3>
  <p>{color} on the outside, {flavor.profile} on the inside</p>
  <img src={image.src} alt={image.altText} />
</article>
```

Output:

```jsx
<article className="fruit-card" data-fruit-id={ $0.id }>
  <h3>{$0.name}, eaten at {$0.ateAt}</h3>
  <p>{$0.color} on the outside, {$0.flavor.profile} on the inside</p>
  <img src={ $0.image.src } alt={ $0.image.altText } />
</article>
```

They can only appear within text nodes or as full attribute values inside a tag.

```
✅ can insert variable {here}, or <img href={here}>.
```

### When

Renders block when expression is truthy:

```html
{?loggedIn}
  <welcome-banner user={user} />
{/loggedIn}
```

Output:

```jsx
($0.loggedIn && <WelcomeBanner user={ $0.user }></WelcomeBanner>)
```

### Unless

Renders block when expression is falsy:

```html
{^hasPermission}
  <p>You do not have access to this section.</p>
{/hasPermission}
```

Output:

```jsx
(!$0.hasPermission && <p>You do not have access to this section.</p>)
```

### Range

Iterates over a list value:

```html
<ul>
  {#fruits}
    <li>
      {name}
      <ul>
        {#vitamins}
          <li>{name}</li>
        {/vitamins}
      </ul>
    </li>
  {/fruits}
</ul>
```

Output:

```jsx
<ul>
  {$0.fruits.map(($1) => (
    <li key={$1.key}>
      {$1.name}
      <ul>
        {$1.vitamins.map(($2) => (
          <li key={$2.key}>{$2.name}</li>
        ))}
      </ul>
    </li>
  ))}
</ul>
```

**Range blocks** create a new lexical scope. Inside the block, `{name}` refers to the local object in context; outer scope variables are not accessible.

> ⚠️ Control structures must wrap well-formed elements (or other well-formed control constructs), and cannot appear inside tags.

### Comments

There are two types of comments:

```html
<!-- Comment example -->

<span>{! TODO: fix bugs } hi</span>
```

Output:

```jsx
<span>{ /* TODO: fix bugs */ } hi</span>
```

Standard HTML comments are removed from the generated output.

restache comments compile into JSX comments.

## JSX generation

Below are some of the JSX-specific quirks that restache handles where necessary.

### Fragment wrapping

Multiple root elements are wrapped in a Fragment:

```html
<h1>{title}</h1>
<p>{description}</p>
```

Output:

```jsx
<>
  <h1>{$0.title}</h1>
  <p>{$0.description}</p>
</>
```

- This also applies within control blocks.
- If you only have one root element, then a Fragment is omitted.

### Case conversion

React requires component names to start with a capital letter and prop names to use camelCase. In contrast, HTML tags and attribute names are not case-sensitive.

To ensure compatibility, restache applies the following transformations:

- Elements written in kebab-case (e.g. `<my-button>`) are automatically converted to PascalCase (`MyButton`) in the output.
- Similarly, kebab-case attributes (like `disable-padding`) are converted to camelCase (`disablePadding`).

kebab-case 🔜 React case:

```html
<search-bar
  hint-text="Type to search..."
  data-max-items="10"
  aria-label="Site search"
/>
```

Output:

```jsx
<SearchBar hintText="Type to search..." data-max-items="10" aria-label="Site search"></SearchBar>
```

ℹ️ _Attributes starting with `data-` or `aria-` are **preserved as-is**, in line with React's conventions._

### Attribute name normalization

Certain attributes are automatically renamed for React compatibility:

```html
<form enctype="multipart/form-data" accept-charset="UTF-8">
  <input name="username" popovertarget="hint">
  <textarea maxlength="200" autocapitalize="sentences"></textarea>
  <button formaction="/submit" formtarget="_blank">Submit</button>
</form>

<video controlslist="nodownload">
  <source src="video.mp4" srcset="video-480.mp4 480w, video-720.mp4 720w">
</video>
```

Output:

```jsx
<>
  <form encType="multipart/form-data" acceptCharset="UTF-8">
    <input name="username" popoverTarget="hint" />
    <textarea maxLength="200" autoCapitalize="sentences"></textarea>
    <button formAction="/submit" formTarget="_blank">
      Submit
    </button>
  </form>
  <video controlsList="nodownload">
    <source
      src="video.mp4"
      srcSet="video-480.mp4 480w, video-720.mp4 720w"
    />
  </video>
</>
```

Attribute renaming only occurs when the attribute is valid for the tag. For instance, `formaction` isn't renamed on `<img>` since it isn't valid there.

However, some attributes are renamed globally, regardless of which element they're used on. These include:

- All standard event handler attributes (`onclick`, `onchange`, etc.), which are converted to their camelCased React equivalents (e.g. `onClick`, `onChange`)
- Common HTML aliases and reserved keywords like `class` and `for`, which are renamed to `className` and `htmlFor`
- Certain accessibility- and editing-related attributes, such as `spellcheck` and `tabindex`

_See [`table.go`](https://github.com/tetsuo/restache/blob/v0.x/table.go) for the full list._

### Implicit key insertion in loops

When rendering lists, restache inserts a `key` prop automatically, assigning it to the top-level element or to a wrapping Fragment if there are multiple root elements.

Key is passed to the root element inside a loop:

```html
{#images}
  <img src={src}>
{/images}
```

Output:

```jsx
$0.images.map($1 => <img key={ $1.key } src={ $1.src } />)
```

If there are multiple roots, it goes on the Fragment:

```html
{#items}
  <h1>{title}</h1>
  <h3>{description}</h3>
{/items}
```

Output:

```jsx
$0.items.map(($1) => (
  <React.Fragment key={$1.key}>
    <h1>{$1.title}</h1>
    <h3>{$1.description}</h3>
  </React.Fragment>
))
```

Manually set key when there's single root:

```html
{#images}
  <img key={id} src={src}>
{/images}
```

Output:

```jsx
$0.images.map($1 => <img key={ $1.id } src={ $1.src } />)
```

## Importing other components

Component imports are inferred from the tag names. The following examples show how different components are resolved:

| HTML               | JSX              | Import path                |
| ------------------ | ---------------- | -------------------------- |
| `<my-button>`      | `<MyButton>`     | `./MyButton`               |
| `<ui:card-header>` | `<UiCardHeader>` | `./ui/CardHeader`          |
| `<main>`           | `<main>`         | Not resolved, standard tag |
| `<ui:div>`         | `<UiDiv>`        | `./ui/div`                 |

### Component resolution

Any tag that isn't a known HTML element is treated as a component.

When the parser encounters such a tag, it follows these steps:

It first determines whether the tag uses a namespace. Namespaced tags contain a prefix and a component name (e.g., `<ui:button>`).

#### Standard custom tags

If the tag **does not** contain a namespace (e.g., `<my-button>`):

- The parser looks for an exact match in the build configuration's `tagMappings`.
- If no mapping is found, it falls back to searching in the current directory.
  For example, `<my-button>` could resolve to either `./my-button` or `./MyButton`.

#### Namespaced tags

If the tag **does** contain a namespace (e.g., `<ui:button>`):

- The parser checks the `tagPrefixes` configuration this time. If a prefix (e.g., `ui`) is defined, it uses the mapped path. For example, if `mui` is mapped to `@mui/material`, then `<mui:app-bar>` resolves to `@mui/material/AppBar`.
- If no mapping is found, it attempts to resolve the component from a subdirectory:
  e.g., `<ui:button>` then becomes one of `ui/Button.js`, `ui/button.jsx`, `ui/button.tsx`, etc.

⚠️ Standard HTML tags are not resolved as components, even if identically named files exist in the current directory.
However, **namespacing can override this behavior**. For example, `<ui:div>` will resolve to `./ui/div` (or `./ui/Div`), even though `<div>` is a native HTML element.
