package gft

import (
	"image"
	"image/draw"
)

type combineFilter struct {
	filters    []Filter
	mergeCount uint
}

func (f *combineFilter) CanMerge(filter Filter) bool {
	filt, ok := filter.(*combineFilter)
	if !ok {
		return false
	}

	if len(f.filters) != len(filt.filters) {
		return false
	}

	for i := 0; i < len(f.filters); i++ {
		if f.filters[i] == nil || filt.filters[i] == nil {
			continue
		}

		if fi, ok := f.filters[i].(MergingFilter); ok {
			if !fi.CanMerge(filt.filters[i]) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineFilter) Merge(filter Filter) {
	filt := filter.(*combineFilter)

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		if f.filters[i] == nil {
			f.filters[i] = filt.filters[i]
			continue
		}

		fi := f.filters[i].(MergingFilter)
		fi.Merge(filt.filters[i])
	}

	f.mergeCount++
}

func (f *combineFilter) CanUndo(filter Filter) bool {
	filt, ok := filter.(*combineFilter)
	if !ok {
		return false
	}

	if len(f.filters) != len(filt.filters) {
		return false
	}

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		if f.filters[i] == nil {
			return false
		}

		if fi, ok := f.filters[i].(MergingFilter); ok {
			if !fi.CanUndo(filt.filters[i]) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineFilter) Undo(filter Filter) bool {
	filt := filter.(*combineFilter)

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		fi := f.filters[i].(MergingFilter)
		fi.Undo(filt.filters[i])
	}

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *combineFilter) Skip() bool {
	for _, filt := range f.filters {
		if filt == nil {
			continue
		}

		if filt, ok := filt.(MergingFilter); ok {
			if !filt.Skip() {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineFilter) Copy() Filter {
	r := &combineFilter{
		mergeCount: f.mergeCount,
	}

	r.filters = make([]Filter, len(f.filters))
	for i := 0; i < len(f.filters); i++ {
		if f.filters[i] == nil {
			r.filters[i] = nil
			continue
		}

		if fi, ok := f.filters[i].(MergingFilter); ok {
			r.filters[i] = fi.Copy()
		} else {
			r.filters[i] = f.filters[i]
		}
	}

	return r
}

func (f *combineFilter) Bounds(src image.Rectangle) image.Rectangle {
	dst := src
	for _, filt := range f.filters {
		if filt, ok := filt.(MergingFilter); ok {
			if filt.Skip() {
				continue
			}
		}

		dst = filt.Bounds(dst)
	}
	return dst
}

func (f *combineFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	first, last := 0, len(f.filters)-1
	var tmpDst draw.Image
	var tmpSrc image.Image

	for i, filt := range f.filters {
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

// Creates combination of filters and returns filter.
func CombineFilters(filters ...Filter) MergingFilter {
	if len(filters) == 0 {
		return nil
	}

	numNil := 0
	for _, filt := range filters {
		if filt == nil {
			numNil++
		}
	}

	if numNil == len(filters) {
		return nil
	}

	return &combineFilter{
		filters:    filters,
		mergeCount: 1,
	}
}
