package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/tools"
)

type combineColorFilter struct {
	filters    []ColorFilter
	mergeCount uint
}

func (f *combineColorFilter) CanMerge(filter Filter) bool {
	filt, ok := filter.(*combineColorFilter)
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

		if fi, ok := f.filters[i].(MergingColorFilter); ok {
			if !fi.CanMerge(filt.filters[i]) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineColorFilter) Merge(filter Filter) {
	filt := filter.(*combineColorFilter)

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		if f.filters[i] == nil {
			f.filters[i] = filt.filters[i]
			continue
		}

		fi := f.filters[i].(MergingColorFilter)
		fi.Merge(filt.filters[i])
	}

	f.mergeCount++
}

func (f *combineColorFilter) CanUndo(filter Filter) bool {
	filt, ok := filter.(*combineColorFilter)
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

		if fi, ok := f.filters[i].(MergingColorFilter); ok {
			if !fi.CanUndo(filt.filters[i]) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineColorFilter) Undo(filter Filter) bool {
	filt := filter.(*combineColorFilter)

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		fi := f.filters[i].(MergingColorFilter)
		fi.Undo(filt.filters[i])
	}

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *combineColorFilter) Skip() bool {
	for _, filt := range f.filters {
		if filt == nil {
			continue
		}

		if filt, ok := filt.(MergingColorFilter); ok {
			if !filt.Skip() {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineColorFilter) Copy() Filter {
	r := &combineColorFilter{
		mergeCount: f.mergeCount,
	}

	r.filters = make([]ColorFilter, len(f.filters))
	for i := 0; i < len(f.filters); i++ {
		if f.filters[i] == nil {
			r.filters[i] = nil
			continue
		}

		if fi, ok := f.filters[i].(MergingColorFilter); ok {
			r.filters[i] = fi.Copy()
		} else {
			r.filters[i] = f.filters[i]
		}
	}

	return r
}

func (f *combineColorFilter) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (f *combineColorFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	for _, filt := range f.filters {
		if filt == nil {
			continue
		}

		if filt, ok := filt.(MergingColorFilter); ok {
			filt.Prepare()
		}
	}

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.Y, srcb.Max.Y, 1, func(start, end int) {
		for y := start; y < end; y++ {
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				pix := pixGetter.getPixel(x, y)

				for _, filt := range f.filters {
					if filt == nil {
						continue
					}

					pix = filt.Fn(pix)
				}

				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pix)
			}
		}
	})
}

// Creates combination of color filters and returns filter.
func CombineColorFilters(filters ...ColorFilter) MergingFilter {
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

	return &combineColorFilter{
		filters:    filters,
		mergeCount: 1,
	}
}
