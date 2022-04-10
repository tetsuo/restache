const React = require('react')
const ReactDOM = require('react-dom')
const createComponent = require('@onur1/stache')

const data = {
  now: new Date().toString(),
  items: [
    { name: 'Apple', items: [{ name: 'Vitamin D' }, { name: 'Vitamin C' }] },
    { name: 'Orange' },
    { name: 'Watermelon' },
  ],
}

fetch('http://localhost:7882/')
  .then(response => response.json())
  .then(templates =>
    ReactDOM.render(
      React.createElement(React.StrictMode, {}, [createComponent(templates)(data)]),
      document.getElementById('content')
    )
  )
