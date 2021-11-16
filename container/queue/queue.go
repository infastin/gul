package queue

type Element struct {
	next, prev *Element
	Value      interface{}
}

type Queue struct {
	front, back *Element
	len         int
}

func New() *Queue {
	return &Queue{}
}

func (q *Queue) Push(value interface{}) *Element {
	el := &Element{
		next:  nil,
		prev:  nil,
		Value: value,
	}

	q.len++

	if q.back == nil {
		q.front = el
		q.back = el
	} else {
		q.back.next = el
		el.prev = q.back
		q.back = el
	}

	return el
}

func (q *Queue) Pop() *Element {
	if q.front == nil {
		return nil
	}

	q.len--
	ret := q.front

	q.front = q.front.next
	if q.front == nil {
		q.back = nil
	} else {
		q.front.prev = nil
	}

	ret.next = nil
	ret.prev = nil

	return ret
}

func (q *Queue) Clear() {
	q.len = 0
	q.front = nil
	q.back = nil
}

func (q *Queue) Empty() bool {
	return q.len == 0
}

func (q *Queue) Front() *Element {
	return q.front
}

func (q *Queue) Back() *Element {
	return q.back
}

func (q *Queue) Len() int {
	return q.len
}
