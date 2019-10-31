package datastructs

type Stack struct {
	s []interface{}
}

func New() *Stack {
	return &Stack{s: []interface{}{}}
}

func (s *Stack) Push(item interface{}) {
	s.s = append(s.s, item)
}

func (s *Stack) Pop() interface{} {
	var ret interface{}
	ret, s.s = s.s[len(s.s)-1], s.s[:len(s.s)-1]
	return ret
}

func (s *Stack) Peek() interface{} {
	var ret interface{}
	ret = s.s[len(s.s)-1]
	return ret
}
