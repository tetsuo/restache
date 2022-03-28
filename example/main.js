const React = require("elm-ts/lib/React");
const cmd = require("elm-ts/lib/Cmd");
const { render } = require("react-dom");
const { Layout, Template } = require("@onur1/stache/layout");
const createComponent = require("@onur1/stache/react");
const update = (_, model) => [model, cmd.none];
const view = (layout) => (model) => (_) => layout(model);
const init = [
  {
    items: [
      { name: "Apple" },
      { name: "Orange" },
      { name: "Watermelon" },
    ],
  },
  cmd.none,
];
fetch("http://localhost:7882/layout")
  .then((response) => response.json())
  .then((templates) => {
    return React.run(
      React.program(
        init,
        update,
        view(
          createComponent(
            Layout(templates.map(([name, trees]) => Template(name, trees)))
          )
        )
      ),
      (dom) => render(dom, document.getElementById("content"))
    );
  });
