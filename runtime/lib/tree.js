const Tree = (value, forest = []) => ({ value, forest })

const Section = (name, forest) => Tree({ _type: 'Section', name }, forest)

const InvertedSection = (name, forest) => Tree({ _type: 'InvertedSection', name }, forest)

const Element = (name, props = {}, forest) => Tree({ _type: 'Element', name, props }, forest)

const Variable = name => Tree({ _type: 'Variable', name })

const Text = text => Tree({ _type: 'Text', text })

const Comment = comment => Tree({ _type: 'Comment', comment })

module.exports = {
  Tree,
  Section,
  InvertedSection,
  Element,
  Variable,
  Text,
  Comment,
}
