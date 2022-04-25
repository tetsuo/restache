const { INVERTED_SECTION, SECTION, ELEMENT, VARIABLE, COMMENT, TEXT } = require('./tree')

const hasOwnProperty = Object.prototype.hasOwnProperty

const Vertex = (id, afters) => ({ id, afters })

/** topological sort */
const tsort = graph => {
  const sorted = []
  const visited = {}
  const recursive = {}
  Object.keys(graph).forEach(function visit(id, ancestors) {
    if (visited[id]) {
      return
    }
    if (!hasOwnProperty.call(graph, id)) {
      return
    }
    const vertex = graph[id]
    if (!Array.isArray(ancestors)) {
      ancestors = []
    }
    ancestors.push(id)
    visited[id] = true
    vertex.afters.forEach(afterId => {
      if (ancestors.indexOf(afterId) >= 0) {
        recursive[id] = true
        recursive[afterId] = true
      } else {
        visit(afterId, ancestors.slice())
      }
    })
    sorted.unshift(id)
  })
  return {
    sorted: sorted.filter(id => !hasOwnProperty.call(recursive, id)),
    recursive,
  }
}

const uniq = vals => {
  let r = []
  let h = {}
  vals.forEach(v => {
    if (!hasOwnProperty.call(h, v)) {
      r.push(v)
    }
    h[v] = true
  })
  return r
}

const getDependencies = e => {
  switch (e.value._type) {
    case INVERTED_SECTION:
    case SECTION:
      return e.forest.flatMap(getDependencies)
    case ELEMENT:
      return [e.value.name].concat(...e.forest.flatMap(getDependencies))
    case VARIABLE:
    case COMMENT:
    case TEXT:
      return []
  }
}
const getDependencyGraph = trees => {
  const graph = {}
  let deps, name
  trees.forEach(({ forest, value }) => {
    name = value.name
    deps = forest.flatMap(getDependencies)
    graph[name] = Vertex(name, deps.length > 1 ? uniq(deps) : deps)
  })
  return graph
}

const sort = roots => {
  const graph = getDependencyGraph(roots)
  const { sorted } = tsort(graph)
  const o = {}
  roots.forEach(d => {
    o[d.value.name] = d
  })
  return sorted.reverse().map(name => o[name])
}

module.exports = sort
