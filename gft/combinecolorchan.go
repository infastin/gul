package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/gm32"
	"github.com/infastin/gul/tools"
)

type combineColorchanFilter struct {
	filters    []ColorchanFilter
	luts       [][]float32
	mergeCount uint
}

func (f *combineColorchanFilter) CanMerge(filter Filter) bool {
	filt, ok := filter.(*combineColorchanFilter)
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

		if fi, ok := f.filters[i].(MergingColorchanFilter); ok {
			if !fi.CanMerge(filt.filters[i]) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineColorchanFilter) Merge(filter Filter) {
	filt := filter.(*combineColorchanFilter)

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		if f.filters[i] == nil {
			f.filters[i] = filt.filters[i]
			continue
		}

		fi := f.filters[i].(MergingColorchanFilter)
		fi.Merge(filt.filters[i])
	}

	f.mergeCount++
}

func (f *combineColorchanFilter) CanUndo(filter Filter) bool {
	filt, ok := filter.(*combineColorchanFilter)
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

		if fi, ok := f.filters[i].(MergingColorchanFilter); ok {
			if !fi.CanUndo(filt.filters[i]) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineColorchanFilter) Undo(filter Filter) bool {
	filt := filter.(*combineColorchanFilter)

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		fi := f.filters[i].(MergingColorchanFilter)
		fi.Undo(filt.filters[i])
	}

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *combineColorchanFilter) Skip() bool {
	for _, filt := range f.filters {
		if filt == nil {
			continue
		}

		if filt, ok := filt.(MergingColorchanFilter); ok {
			if !filt.Skip() {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (f *combineColorchanFilter) Copy() Filter {
	r := &combineColorchanFilter{
		mergeCount: f.mergeCount,
	}

	r.filters = make([]ColorchanFilter, len(f.filters))
	for i := 0; i < len(f.filters); i++ {
		if f.filters[i] == nil {
			r.filters[i] = nil
			continue
		}

		if fi, ok := f.filters[i].(MergingColorchanFilter); ok {
			r.filters[i] = fi.Copy()
		} else {
			r.filters[i] = f.filters[i]
		}
	}

	return r
}

func (f *combineColorchanFilter) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (f *combineColorchanFilter) makeLut(lutSize int, index int) {
	lutLen := len(f.luts[index])
	start := 0

	if lutLen == 0 {
		f.luts[index] = make([]float32, lutSize)
	} else if lutLen < lutSize {
		newLut := make([]float32, lutSize)
		copy(newLut[:lutLen], f.luts[index])
		f.luts[index] = newLut
		start = lutLen
	}

	q := float32(1) / float32(lutSize-1)
	for i := start; i < lutSize; i++ {
		v := float32(i) * q
		f.luts[index][i] = f.filters[index].Fn(v)
	}
}

func (f *combineColorchanFilter) getFromLut(x float32, index int) float32 {
	i := int(gm32.Round(x * float32(len(f.luts[index])-1)))
	return f.luts[index][i]
}

func (f *combineColorchanFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	useLut := make([]bool, len(f.filters))

	for i, filt := range f.filters {
		if filt == nil {
			continue
		}

		if filt, ok := filt.(MergingColorchanFilter); ok {
			filt.Prepare()
		}

		if filt.UseLut() {
			lutSize := len(f.luts[i])
			neededLutSize := 0

			switch pixGetter.img.(type) {
			case *image.RGBA, *image.NRGBA, *image.YCbCr, *image.Gray, *image.CMYK:
				neededLutSize = 0xff + 1
			default:
				neededLutSize = 0xffff + 1
			}

			numCalc := srcb.Dx() * srcb.Dy() * 3
			if numCalc > neededLutSize*2 {
				if lutSize != neededLutSize {
					f.makeLut(neededLutSize, i)
				}

				useLut[i] = true
			}
		} else {
			useLut[i] = false
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

				for i, filt := range f.filters {
					if filt == nil {
						continue
					}

					if useLut[i] {
						pix.r = f.getFromLut(pix.r, i)
						pix.g = f.getFromLut(pix.g, i)
						pix.b = f.getFromLut(pix.b, i)
					} else {
						pix.r = filt.Fn(pix.r)
						pix.g = filt.Fn(pix.g)
						pix.b = filt.Fn(pix.b)
					}
				}

				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pix)
			}
		}
	})
}

// Creates combination of colorchan filters and returns filter.
func CombineColorchanFilters(filters ...ColorchanFilter) MergingFilter {
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

	return &combineColorchanFilter{
		filters:    filters,
		luts:       make([][]float32, len(filters)),
		mergeCount: 1,
	}
}
