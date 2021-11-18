package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/gm32"
)

type scaleFilter struct {
	scaleX   float32
	scaleY   float32
	additive bool

	rfilt       ResamplingFilter
	oldrfilt    ResamplingFilter
	rfiltScaleX float32
	rfiltScaleY float32

	mergeCount uint
}

func (f *scaleFilter) Bounds(src image.Rectangle) image.Rectangle {
	srcb := src.Bounds()
	srcWidth := float32(srcb.Dx())
	srcHeight := float32(srcb.Dy())

	dstWidth := int(gm32.Round(srcWidth * f.scaleX))
	dstHeight := int(gm32.Round(srcHeight * f.scaleY))

	return image.Rect(0, 0, dstWidth, dstHeight)
}

func (f *scaleFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	resamp := newResampler(f.rfilt, f.rfiltScaleX, f.rfiltScaleY)
	resamp.resample(dst, src, true)
}

func (f *scaleFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*scaleFilter)
	if !ok {
		return false
	}

	if filt.additive {
		f.scaleX += filt.scaleX
		f.scaleY += filt.scaleY
	} else {
		f.scaleX *= filt.scaleX
		f.scaleY *= filt.scaleY
	}

	filt.oldrfilt = f.rfilt
	f.rfilt = filt.rfilt

	f.rfiltScaleX = filt.rfiltScaleX
	f.rfiltScaleY = filt.rfiltScaleY

	f.mergeCount++

	return true
}

func (f *scaleFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*scaleFilter)
	if !ok {
		return false
	}

	if filt.additive {
		f.scaleX = gm32.Max(1.0e-5, f.scaleX-filt.scaleX)
		f.scaleY = gm32.Max(1.0e-5, f.scaleY-filt.scaleY)
	} else {
		f.scaleX /= filt.scaleX
		f.scaleY /= filt.scaleY
	}

	f.rfilt = filt.oldrfilt

	f.rfiltScaleX = filt.rfiltScaleX
	f.rfiltScaleY = filt.rfiltScaleY

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *scaleFilter) Skip() bool {
	return f.scaleX == 1 && f.scaleY == 1
}

func (f *scaleFilter) Copy() Filter {
	return &scaleFilter{
		scaleX:      f.scaleX,
		scaleY:      f.scaleY,
		additive:    f.additive,
		rfilt:       f.rfilt,
		oldrfilt:    f.oldrfilt,
		rfiltScaleX: f.rfiltScaleX,
		rfiltScaleY: f.rfiltScaleY,
		mergeCount:  f.mergeCount,
	}
}

// Scales an image by scaleX horizontally and by scaleY vertically using given ResamplingFilter.
// The scaleX, scaleY, rfiltScaleX, rfiltScaleY parameters must be greater than 0.
//
// If additive true, then scaleX and scaleY parameters will be summed  when merging this filter into another.
// Otherwise, parameters will be multiplied.
//
// The rfiltScaleX and rfiltScaleY values less than 1.0 cause aliasing, but create sharper looking mips.
// The values greater than 1.0 cause anti-aliasing, but create more blurred looking mips.
func Scale(scaleX, scaleY float32, additive bool, rfilt ResamplingFilter, rfiltScaleX, rfiltScaleY float32) MergingFilter {
	if scaleX == 1 && scaleY == 1 {
		return nil
	}

	scaleX = gm32.Max(1.0e-5, scaleX)
	scaleY = gm32.Max(1.0e-5, scaleY)
	rfiltScaleX = gm32.Max(1.0e-5, rfiltScaleX)
	rfiltScaleY = gm32.Max(1.0e-5, rfiltScaleY)

	return &scaleFilter{
		scaleX:      scaleX,
		scaleY:      scaleY,
		additive:    additive,
		rfilt:       rfilt,
		oldrfilt:    rfilt,
		rfiltScaleX: rfiltScaleX,
		rfiltScaleY: rfiltScaleY,
		mergeCount:  1,
	}
}
