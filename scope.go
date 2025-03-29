package stache

import "bytes"

type scopeStack [][]byte

func (s *scopeStack) pushSegments(expr []byte) [][]byte {
	segments := bytes.Split(expr, []byte("."))
	*s = append(*s, segments...)
	return segments
}

func (s *scopeStack) popN(n int) {
	*s = (*s)[:len(*s)-n]
}
