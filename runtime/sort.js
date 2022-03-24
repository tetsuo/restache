const { tsort, Vertex } = require('./graph')
const { uniq, flatten, has } = require('./function')

const getDependencies = e => {
  switch (e._tag) {
    case 'InvertedSection':
    case 'Section':
      return flatten(e.children.map(getDependencies))
    case 'Element':
      if (e.external) {
        return [e.name].concat(flatten(e.children.map(getDependencies)))
      } else {
        return flatten(e.children.map(getDependencies))
      }
    case 'Variable':
    case 'Comment':
    case 'Text':
      return []
  }
}

const getElementGraph = roots => {
  const graph = {}
  let deps, name, vertex
  roots.forEach(d => {
    name = d.name
    if (has(graph, name)) {
      throw new Error(`duplicated name: ${JSON.stringify(name)}`)
    }
    vertex = graph[name] = Vertex(name, [])
    deps = flatten(d.children.map(getDependencies))
    if (deps.length > 1) {
      deps = uniq(deps)
    }
    vertex.afters.push(...deps)
  })
  return graph
}

const getElementMap = roots => {
  const map = {}
  roots.forEach(d => {
    map[d.name] = d
  })
  return map
}

const sort = roots => {
  const graph = getElementGraph(roots)
  const { sorted, recursive } = tsort(graph)
  const map = getElementMap(roots)
  return [sorted.reverse().map(name => map[name]), recursive]
}

module.exports = sort
