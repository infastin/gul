package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/gm32"
	"github.com/infastin/gul/tools"
)

type segment struct {
	Min, Max int
}

func makeSegment(min, max int) segment {
	return segment{min, max}
}

func (l segment) size() int {
	return l.Max - l.Min
}

type contrib struct {
	index  int
	weight float32
}

type contribBounds struct {
	center      float32
	left, right int
}

type resampler struct {
	filtScaleX float32
	filtScaleY float32

	filt ResamplingFilter
}

func newResampler(filt ResamplingFilter, filtScaleX, filtScaleY float32) *resampler {
	return &resampler{
		filt:       filt,
		filtScaleX: filtScaleX,
		filtScaleY: filtScaleY,
	}
}

func (r *resampler) makeCList(dst, src segment, filtScale float32) [][]contrib {
	srcSize := src.size()
	dstSize := dst.size()

	result := make([][]contrib, dstSize)
	cb := make([]contribBounds, dstSize)

	delta := float32(dstSize) / float32(srcSize)
	scale := delta
	if scale > 1 {
		scale = 1
	}

	radius := (r.filt.Support() / scale) * filtScale
	n := 0

	for i := dst.Min; i < dst.Max; i++ {
		center := (float32(i)+0.5)/delta - 0.5
		left := int(gm32.Floor(center - radius))
		if left < 0 {
			left = 0
		}

		right := int(gm32.Ceil(center + radius))
		if right > srcSize-1 {
			right = srcSize - 1
		}

		cb[i-dst.Min] = contribBounds{
			center: center,
			left:   left,
			right:  right,
		}

		n += int(right - left + 1)
	}

	if n == 0 {
		return nil
	}

	ooFiltScale := 1.0 / filtScale
	tmp := make([]contrib, 0, n)

	for i := dst.Min; i < dst.Max; i++ {
		center := cb[i].center
		left := cb[i].left
		right := cb[i].right

		var sum float32
		for j := left; j <= right; j++ {
			weight := r.filt.Kernel((center - float32(j)) * scale * ooFiltScale)
			if weight == 0 {
				continue
			}

			tmp = append(tmp, contrib{
				index:  j,
				weight: weight,
			})
			sum += weight
		}

		for j := range tmp {
			tmp[j].weight /= sum
		}

		result[i-dst.Min] = tmp
		tmp = tmp[len(tmp):]
	}

	return result
}

func (res *resampler) resampleSegment(dst []pixel, src []pixel, clist [][]contrib) {
	for i := 0; i < len(dst); i++ {
		var r, g, b, a float32
		for _, c := range clist[i] {
			col := src[c.index]
			r += col.r * c.weight
			g += col.g * c.weight
			b += col.b * c.weight
			a += col.a * c.weight
		}

		dst[i] = pixel{r, g, b, a}
	}
}

func (r *resampler) resampleX(dst draw.Image, src image.Image, parallel bool) {
	dstb := dst.Bounds()
	srcb := src.Bounds()

	clistx := r.makeCList(
		makeSegment(dstb.Min.X, dstb.Max.X),
		makeSegment(srcb.Min.X, srcb.Max.X),
		r.filtScaleX,
	)

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.Y, srcb.Max.Y, 1, func(start, end int) {
		srcBuf := make([]pixel, srcb.Dx())
		dstBuf := make([]pixel, dstb.Dx())

		for y := start; y < end; y++ {
			pixGetter.getPixelRow(y, &srcBuf)
			r.resampleSegment(dstBuf, srcBuf, clistx)
			pixSetter.setPixelRow(dstb.Min.Y+y-srcb.Min.Y, dstBuf)
		}
	})
}

func (r *resampler) resampleY(dst draw.Image, src image.Image, parallel bool) {
	dstb := dst.Bounds()
	srcb := src.Bounds()

	clisty := r.makeCList(
		makeSegment(dstb.Min.Y, dstb.Max.Y),
		makeSegment(srcb.Min.Y, srcb.Max.Y),
		r.filtScaleY,
	)

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.X, srcb.Max.X, 1, func(start, end int) {
		srcBuf := make([]pixel, srcb.Dy())
		dstBuf := make([]pixel, dstb.Dy())

		for x := start; x < end; x++ {
			pixGetter.getPixelColumn(x, &srcBuf)
			r.resampleSegment(dstBuf, srcBuf, clisty)
			pixSetter.setPixelColumn(dstb.Min.X+x-srcb.Min.X, dstBuf)
		}
	})
}

func resampleNearestNeightbor(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	srcWidth := srcb.Dx()
	srcHeight := srcb.Dy()

	dstWidth := dstb.Dx()
	dstHeight := dstb.Dy()

	scaleX := float32(dstWidth) / float32(srcWidth)
	scaleY := float32(dstHeight) / float32(srcHeight)

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, dstb.Min.Y, dstb.Max.Y, 1, func(start, end int) {
		for yi := start; yi < end; yi++ {
			for xi := dstb.Min.X; xi < dstb.Max.X; xi++ {
				x := float32(xi)/scaleX - float32(srcb.Min.X)
				y := float32(yi)/scaleY - float32(srcb.Min.Y)

				rgba := nearestNeighbor(pixGetter, x, y)
				pixSetter.setPixel(xi, yi, rgba)
			}
		}
	})
}

func (r *resampler) resample(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	srcWidth := srcb.Dx()
	dstWidth := dstb.Dx()

	srcHeight := srcb.Dy()
	dstHeight := dstb.Dy()

	if srcWidth == dstWidth && srcHeight == dstHeight {
		draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)
		return
	}

	if r.filt.Support() <= 0 {
		resampleNearestNeightbor(dst, src, parallel)
		return
	}

	if srcWidth == dstWidth {
		r.resampleY(dst, src, parallel)
		return
	}

	if srcHeight == dstHeight {
		r.resampleX(dst, src, parallel)
		return
	}

	tmp := image.NewRGBA(image.Rect(0, 0, dstWidth, srcHeight))
	r.resampleX(tmp, src, parallel)
	r.resampleY(dst, tmp, parallel)
}
