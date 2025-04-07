package jsx

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tetsuo/restache"
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

type Reader struct {
	root  *restache.Node
	phase renderPhase
	buf   *bytes.Buffer
	rd    *Renderer
}

func NewReader(root *restache.Node) *Reader {
	if root.Type != restache.ComponentNode {
		panic("restache: root must be a component node")
	}
	c := &Reader{root: root}
	c.buf = &bytes.Buffer{}
	c.phase = phaseStart
	return c
}

// Read implements io.Reader for streaming JSX rendering.
func (c *Reader) Read(p []byte) (int, error) {
	for c.buf.Len() == 0 && c.phase != phaseDone {
		switch c.phase {
		case phaseStart:
			c.phase = phaseDeps

		case phaseDeps:
			fmt.Fprintln(c.buf, "import * as React from 'react';")
			for _, attr := range c.root.Attr {
				fmt.Fprintf(c.buf, "import %s from \"./%s.jsx\";\n", attr.Key, attr.Val)
			}
			fmt.Fprintln(c.buf)
			c.phase = phaseOpen

		case phaseOpen:
			fmt.Fprint(c.buf, "export default function")
			if len(c.root.Data) != 0 {
				c.buf.WriteRune(' ')
				c.buf.Write(c.root.Data)
			}
			fmt.Fprint(c.buf, "(props) {\n")
			fmt.Fprintln(c.buf, "  return (")
			c.phase = phaseBody

		case phaseBody:
			if c.rd == nil {
				c.rd = NewRenderer(c.buf, 2, c.root)
			}
			if !c.rd.RenderNext() {
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

