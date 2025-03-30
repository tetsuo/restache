# restache

Restache is a Mustache-inspired templating language designed for integration with React.

⚠️ **Work In Progress** ⚠️

## Example template

```html
<div>
  {#cart.items}
    <span>{name}</span>
  {/cart.items}
  {?cart.empty}
    <p>Your cart is empty</p>
  {/cart.empty}
  {! This is a comment and will not be rendered }
</div>
```

## Grammar

A simplified, informal grammar for Restache is as follows:

```
Template     ::= (Element | ControlBlock | Variable | Comment | Text)*

Element      ::= StartTag Template EndTag | SelfClosingTag

ControlBlock ::= RangeBlock | WhenBlock | UnlessBlock

RangeBlock   ::= "{#" Expression "}" Template EndControl
WhenBlock    ::= "{?" Expression "}" Template EndControl
UnlessBlock  ::= "{^" Expression "}" Template EndControl

EndControl   ::= "{/" Expression "}"

Variable     ::= "{" Expression "}"
Comment      ::= "{!" CommentText "}"

Text         ::= Any sequence of characters that does not form part of a recognized Stache token.

Expression   ::= Identifier ('.' Identifier)*
```

## Semantic rules

### Parse tree and node types

During parsing, the template is converted into a tree of nodes. The primary node types include:

- **ComponentNode**: The root node of the entire document.
- **ElementNode**: Represents HTML-like elements.
- **TextNode**: Contains textual content.
- **VariableNode**: Represents a dynamic value inserted into the template.
- **CommentNode**: Contains comments that are not rendered.
- **RangeNode, WhenNode, UnlessNode**: Represent control structures that conditionally render or iterate over content.

### Scope and context

- **Range**:
  When a `{# expression }` token is encountered, the parser splits the expression by `.` and appends the parts to the current scope (stored in the node's `Path` field). All child nodes within this block use the extended scope for variable resolution.

- **When and Unless**:
  These constructs do not modify the scope. They evaluate their expressions against the current context, rendering their children only when the condition holds true.

### Attribute parsing

- **Expression attributes**:
  If an attribute's value is entirely wrapped in `{ ... }`, it is treated as an expression and evaluated at runtime.

- **Static attributes**:
  Attribute values that contain additional text (outside of a full expression) are treated as literal strings.

