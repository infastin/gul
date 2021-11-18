// Go (Image) Filtering Toolkit.
// It is basically non-direct fork of GIFT (https://github.com/disintegration/gift).
// Also some filters are taken from GIMP (https://github.com/GNOME/gimp).
package gft

import (
	"image"
	"image/draw"
)

// Filter is an image filter.
// Must be a pointer.
type Filter interface {
	// Returns the bounds after applying filter.
	Bounds(src image.Rectangle) image.Rectangle

	// Applies the filter to the src image and draws the result to the dst image.
	Apply(dst draw.Image, src image.Image, parallel bool)
}

// This filter can merge other filters into itself.
type MergingFilter interface {
	Filter

	// If possible, merges one filter into an instance of interface and returns true.
	// Otherwise, returns false.
	Merge(filter Filter) bool

	// Operation opposite to Merge.
	// If possible, demerges one filter from an instance of interface.
	// If not, returns false.
	// If got nothing after decombination, returns true.
	// Otherwise, returns false.
	Undo(filter Filter) bool

	// Returns true, if nothing will change after applying the filter.
	// Otherwise, returns false.
	Skip() bool

	// Returns a copy of the filter.
	Copy() Filter
}

// List of filters, which allows applying multiple filters at once.
// And makes use of filters' Merge, Undo and Skip methods.
type List struct {
	filters []Filter
}

func MakeList(filters ...Filter) List {
	l := List{}
	for _, filt := range filters {
		l.Add(filt)
	}

	return l
}

func NewList(filters ...Filter) *List {
	l := &List{}
	for _, filt := range filters {
		l.Add(filt)
	}

	return l
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

		if last, ok := last.(MergingFilter); ok {
			if last.Merge(filt) {
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
	if last, ok := last.(MergingFilter); ok {
		if last.Undo(filt) {
			l.filters = l.filters[:len(l.filters)-1]
		}
	} else {
		l.filters = l.filters[:len(l.filters)-1]
	}
}

func (l *List) Bounds(src image.Rectangle) image.Rectangle {
	dst := src
	for _, filt := range l.filters {
		if filt, ok := filt.(MergingFilter); ok {
			if filt.Skip() {
				continue
			}
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
		if filt, ok := filt.(MergingFilter); ok {
			if filt.Skip() {
				if i == first {
					first++
				}

				continue
			}
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
