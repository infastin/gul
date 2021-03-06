package gft

import (
	"image"
	"image/draw"
	"math"

	"github.com/infastin/gul/gm32"
	"github.com/infastin/gul/tools"
)

type Interpolation int

const (
	NearestNeighborInterpolation Interpolation = iota
	BilinearInterpolation
	BicubicInterpolation
)

type rotateFilter struct {
	rad              float32
	interpolation    Interpolation
	oldInterpolation Interpolation
	mergeCount       uint
}

func (f *rotateFilter) Bounds(src image.Rectangle) image.Rectangle {
	srcb := src.Bounds()
	srcWidth := srcb.Dx()
	srcHeight := srcb.Dy()

	rad := gm32.Mod(f.rad, math.Pi)

	switch {
	case rad <= -math.Pi/2:
		srcWidth, srcHeight = srcHeight, srcWidth
		rad += math.Pi / 2
	case rad >= math.Pi/2:
		srcWidth, srcHeight = srcHeight, srcWidth
		rad -= math.Pi / 2
	}

	if rad == 0 {
		return image.Rect(0, 0, srcWidth, srcHeight)
	}

	rad = gm32.Clamp(-rad, -math.Pi/2, math.Pi/2)
	sine, cosine := gm32.Sincos(rad)

	fdstWidth := float32(srcWidth)*cosine + float32(srcHeight)*gm32.Abs(sine)
	fdstHeight := float32(srcWidth)*gm32.Abs(sine) + float32(srcHeight)*cosine

	dstWidth := int(gm32.Round(fdstWidth))
	dstHeight := int(gm32.Round(fdstHeight))

	return image.Rect(0, 0, dstWidth, dstHeight)
}

func (f *rotateFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	f.rad = gm32.Mod(f.rad, 2*math.Pi)
	if f.rad == 0 {
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
		return
	}

	srcb := src.Bounds()
	srcWidth := srcb.Dx()
	srcHeight := srcb.Dy()

	dstb := dst.Bounds()
	dstWidth := dstb.Dx()
	dstHeight := dstb.Dy()

	rad := -f.rad
	sine, cosine := gm32.Sincos(rad)

	halfSrcWidth := float32(srcWidth) / 2
	halfSrcHeight := float32(srcHeight) / 2
	halfDstWidth := float32(dstWidth) / 2
	halfDstHeight := float32(dstHeight) / 2

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, dstb.Min.Y, dstb.Max.Y, 1, func(start, end int) {
		for yi, y := start, -halfDstHeight+float32(start); yi < end; yi, y = yi+1, y+1 {
			for xi, x := dstb.Min.X, -halfDstWidth; xi < dstb.Max.X; xi, x = xi+1, x+1 {
				x2 := cosine*x - sine*y + halfSrcWidth
				y2 := sine*x + cosine*y + halfSrcHeight

				var rgba pixel

				switch f.interpolation {
				default:
				case NearestNeighborInterpolation:
					rgba = nearestNeighbor(pixGetter, x2, y2)
				case BilinearInterpolation:
					rgba = bilinearInterpolation(pixGetter, x2, y2)
				case BicubicInterpolation:
					rgba = bicubicInterpolation(pixGetter, x2, y2, -0.75)
				}

				pixSetter.setPixel(xi, yi, rgba)
			}
		}
	})
}

func (f *rotateFilter) CanMerge(filter Filter) bool {
	switch filt := filter.(type) {
	case *rotateFilter:
		return true
	case *transformFilter:
		switch filt.transformer {
		case Rotate90Transformer, Rotate180Transformer, Rotate270Transformer:
			return true
		}
	}

	return false
}

func (f *rotateFilter) Merge(filter Filter) {
	switch filt := filter.(type) {
	case *rotateFilter:
		f.rad = gm32.Mod(f.rad+filt.rad, 2*math.Pi)
		filt.oldInterpolation = f.interpolation
		f.interpolation = filt.interpolation
	case *transformFilter:
		switch filt.transformer {
		case Rotate90Transformer:
			f.rad = gm32.Mod(f.rad+math.Pi/2, 2*math.Pi)
		case Rotate180Transformer:
			f.rad = gm32.Mod(f.rad+math.Pi, 2*math.Pi)
		case Rotate270Transformer:
			f.rad = gm32.Mod(f.rad+3*math.Pi/2, 2*math.Pi)
		default:
			return
		}
	default:
		return
	}

	f.mergeCount++
}

func (f *rotateFilter) Skip() bool {
	return f.rad == 0
}

func (f *rotateFilter) CanUndo(filter Filter) bool {
	switch filt := filter.(type) {
	case *rotateFilter:
		return true
	case *transformFilter:
		switch filt.transformer {
		case Rotate90Transformer, Rotate180Transformer, Rotate270Transformer:
			return true
		}
	}

	return false
}

func (f *rotateFilter) Undo(filter Filter) bool {
	switch filt := filter.(type) {
	case *rotateFilter:
		f.rad = gm32.Mod(f.rad-filt.rad, 2*math.Pi)
		f.interpolation = filt.oldInterpolation
	case *transformFilter:
		switch filt.transformer {
		case Rotate90Transformer:
			f.rad = gm32.Mod(f.rad-math.Pi/2, 2*math.Pi)
		case Rotate180Transformer:
			f.rad = gm32.Mod(f.rad-math.Pi, 2*math.Pi)
		case Rotate270Transformer:
			f.rad = gm32.Mod(f.rad-3*math.Pi/2, 2*math.Pi)
		default:
			return false
		}
	default:
		return false
	}

	f.mergeCount--
	return f.mergeCount == 0
}

func (f *rotateFilter) Copy() Filter {
	return &rotateFilter{
		rad:              f.rad,
		interpolation:    f.interpolation,
		oldInterpolation: f.oldInterpolation,
		mergeCount:       f.mergeCount,
	}
}

// Rotates the image using given interpolation method.
// The angle is given in radians.
func Rotate(rad float32, interpolation Interpolation) MergingFilter {
	if gm32.Mod(rad, 2*math.Pi) == 0 {
		return nil
	}

	return &rotateFilter{
		rad:              rad,
		interpolation:    interpolation,
		oldInterpolation: interpolation,
		mergeCount:       1,
	}
}
