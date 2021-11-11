package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/gm32"
)

type scaleFilter struct {
	scaleX      float32
	scaleY      float32
	rfilt       ResamplingFilter
	rfiltScaleX float32
	rfiltScaleY float32
	mergeCount  uint
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
	filt := filter.(*scaleFilter)

	f.scaleX *= filt.scaleX
	f.scaleY *= filt.scaleY

	f.rfilt = filt.rfilt
	f.rfiltScaleX = filt.rfiltScaleX
	f.rfiltScaleY = filt.rfiltScaleY

	f.mergeCount++

	return true
}

func (f *scaleFilter) Undo(filter Filter) bool {
	filt := filter.(*scaleFilter)

	f.scaleX /= filt.scaleX
	f.scaleY /= filt.scaleY

	f.rfilt = filt.rfilt
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
		rfilt:       f.rfilt,
		rfiltScaleX: f.rfiltScaleX,
		rfiltScaleY: f.rfiltScaleY,
		mergeCount:  f.mergeCount,
	}
}

func Scale(scaleX, scaleY float32, rfilt ResamplingFilter, rfiltScaleX, rfiltScaleY float32) Filter {
	if scaleX == 1 && scaleY == 1 {
		return nil
	}

	return &scaleFilter{
		scaleX:      scaleX,
		scaleY:      scaleY,
		rfilt:       rfilt,
		rfiltScaleX: rfiltScaleX,
		rfiltScaleY: rfiltScaleY,
		mergeCount:  1,
	}
}

type scaleFilterAdditive struct {
	scaleFilter
}

func (f *scaleFilterAdditive) Merge(filter Filter) bool {
	filt := filter.(*scaleFilterAdditive)

	f.scaleX += filt.scaleX
	f.scaleY += filt.scaleY

	f.rfilt = filt.rfilt
	f.rfiltScaleX = filt.rfiltScaleX
	f.rfiltScaleY = filt.rfiltScaleY

	f.mergeCount++

	return true
}

func (f *scaleFilterAdditive) Undo(filter Filter) bool {
	filt := filter.(*scaleFilterAdditive)

	f.scaleX -= filt.scaleX
	f.scaleY -= filt.scaleY

	f.rfilt = filt.rfilt
	f.rfiltScaleX = filt.rfiltScaleX
	f.rfiltScaleY = filt.rfiltScaleY

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *scaleFilterAdditive) Copy() Filter {
	return &scaleFilterAdditive{
		scaleFilter: f.scaleFilter,
	}
}

func ScaleAdditive(scaleX, scaleY float32, rfilt ResamplingFilter, rfiltScaleX, rfiltScaleY float32) Filter {
	if scaleX == 1 && scaleY == 1 {
		return nil
	}

	return &scaleFilterAdditive{
		scaleFilter: scaleFilter{
			scaleX:      scaleX,
			scaleY:      scaleY,
			rfilt:       rfilt,
			rfiltScaleX: rfiltScaleX,
			rfiltScaleY: rfiltScaleY,
			mergeCount:  1,
		},
	}
}
