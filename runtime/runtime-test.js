const { createContext, Template } = require("./runtime");

function build(tree, state, options) {
  const render = createContext(tree);
  return render.call(new Template(), state, options);
}

const testData = [
  {
    // object is accessible within the scope
    html: "<x>{yar}{#bar}{nor}<k>{ssw}</k>{/bar}</x>",
    tree: [
      "x",
      {},
      [
        [3, "yar"],
        [
          2,
          "bar",
          [
            [3, "nor"],
            ["k", {}, [[3, "ssw"]]],
          ],
        ],
      ],
    ],
    state: {
      yar: "555",
      bar: [{ nor: "xyz", ssw: "rre" }],
    },
    expected: {
      tag: "x",
      props: {},
      children: ["555", "xyz", { tag: "k", props: {}, children: ["rre"] }],
    },
  },
  {
    // nested array sections
    html: "<x>{#bar}<u>{#da}<k>{m}{#ow}8<z></z>{/ow}</k>{/da}</u>{#dor}{kw}{#mor}{sd}55{/mor}df{ke}1s{/dor}63{jh}{/bar}{kq}</x>",
    tree: [
      "x",
      {},
      [
        [
          2,
          "bar",
          [
            [
              "u",
              {},
              [
                [
                  2,
                  "da",
                  [
                    [
                      "k",
                      {},
                      [
                        [3, "m"],
                        [2, "ow", ["8", ["z", {}, []]]],
                      ],
                    ],
                  ],
                ],
              ],
            ],
            [
              2,
              "dor",
              [[3, "kw"], [2, "mor", [[3, "sd"], "55"]], "df", [3, "ke"], "1s"],
            ],
            "63",
            [3, "jh"],
          ],
        ],
        [3, "kq"],
      ],
    ],
    state: {
      kq: "qkw",
      bar: {
        da: {
          m: "ram",
          ow: [1, 1, 1],
        },
        dor: {
          kw: "ghj",
          mor: [{ sd: "77" }, { sd: "55" }],
          ke: "s1",
        },
        jh: "6734",
      },
    },
    expected: {
      tag: "x",
      props: {},
      children: [
        {
          tag: "u",
          props: {},
          children: [
            {
              tag: "k",
              props: {},
              children: [
                "ram",
                "8",
                { tag: "z", props: {}, children: [] },
                "8",
                { tag: "z", props: {}, children: [] },
                "8",
                { tag: "z", props: {}, children: [] },
              ],
            },
          ],
        },
        "ghj",
        "77",
        "55",
        "55",
        "55",
        "df",
        "s1",
        "1s",
        "63",
        "6734",
        "qkw",
      ],
    },
  },
  {
    html: "<x>{#s}{s}{#z}{c}{/z}{#u}{k}{#r}{u}{#y}55{/y}{/r}{/u}{/s}</x>",
    tree: [
      "x",
      {},
      [
        [
          2,
          "s",
          [
            [3, "s"],
            [2, "z", [[3, "c"]]],
            [
              2,
              "u",
              [
                [3, "k"],
                [
                  2,
                  "r",
                  [
                    [3, "u"],
                    [2, "y", ["55"]],
                  ],
                ],
              ],
            ],
          ],
        ],
      ],
    ],
    state: {
      s: {
        s: "sd",
        z: false,
        u: [
          {
            k: "12",
            r: [{ u: "34" }, { u: "23", y: true }],
          },
          {
            k: "65",
            r: [{ u: "89" }],
          },
        ],
      },
    },
    expected: {
      tag: "x",
      props: {},
      children: ["sd", "12", "34", "23", "55", "65", "89"],
    },
  },
  {
    html: "<x>{^s}foo{/s}</x>",
    tree: ["x", {}, [[4, "s", ["foo"]]]],
    state: {
      s: true,
    },
    expected: {
      tag: "x",
      props: {},
      children: [],
    },
  },
  {
    html: "<x>{^s}foo{/s}</x>",
    tree: ["x", {}, [[4, "s", ["foo"]]]],
    state: {
      s: false,
    },
    expected: {
      tag: "x",
      props: {},
      children: ["foo"],
    },
  },
  {
    html: "<x>{^s}foo{/s}</x>",
    tree: ["x", {}, [[4, "s", ["foo"]]]],
    state: {
      s: [],
    },
    expected: {
      tag: "x",
      props: {},
      children: ["foo"],
    },
  },
  {
    html: "<x>{^s}{x}{/s}</x>",
    tree: ["x", {}, [[4, "s", [[3, "x"]]]]],
    state: {
      s: { x: 1 },
    },
    expected: {
      tag: "x",
      props: {},
      children: [],
    },
  },
];

testData.forEach((d) => {
  const actual = build(d.tree, d.state, {
    createElement: (tag, props, children) => {
      delete props.parseTree;
      delete props.traverseChildren;
      delete props.state;
      return { tag, props, children };
    },
  });
  if (!deepEqual(actual, d.expected)) {
    throw new Error(
      "not equal: " +
        JSON.stringify(actual) +
        " != " +
        JSON.stringify(d.expected)
    );
  }
});

function deepEqual(obj1, obj2) {
  if (obj1 === obj2) {
    return true;
  } else if (isObject(obj1) && isObject(obj2)) {
    if (Object.keys(obj1).length !== Object.keys(obj2).length) {
      return false;
    }
    for (var prop in obj1) {
      if (!deepEqual(obj1[prop], obj2[prop])) {
        return false;
      }
    }
    return true;
  }
  function isObject(obj) {
    if (typeof obj === "object" && obj != null) {
      return true;
    } else {
      return false;
    }
  }
}
