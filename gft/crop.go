package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/gm32"
	"github.com/srwiley/rasterx"
)

type cropRectangleFilter struct {
	startX, startY float32
	width, height  float32
}

func (f *cropRectangleFilter) Bounds(src image.Rectangle) image.Rectangle {
	srcb := src.Bounds()
	srcWidth := float32(srcb.Dx())
	srcHeight := float32(srcb.Dy())

	dstWidth := int(gm32.Round(srcWidth * f.width))
	dstHeight := int(gm32.Round(srcHeight * f.height))

	return image.Rect(0, 0, dstWidth, dstHeight)
}

func (f *cropRectangleFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	srcWidth := srcb.Dx()
	srcHeight := srcb.Dy()

	dstb := dst.Bounds()
	startX := int(gm32.Round(float32(srcWidth) * f.startX))
	startY := int(gm32.Round(float32(srcHeight) * f.startY))

	draw.Draw(dst, dstb, src, image.Pt(startX, startY), draw.Over)
}

func (f *cropRectangleFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*cropRectangleFilter)
	if !ok {
		return false
	}

	f.startX += filt.startX * f.width
	f.startY += filt.startY * f.height

	f.width *= filt.width
	f.height *= filt.height

	return true
}

func (f *cropRectangleFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*cropRectangleFilter)
	if !ok {
		return false
	}

	f.height /= filt.height
	f.width /= filt.width

	f.startX -= filt.startX * f.width
	f.startY -= filt.startY * f.height

	return false
}

func (f *cropRectangleFilter) Skip() bool {
	return f.startX == 0 && f.startY == 0 && f.height == 1 && f.width == 1
}

func (f *cropRectangleFilter) Copy() Filter {
	return &cropRectangleFilter{
		startX: f.startX,
		startY: f.startY,
		width:  f.width,
		height: f.height,
	}
}

// Crops an image starting at a given position (startX, startY) with a rectangle of a given size (width, height).
// The position and size parameters should be in the range [0, 1].
// Example: You have an image and you want to crop the bottom-right quarter of it.
// Then pos will be (0.5, 0.5) and size will be (0.5, 0.5).
func CropRectangle(startX, startY, width, height float32) Filter {
	if startX == 0 && startY == 0 && height == 1 && width == 1 {
		return nil
	}

	startX = gm32.Clamp(startX, 0, 1)
	startY = gm32.Clamp(startY, 0, 1)

	width = gm32.Clamp(width, 0, 1)
	height = gm32.Clamp(height, 0, 1)

	if startX+width > 1 {
		width = 1 - startX
	}

	if startY+height > 1 {
		height = 1 - startY
	}

	return &cropRectangleFilter{
		startX: startX,
		startY: startY,
		width:  width,
		height: height,
	}
}

type cropEllipseFilter struct {
	cx, cy float32
	rx, ry float32
}

func (f *cropEllipseFilter) Bounds(src image.Rectangle) image.Rectangle {
	srcb := src.Bounds()
	srcWidth := float32(srcb.Dx())
	srcHeight := float32(srcb.Dy())

	leftX := gm32.Min(f.rx, f.cx)
	rightX := gm32.Min(f.rx, 1-f.cx)

	topY := gm32.Min(f.ry, f.cy)
	botY := gm32.Min(f.ry, 1-f.cy)

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

	leftX := gm32.Min(f.rx, f.cx)
	topY := gm32.Min(f.ry, f.cy)

	startX := f.cx - leftX
	startY := f.cy - topY

	offset := image.Point{
		X: int(gm32.Round(startX * srcWidth)),
		Y: int(gm32.Round(startY * srcHeight)),
	}

	dstb := dst.Bounds()
	dstWidth := dstb.Dx()
	dstHeight := dstb.Dy()

	cx := float64(srcWidth * (f.cx - startX))
	cy := float64(srcHeight * (f.cy - startY))
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
		cx: f.cx,
		cy: f.cy,
		rx: f.rx,
		ry: f.ry,
	}
}

// Crops an image with an ellipse of a radii (rx, ry) with the center at a given position (cx, cy).
// The position and radii parameters should be in the range [0, 1].
func CropEllipse(cx, cy, rx, ry float32) Filter {
	cx = gm32.Clamp(cx, 0, 1)
	cy = gm32.Clamp(cy, 0, 1)

	rx = gm32.Clamp(rx, 0, 1)
	ry = gm32.Clamp(ry, 0, 1)

	return &cropEllipseFilter{
		cx: cx,
		cy: cy,
		rx: rx,
		ry: ry,
	}
}
