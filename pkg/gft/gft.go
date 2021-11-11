package gft

import (
	"image"
	"image/draw"
	"reflect"
)

type Filter interface {
	Bounds(src image.Rectangle) image.Rectangle
	Apply(dst draw.Image, src image.Image, parallel bool)
	Merge(filter Filter) bool
	Skip() bool
	Undo(filter Filter) bool
	Copy() Filter
}

type List struct {
	filters []Filter
}

func MakeList(filters ...Filter) List {
	return List{
		filters: filters,
	}
}

func NewList(filters ...Filter) *List {
	return &List{
		filters: filters,
	}
}

func (l *List) Clear() {
	l.filters = nil
}

func (l *List) Empty() bool {
	return len(l.filters) == 0
}

func (l *List) Add(filt Filter) {
	if len(l.filters) != 0 {
		last := len(l.filters) - 1
		if reflect.TypeOf(l.filters[last]) == reflect.TypeOf(filt) {
			if l.filters[last].Merge(filt) {
				return
			}
		}
	}

	l.filters = append(l.filters, filt)
}

func (l *List) Undo(filt Filter) {
	if len(l.filters) == 0 {
		return
	}

	last := l.filters[len(l.filters)-1]
	if reflect.TypeOf(last) != reflect.TypeOf(filt) {
		return
	}

	if last.Undo(filt) {
		l.filters = l.filters[:len(l.filters)-1]
	}
}

func (l *List) Bounds(src image.Rectangle) image.Rectangle {
	dst := src
	for _, filt := range l.filters {
		if filt.Skip() {
			continue
		}

		dst = filt.Bounds(dst)
	}
	return dst
}

func (l *List) Apply(dst draw.Image, src image.Image, parallel bool) {
	if len(l.filters) == 0 {
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
		return
	}

	first, last := 0, len(l.filters)-1
	var tmpDst draw.Image
	var tmpSrc image.Image

	for i, filt := range l.filters {
		if filt.Skip() {
			if i == first {
				first++
			}

			continue
		}

		if i == first {
			tmpSrc = src
		} else {
			tmpSrc = tmpDst
		}

		if i == last {
			tmpDst = dst
		} else {
			tmpDst = image.NewRGBA(filt.Bounds(tmpSrc.Bounds()))
		}

		filt.Apply(tmpDst, tmpSrc, parallel)
	}

	if tmpDst != dst {
		if tmpSrc == nil {
			tmpSrc = src
		} else {
			tmpSrc = tmpDst
		}

		draw.Draw(dst, dst.Bounds(), tmpSrc, tmpSrc.Bounds().Min, draw.Over)
	}
}
