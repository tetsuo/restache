const t = require('tap')
const ReactDOMServer = require('react-dom/server')
const decode = require('./lib/decode')
const createComponent = require('./lib/render')
const sort = require('./lib/sort')
const main = require('./lib/main')

const isObject = o => typeof o === 'object' && o !== null

const deepEqual = (obj1, obj2) => {
  if (obj1 === obj2) {
    return true
  } else if (isObject(obj1) && isObject(obj2)) {
    if (Object.keys(obj1).length !== Object.keys(obj2).length) {
      return false
    }
    let prop
    for (prop in obj1) {
      if (!deepEqual(obj1[prop], obj2[prop])) {
        return false
      }
    }
    return true
  }
}

const run = (fn, suite) => tc => {
  let res, err
  try {
    res = fn(tc.input)
  } catch (e) {
    err = e
  }
  if (tc.expectedErr) {
    if (err && err instanceof TypeError && err.message.startsWith(tc.expectedErr)) {
      t.pass(`${suite}: ${tc.desc}`)
    } else {
      t.fail(
        `${suite}: ${tc.desc}: expected the error message to start with "${tc.expectedErr}", got "${
          (err || {}).message
        }"`
      )
    }
    return
  } else {
    if (err) {
      t.fail(`${suite}: ${tc.desc}: ${(err || {}).message}`)
      return
    }
  }
  if (deepEqual(res, tc.expected)) {
    t.pass(`${suite}: ${tc.desc}`)
  } else {
    t.fail(`${suite}: ${tc.desc}:\nexpected:\n${JSON.stringify(tc.expected)}\ngot:\n${JSON.stringify(res)}`)
  }
}

const runDecodeTest = run(decode, 'decode')

const runRenderTest = run(
  d => ReactDOMServer.renderToStaticMarkup(createComponent(d.trees, d.opts)(d.props, d.children)),
  'render'
)

const runMainTest = run(d => ReactDOMServer.renderToStaticMarkup(main(d.roots, d.opts)(d.props, d.children)), 'main')

const runSortTest = run(d => sort(d), 'sort')

module.exports = {
  runDecodeTest,
  runSortTest,
  runRenderTest,
  runMainTest,
}
