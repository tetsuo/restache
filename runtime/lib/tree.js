const Tree = (value, forest = []) => ({ value, forest })

const SECTION = 'Section'

const Section = (name, forest) => Tree({ _type: SECTION, name }, forest)

const INVERTED_SECTION = 'InvertedSection'

const InvertedSection = (name, forest) => Tree({ _type: INVERTED_SECTION, name }, forest)

const ELEMENT = 'Element'

const Element = (name, props = {}, forest) => Tree({ _type: ELEMENT, name, props }, forest)

const VARIABLE = 'Variable'

const Variable = name => Tree({ _type: VARIABLE, name })

const TEXT = 'Text'

const Text = text => Tree({ _type: TEXT, text })

const COMMENT = 'Comment'

const Comment = comment => Tree({ _type: COMMENT, comment })

const fold = (tree, b, f, i = 0, level = 0) => {
  let r = b
  const len = tree.forest.length
  for (let j = 0; j < len; j++) {
    r = fold(tree.forest[j], r, f, j, level + 1)
  }
  return f(tree.value, r, i, len, level)
}

module.exports = {
  Tree,
  Section,
  InvertedSection,
  Element,
  Variable,
  Text,
  Comment,
  SECTION,
  INVERTED_SECTION,
  ELEMENT,
  VARIABLE,
  TEXT,
  COMMENT,
  fold,
}
