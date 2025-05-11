#!/usr/bin/env node

// Formats expected AST-like structures for parser test cases defined in parser_ast.txt.
function prettyPrint(input) {
  let i = 0;
  const indentStep = 2;

  function skipWhitespace() {
    while (i < input.length && /\s/.test(input[i])) i++;
  }

  function parseBlock() {
    skipWhitespace();
    if (input[i] !== '[') return null;
    i++; // skip '['

    const items = [];
    let token = '';
    let inQuotes = false;

    while (i < input.length) {
      const char = input[i];

      if (char === '"' && !inQuotes) {
        inQuotes = true;
        token += char;
        i++;
        continue;
      }

      if (char === '"' && inQuotes && input[i - 1] !== '\\') {
        inQuotes = false;
        token += char;
        i++;
        continue;
      }

      if (char === '[' && !inQuotes) {
        if (token.trim()) {
          items.push({ type: 'token', value: token.trim() });
          token = '';
        }
        const child = parseBlock();
        if (child) items.push({ type: 'block', value: child });
      } else if (char === ']' && !inQuotes) {
        if (token.trim()) {
          items.push({ type: 'token', value: token.trim() });
        }
        i++;
        break;
      } else if (char === ',' && !inQuotes) {
        if (token.trim()) {
          items.push({ type: 'token', value: token.trim() });
          token = '';
        }
        i++;
      } else {
        token += char;
        i++;
      }
    }
    return items;
  }

  function formatBlock(block, depth = 0) {
    const indent       = ' '.repeat(depth * indentStep);
    const childIndent  = ' '.repeat((depth + 1) * indentStep);

    if (block.length === 0) return `${indent}[]\n`;

    let output = `${indent}[\n`;

    block.forEach((item, idx) => {
      const next = block[idx + 1];            // look-ahead

      if (item.type === 'token') {
        output += `${childIndent}${item.value}`;
        if (next && next.type === 'token')    // token -> token
          output += ',';
        output += '\n';
      } else {                                // item is a sub-block
        output += formatBlock(item.value, depth + 1);
        if (next && next.type === 'token') {  // block -> token
          output = output.replace(/\n$/, ''); // remove the newline just written
          output += ',\n';
        }
      }
    });

    output += `${indent}]\n`;
    return output;
  }

  const tree = parseBlock();
  return formatBlock(tree).trim();
}


// Read from stdin
let input = '';
process.stdin.setEncoding('utf8');

process.stdin.on('data', chunk => input += chunk);
process.stdin.on('end', () => {
  try {
    const output = prettyPrint(input);
    console.log(output);
  } catch (err) {
    console.error('Error:', err.message);
    process.exit(1);
  }
});
