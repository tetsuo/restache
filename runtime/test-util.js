const ReactDOMServer = require('react-dom/server')
const { exec } = require('child_process')
// const createComponent = require('./react')
// const { Layout, Template } = require('./layout')
const t = require('tap')
const decode = require('./lib/decode')
const createComponent = require('./lib/render')
const sort = require('./lib/sort')
// const { getDependencyGraph } = require('./lib/graph')

const parse = s =>
  new Promise((resolve, reject) =>
    exec(`echo '${s}' | ${__dirname}/../stache`, (err, stdout) => {
      if (err) {
        reject(err)
        return
      }
      let data
      try {
        data = JSON.parse(stdout)
      } catch (e) {
        reject(e)
      }
      resolve(data)
    })
  )

// async function runRenderTest(tc) {
//   const roots = await Promise.all(Object.entries(tc.input).map(([k, v]) => parse(v).then(o => [k, [o]])))
//   const T = createComponent(
//     Layout(
//       roots.map(([name, trees]) => Template(name, trees)),
//       tc.opts
//     )
//   )
//   const markup = ReactDOMServer.renderToStaticMarkup(T(tc.props || {}))
//   if (markup === tc.expected) {
//     t.pass(tc.desc)
//   } else {
//     t.fail(`expected "${tc.expected}", got "${markup}"`)
//   }
// }

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

// const runGetDependencyGraphTest = run(getDependencyGraph)

const runSortTest = run(d => sort(d)[0], 'sort')

// const runSortTest = run(input => {
//   const d = decode(input)
// })

// async function runSortTest(tc) {
//   const d = decode(tc.input)
//   Object.entries(tc.input).map(([k, v]) => decode(v).then(o => [k, [o]]))
// }

module.exports = {
  parse,
  runDecodeTest,
  runSortTest,
  runRenderTest,
}
