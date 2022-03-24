const { has } = require('./function')

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
    if (!has(graph, id)) {
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
    sorted: sorted.filter(id => !has(recursive, id)),
    recursive,
  }
}

module.exports = {
  tsort,
  Vertex,
}
