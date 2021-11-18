package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/tools"
)

type transformFilter struct {
	transformer Transformer
	mergeCount  uint
}

func (f *transformFilter) Bounds(src image.Rectangle) image.Rectangle {
	return f.transformer.Bounds(src)
}

func (f *transformFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.Y, srcb.Max.Y, 1, func(start, end int) {
		for sy := start; sy < end; sy++ {
			for sx := srcb.Min.X; sx < srcb.Max.X; sx++ {
				dx, dy, oppX, oppY := f.transformer.Transform(sx, sy)
				if oppX {
					dx = (dstb.Max.X - 1) - dx
				}
				if oppY {
					dy = (dstb.Max.Y - 1) - dy
				}

				pix := pixGetter.getPixel(sx, sy)
				pixSetter.setPixel(dx, dy, pix)
			}
		}
	})
}

func (f *transformFilter) CanMerge(filter Filter) bool {
	filt, ok := filter.(*transformFilter)
	if !ok {
		return false
	}

	if f.transformer == nil {
		return true
	}

	if _, ok := f.transformer.Merge(filt.transformer); ok {
		return true
	}

	return false
}

func (f *transformFilter) Merge(filter Filter) {
	filt := filter.(*transformFilter)

	f.mergeCount++

	if f.transformer == nil {
		f.transformer = filt.transformer
		return
	}

	out, _ := f.transformer.Merge(filt.transformer)
	f.transformer = out
}

func (f *transformFilter) Skip() bool {
	return f.transformer == nil
}

func (f *transformFilter) CanUndo(filter Filter) bool {
	filt, ok := filter.(*transformFilter)
	if !ok {
		return false
	}

	if filt.transformer == nil {
		return false
	}

	if _, ok := filt.transformer.Recreate(f.transformer); ok {
		return true
	}

	return false
}

func (f *transformFilter) Undo(filter Filter) bool {
	filt := filter.(*transformFilter)

	out, _ := filt.transformer.Recreate(f.transformer)
	f.transformer = out

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *transformFilter) Copy() Filter {
	return &transformFilter{
		transformer: f.transformer,
		mergeCount:  f.mergeCount,
	}
}

// Transform an image using given Transformer.
func Transform(transformer Transformer) MergingFilter {
	return &transformFilter{
		transformer: transformer,
		mergeCount:  1,
	}
}
