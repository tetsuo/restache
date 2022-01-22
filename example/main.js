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
