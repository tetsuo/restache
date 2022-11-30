# stache

mustache parser.

![Build](https://github.com/onur1/stache/actions/workflows/main.yml/badge.svg)

[Report a bug or suggest a feature](https://github.com/onur1/stache/issues)

## Usage

To parse files that have the `.stache` extension under `some/folder`:

```
stache some/folder/*.stache
```

The output is an array of JSON-encoded syntax trees which will contain top-level element objects with filenames, without the extension part, as element names (e.g. the contents of `sidebar.stache` becomes the `<sidebar>` element).

## Syntax

Supports the [mustache spec](https://mustache.github.io/mustache.5.html) with the exception of partials. Partials are not needed in this implementation, simply because you can include another file as an element, the parse result will be sorted accordingly.

### Elements

Template:

```html
<h1 class="header">
  <b>Eat</b> fruits
</h1>
```

Output:

```json
[
  "h1",
  {
    "class": ["header"]
  },
  [
    [
      "b", {}, ["Eat"]
    ],
    " fruits"
  ]
]
```

### Variables

Template:

```html
{name}

<h1 id={foo} class="{bar}-small">
  {label}
</h1>
```

Output:

```json
[3, "name"]

[
  "h1",
  {
    "class": [
      [3, "bar"],
      "-small"
    ],
    "id": [
      [3, "foo"]
    ]
  },
  [
    [3, "label"]
  ]
]
```

### Sections

Template:

```html
{#items}
  <li>{name}</li>
{/items}

{#apple : some arbitrary text}
  {energy}
{/apple}
```

Output:

```json
[
  2,
  "items",
  [
    [
      "li",
      {},
      [
        [3, "name"]
      ]
    ]
  ],
  ""
]

[
  2,
  "apple",
  [
    [3, "energy"]
  ],
  "some arbitrary text"
]
```

### Inverted sections

Template:

```html
<div>
  {^fruits}
    No fruits :(
  {/fruits}
</div>
```

Output:

```json
[
  "div",
  {},
  [
    [
      4,
      "fruits",
      [
        "No fruits :("
      ],
      ""
    ]
  ]
]
```

### Comments

Template:

```html
<h1>Hello{! some comment }.</h1>
```

Output:

```json
[
  "h1",
  {},
  [
    "Hello",
    [
      5,
      " some comment "
    ],
    "."
  ]
]
```
