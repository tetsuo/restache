const ReactDOM = require('react-dom')
const createComponent = require('@onur1/stache')

const data = {
  items: [{ name: 'Apple' }, { name: 'Orange' }, { name: 'Watermelon' }],
}

fetch('http://localhost:7882/')
  .then(response => response.json())
  .then(templates => ReactDOM.render(createComponent(templates)(data), document.getElementById('content')))
