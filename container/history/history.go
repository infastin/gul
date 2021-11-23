package history

type History struct {
	index  int
	values []interface{}
	Defv   interface{}
}

func Make(defv interface{}) History {
	return History{
		index: -1,
		Defv:  defv,
	}
}

func New(defv interface{}) *History {
	return &History{
		index: -1,
		Defv:  defv,
	}
}

func (h *History) Set(val interface{}) {
	vlen := len(h.values)

	if vlen > 0 && h.index < vlen-1 {
		h.values = h.values[:h.index+1]
	}

	h.values = append(h.values, val)
	h.index++
}

func (h *History) Back() interface{} {
	if h.index <= 0 {
		if h.index == 0 {
			h.index--
		}
		return h.Defv
	}

	h.index--
	return h.values[h.index]
}

func (h *History) Forward() interface{} {
	if h.index >= len(h.values)-1 {
		return h.Defv
	}

	h.index++
	return h.values[h.index]
}

func (h *History) Get() interface{} {
	if h.index == -1 {
		return h.Defv
	}

	return h.values[h.index]
}

func (h *History) Len() int {
	return len(h.values)
}

func (h *History) Index() int {
	return h.index
}

func (h *History) Reset() {
	h.index = -1
}

func (h *History) Clear() {
	h.index = -1
	h.values = h.values[:0]
}

func (h *History) Empty() bool {
	return len(h.values) == 0
}

func (h *History) AtBegin() bool {
	return h.index == -1
}

func (h *History) AtEnd() bool {
	return h.index == len(h.values)-1
}
