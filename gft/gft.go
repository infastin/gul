package gft

import (
	"image"
	"image/draw"
	"reflect"
)

// Filter is an image filter
// Must be a pointer
type Filter interface {
	// Returns the bounds after applying filter
	Bounds(src image.Rectangle) image.Rectangle

	// Applies the filter to the src image and draws the result to the dst image
	Apply(dst draw.Image, src image.Image, parallel bool)

	// If possible, combines two filters, writes the result to an instance of interface and returns true
	// Otherwise, returns false
	Merge(filter Filter) bool

	// Returns true, if nothing will change after applying the filter
	// Otherwise, returns false
	Skip() bool

	// If possible, decombines tow filters, writes the result to an instance of interface and returns true
	// Otherwise, returns false
	Undo(filter Filter) bool

	// Returns a copy of the filter
	Copy() Filter
}

// List of filters, which allows applying multiple filters at once
// And makes use of filters' Merge, Undo and Skip methods
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
		last := l.filters[len(l.filters)-1]
		if reflect.TypeOf(last) == reflect.TypeOf(filt) {
			if last.Merge(filt) {
				if last.Skip() {
					l.filters = l.filters[:len(l.filters)-1]
				}

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
		dst = filt.Bounds(dst)
	}
	return dst
}

func (l *List) Apply(dst draw.Image, src image.Image, parallel bool) {
	if len(l.filters) == 0 {
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
		return
	}

	var tmpDst draw.Image
	var tmpSrc image.Image

	for i, filt := range l.filters {
		if i == 0 {
			tmpSrc = src
		} else {
			tmpSrc = tmpDst
		}

		if i == len(l.filters)-1 {
			tmpDst = dst
		} else {
			tmpDst = image.NewRGBA(filt.Bounds(tmpSrc.Bounds()))
		}

		filt.Apply(tmpDst, tmpSrc, parallel)
	}
}
