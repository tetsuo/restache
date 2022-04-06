const Variable = name => ({ _tag: 'Variable', name })

const Comment = comment => ({ _tag: 'Comment', comment })

const Text = text => ({ _tag: 'Text', text })

/** tree types */

const Element = (name, props = {}, children = []) => ({
  _tag: 'Element',
  name,
  props,
  children,
})

const Section = (name, children = []) => ({
  _tag: 'Section',
  name,
  children,
})

const InvertedSection = (name, children = []) => ({
  _tag: 'InvertedSection',
  name,
  children,
})

module.exports = { Section, InvertedSection, Variable, Comment, Text, Element }
