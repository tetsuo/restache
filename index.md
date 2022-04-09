---
title: Home
---


# Mustache on React
{: .fs-9 }

Just the Docs gives your documentation a jumpstart with a responsive Jekyll theme that is easily customizable and hosted on GitHub Pages.
{: .fs-6 .fw-300 }

[Get started now](#getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [View it on GitHub](https://github.com/onur1/stache){: .btn .fs-5 .mb-4 .mb-md-0 }

---


**Table of contents**

- [Getting started](#getting-started)
- [Example](#example)
- [Language](#language)


## Getting started

### Dependencies

Just the Docs is built for [Jekyll](https://jekyllrb.com), a static site generator. View the [quick start guide](https://jekyllrb.com/docs/) for more information. Just the Docs requires no special plugins and can run on GitHub Pages' standard Jekyll compiler. The [Jekyll SEO Tag plugin](https://github.com/jekyll/jekyll-seo-tag) is included by default (no need to run any special installation) to inject SEO and open graph metadata on docs pages. For information on how to configure SEO and open graph metadata visit the [Jekyll SEO Tag usage guide](https://jekyll.github.io/jekyll-seo-tag/usage/).

### Quick start: Use as a GitHub Pages remote theme

1. Add Just the Docs to your Jekyll site's `_config.yml` as a [remote theme](https://blog.github.com/2017-11-29-use-any-theme-with-github-pages/)


## Example

```html
<todo></todo>
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

