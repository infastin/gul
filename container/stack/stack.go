package stack

type Element struct {
	next, prev *Element
	Value      interface{}
}

type Stack struct {
	top *Element
	len int
}

func New() *Stack {
	return &Stack{}
}

func (s *Stack) Push(value interface{}) *Element {
	el := &Element{
		next:  nil,
		prev:  nil,
		Value: value,
	}

	s.len++

	if s.top == nil {
		s.top = el
	} else {
		s.top.next = el
		el.prev = s.top
		s.top = el
	}

	return el
}

func (s *Stack) Pop() *Element {
	if s.top == nil {
		return nil
	}

	s.len--
	ret := s.top

	s.top = s.top.prev
	if s.top != nil {
		s.top.next = nil
	}

	ret.next = nil
	ret.prev = nil

	return ret
}

func (s *Stack) Clear() {
	s.len = 0
	s.top = nil
}

func (s *Stack) Empty() bool {
	return s.len == 0
}

func (s *Stack) Top() *Element {
	return s.top
}

func (s *Stack) Len() int {
	return s.len
}
