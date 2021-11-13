package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/gm32"
)

type cropFilter struct {
	pos  Position
	size Size
}

func (f *cropFilter) Bounds(src image.Rectangle) image.Rectangle {
	srcb := src.Bounds()
	srcWidth := float32(srcb.Dx())
	srcHeight := float32(srcb.Dy())

	dstWidth := int(gm32.Round(srcWidth * f.size.Width))
	dstHeight := int(gm32.Round(srcHeight * f.size.Height))

	return image.Rect(0, 0, dstWidth, dstHeight)
}

func (f *cropFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	srcWidth := srcb.Dx()
	srcHeight := srcb.Dy()

	dstb := dst.Bounds()
	startX := int(gm32.Round(float32(srcWidth) * f.pos.X))
	startY := int(gm32.Round(float32(srcHeight) * f.pos.Y))

	draw.Draw(dst, dstb, src, image.Pt(startX, startY), draw.Over)
}

func (f *cropFilter) Merge(filter Filter) bool {
	filt := filter.(*cropFilter)

	f.pos.X += filt.pos.X * f.size.Width
	f.pos.Y += filt.pos.Y * f.size.Height

	f.size.Width *= filt.size.Width
	f.size.Height *= filt.size.Height

	return true
}

func (f *cropFilter) Undo(filter Filter) bool {
	filt := filter.(*cropFilter)

	f.size.Height /= filt.size.Height
	f.size.Width /= filt.size.Width

	f.pos.X -= filt.pos.X * f.size.Width
	f.pos.Y -= filt.pos.Y * f.size.Height

	return false
}

func (f *cropFilter) Skip() bool {
	return f.pos.X == 0 && f.pos.Y == 0 && f.size.Height == 1 && f.size.Width == 1
}

func (f *cropFilter) Copy() Filter {
	return &cropFilter{
		pos:  f.pos,
		size: f.size,
	}
}

func Crop(pos Position, size Size) Filter {
	if pos.X == 0 && pos.Y == 0 && size.Height == 1 && size.Width == 1 {
		return nil
	}

	return &cropFilter{
		pos:  pos,
		size: size,
	}
}
