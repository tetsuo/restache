package stache

type controlFrame struct {
	typ      NodeType
	name     []byte
	segments [][]byte
}

type controlStack []controlFrame

func (s *controlStack) push(cf controlFrame) {
	*s = append(*s, controlFrame{
		typ:  cf.typ,
		name: cf.name,
	})
}

func (s *controlStack) top() *controlFrame {
	if len(*s) == 0 {
		return nil
	}
	return &(*s)[len(*s)-1]
}

func (s *controlStack) pop() {
	*s = (*s)[:len(*s)-1]
}
