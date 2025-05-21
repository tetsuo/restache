# restache

Mustache-like extension to HTML syntax, designed for use with React and transpiled into JSX.

## Getting started

### Installation

To use Restache in your project, you'll need to integrate the [ESBuild](https://esbuild.github.io/) plugin:

```bash
go get github.com/tetsuo/restache
```

### Usage

#### 1. Define a template

For example, to list fruits, create `fruits.stache`:

```html
<ul>
  {#items}
    <li>{name}</li>
  {/items}
</ul>
```

#### 2. Configure ESBuild with the plugin

```go
import (
  "github.com/evanw/esbuild/pkg/api"
  "github.com/tetsuo/restache"
)

func main() {
  api.Build(api.BuildOptions{
    EntryPoints: []string{"main.mjs"},
    Bundle:      true,
    Plugins:     []api.Plugin{restache.Plugin()},
    Outfile:     "out.js",
  })
}
```

#### 3. Import .stache files as React components

```js
import Fruits from './fruits.stache'
import { createRoot } from 'react-dom/client'

const root = createRoot(document.getElementById('root'))
root.render(Fruits({ items: [{ name: 'Apple', key: 'apple' }] }))
```

#### 4. Build your project

Run your build process, and `.stache` files will be transpiled into JSX automatically.

> A more complete usage example is available in the [example](./example) directory.

## Syntax

### Variables

Use `{variableName}` to interpolate variables.

### Conditionals

#### When

`{?isVisible}...{/isVisible}` renders content when the condition is true.

#### Unless

`{^isHidden}...{/isHidden}` renders content when the condition is false.

### Loops

`{#list}...{/list}` iterates over a list.

### Components

Define components using custom tags, which are resolved based on naming conventions and mappings.

## Component resolution

Restache resolves component tags by:

1. **PascalCasing**: Converting tag names to PascalCase.
2. **Prefix Mapping**: Using configured prefixes to locate components.
3. **Tag Mappings**: Directly mapping tag names to component paths.

These mappings are configured via the ESBuild plugin options.

