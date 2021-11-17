package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/gm32"
	"github.com/srwiley/rasterx"
)

type cropRectangleFilter struct {
	pos  Position
	size Size
}

func (f *cropRectangleFilter) Bounds(src image.Rectangle) image.Rectangle {
	srcb := src.Bounds()
	srcWidth := float32(srcb.Dx())
	srcHeight := float32(srcb.Dy())

	dstWidth := int(gm32.Round(srcWidth * f.size.Width))
	dstHeight := int(gm32.Round(srcHeight * f.size.Height))

	return image.Rect(0, 0, dstWidth, dstHeight)
}

func (f *cropRectangleFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	srcWidth := srcb.Dx()
	srcHeight := srcb.Dy()

	dstb := dst.Bounds()
	startX := int(gm32.Round(float32(srcWidth) * f.pos.X))
	startY := int(gm32.Round(float32(srcHeight) * f.pos.Y))

	draw.Draw(dst, dstb, src, image.Pt(startX, startY), draw.Over)
}

func (f *cropRectangleFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*cropRectangleFilter)
	if !ok {
		return false
	}

	f.pos.X += filt.pos.X * f.size.Width
	f.pos.Y += filt.pos.Y * f.size.Height

	f.size.Width *= filt.size.Width
	f.size.Height *= filt.size.Height

	return true
}

func (f *cropRectangleFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*cropRectangleFilter)
	if !ok {
		return false
	}

	f.size.Height /= filt.size.Height
	f.size.Width /= filt.size.Width

	f.pos.X -= filt.pos.X * f.size.Width
	f.pos.Y -= filt.pos.Y * f.size.Height

	return false
}

func (f *cropRectangleFilter) Skip() bool {
	return f.pos.X == 0 && f.pos.Y == 0 && f.size.Height == 1 && f.size.Width == 1
}

func (f *cropRectangleFilter) Copy() Filter {
	return &cropRectangleFilter{
		pos:  f.pos,
		size: f.size,
	}
}

// Crops an image starting at a given position with a rectangle of a given size.
// The position and size parameters should be in the range [0, 1].
// Example: You have an image and you want to crop the bottom-right quarter of it.
// Then pos will be (0.5, 0.5) and size will be (0.5, 0.5).
func CropRectangle(pos Position, size Size) Filter {
	if pos.X == 0 && pos.Y == 0 && size.Height == 1 && size.Width == 1 {
		return nil
	}

	return &cropRectangleFilter{
		pos:  pos,
		size: size,
	}
}

type cropEllipseFilter struct {
	center Position
	rx, ry float32
}

func (f *cropEllipseFilter) Bounds(src image.Rectangle) image.Rectangle {
	srcb := src.Bounds()
	srcWidth := float32(srcb.Dx())
	srcHeight := float32(srcb.Dy())

	leftX := gm32.Min(f.rx, f.center.X)
	rightX := gm32.Min(f.rx, 1-f.center.X)

	topY := gm32.Min(f.ry, f.center.Y)
	botY := gm32.Min(f.ry, 1-f.center.Y)

	dstWidth := int(gm32.Round(srcWidth * (leftX + rightX)))
	dstHeight := int(gm32.Round(srcHeight * (topY + botY)))

	return image.Rect(0, 0, dstWidth, dstHeight)
}

func (f *cropEllipseFilter) Merge(Filter) bool {
	return false
}

func (f *cropEllipseFilter) Undo(Filter) bool {
	return false
}

func (f *cropEllipseFilter) Skip() bool {
	return false
}

func (f *cropEllipseFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	srcWidth := float32(srcb.Dx())
	srcHeight := float32(srcb.Dy())

	leftX := gm32.Min(f.rx, f.center.X)
	topY := gm32.Min(f.ry, f.center.Y)

	startX := f.center.X - leftX
	startY := f.center.Y - topY

	offset := image.Point{
		X: int(gm32.Round(startX * srcWidth)),
		Y: int(gm32.Round(startY * srcHeight)),
	}

	dstb := dst.Bounds()
	dstWidth := dstb.Dx()
	dstHeight := dstb.Dy()

	cx := float64(srcWidth * (f.center.X - startX))
	cy := float64(srcHeight * (f.center.Y - startY))
	rx := float64(srcWidth * f.rx)
	ry := float64(srcHeight * f.ry)

	scanner := rasterx.NewScannerGV(dstWidth, dstHeight, dst, dstb)
	scanner.Source = src
	scanner.Offset = offset

	filler := rasterx.NewFiller(dstWidth, dstHeight, scanner)
	rasterx.AddEllipse(cx, cy, rx, ry, 0, filler)
	filler.Draw()
}

func (f *cropEllipseFilter) Copy() Filter {
	return &cropEllipseFilter{
		center: f.center,
		rx:     f.rx,
		ry:     f.ry,
	}
}

// Crops an image starting at a given position with an ellipse of the a radii.
// The position and radii parameters should be in the range [0, 1].
func CropEllipse(center Position, rx, ry float32) Filter {
	center.X = gm32.Clamp(center.X, 0, 1)
	center.Y = gm32.Clamp(center.Y, 0, 1)

	rx = gm32.Clamp(rx, 0, 1)
	ry = gm32.Clamp(ry, 0, 1)

	return &cropEllipseFilter{
		center: center,
		rx:     rx,
		ry:     ry,
	}
}
