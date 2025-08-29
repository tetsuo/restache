import { initPlaygrounds } from './playground.mjs'

export default function(runtime) {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch(runtime), go.importObject).then(result => {
    go.run(result.instance);
    initPlaygrounds(input => {
      const result = runExampleCode(input);
      if (runtime.startsWith("restache") && result.output?.startsWith("export default function ($0) {return ")) {
        result.output = result.output.slice(37, -2) // Show the relevant parts only
      }
      return result
    });
  });
}
