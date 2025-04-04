package component

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tetsuo/stache"
)

type renderPhase int

const (
	phaseStart renderPhase = iota
	phaseDeps
	phaseOpen
	phaseBody
	phaseClose
	phaseDone
)

// Component represents a parsed template file.
type Component struct {
	Name string
	Path string
	Root *stache.Node
	Deps []*Component

	// Render state
	phase renderPhase
	buf   *bytes.Buffer
	jsx   *jsxRenderer
}

// Read implements io.Reader for streaming JSX rendering.
func (c *Component) Read(p []byte) (int, error) {
	if c.buf == nil {
		c.buf = &bytes.Buffer{}
		c.phase = phaseStart
	}

	for c.buf.Len() == 0 && c.phase != phaseDone {
		switch c.phase {
		case phaseStart:
			c.phase = phaseDeps

		case phaseDeps:
			fmt.Fprintln(c.buf, "import * as React from 'react';")
			for _, dep := range c.Deps {
				fmt.Fprintf(c.buf, "import %s from \"./%s.jsx\";\n", dep.Name, dep.Name)
			}
			if len(c.Deps) > 0 {
				fmt.Fprintln(c.buf)
			}
			c.phase = phaseOpen

		case phaseOpen:
			fmt.Fprint(c.buf, "export default function")
			if c.Name != "" {
				c.buf.WriteRune(' ')
				c.buf.WriteString(c.Name)
			}
			fmt.Fprint(c.buf, "(props) {\n")
			fmt.Fprintln(c.buf, "  return (")
			c.phase = phaseBody

		case phaseBody:
			if c.jsx == nil {
				c.jsx = &jsxRenderer{
					w:      c.buf,
					indent: 2,
					doc:    c.Root,
				}
			}

			if !c.jsx.renderNext() {
				c.phase = phaseClose
			}

		case phaseClose:
			fmt.Fprintln(c.buf, "  );")
			fmt.Fprintln(c.buf, "}")
			c.phase = phaseDone
		}
	}

	if c.phase == phaseDone && c.buf.Len() == 0 {
		return 0, io.EOF
	}

	return c.buf.Read(p)
}

}
