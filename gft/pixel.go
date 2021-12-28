package gft

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/infastin/gul/gm32"
)

type pixel struct {
	r, g, b, a float32
}

func (p pixel) clampPixel(min, max float32) pixel {
	p.r = gm32.Clamp(p.r, min, max)
	p.g = gm32.Clamp(p.g, min, max)
	p.b = gm32.Clamp(p.b, min, max)
	p.a = gm32.Clamp(p.a, min, max)
	return p
}

type pixelGetter struct {
	img    image.Image
	bounds image.Rectangle
}

const (
	qf8  = 1.0 / 0xff
	qf16 = 1.0 / 0xffff
	epal = qf16 * qf16 / 2
)

func pixelFromColor(c color.Color) (pix pixel) {
	r, g, b, a := c.RGBA()
	switch a {
	case 0:
		pix = pixel{0, 0, 0, 0}
	case 0xffff:
		pix = pixel{
			r: float32(r) * qf16,
			g: float32(g) * qf16,
			b: float32(b) * qf16,
			a: 1,
		}
	default:
		q := float32(1) / float32(a)
		pix = pixel{
			r: float32(r) * q,
			g: float32(g) * q,
			b: float32(b) * q,
			a: float32(a) * qf16,
		}
	}

	return pix
}

func newPixelGetter(img image.Image) *pixelGetter {
	pixGetter := &pixelGetter{
		img:    img,
		bounds: img.Bounds(),
	}

	return pixGetter
}

func (p *pixelGetter) getPixel(x, y int) pixel {
	if !(image.Point{x, y}.In(p.bounds)) {
		return pixel{0, 0, 0, 0}
	}

	switch img := p.img.(type) {
	case *image.RGBA:
		i := img.PixOffset(x, y)
		a := img.Pix[i+3]
		switch a {
		case 0:
			return pixel{0, 0, 0, 0}
		case 0xff:
			return pixel{
				r: float32(img.Pix[i]) * qf8,
				g: float32(img.Pix[i+1]) * qf8,
				b: float32(img.Pix[i+2]) * qf8,
				a: 1,
			}
		default:
			q := float32(1) / float32(a)
			return pixel{
				r: float32(img.Pix[i]) * q,
				g: float32(img.Pix[i+1]) * q,
				b: float32(img.Pix[i+2]) * q,
				a: float32(a) * qf8,
			}
		}
	case *image.RGBA64:
		i := img.PixOffset(x, y)
		a := uint16(img.Pix[i+6])<<8 | uint16(img.Pix[i+7])
		switch a {
		case 0:
			return pixel{0, 0, 0, 0}
		case 0xffff:
			return pixel{
				r: float32(uint16(img.Pix[i])<<8|uint16(img.Pix[i+1])) * qf16,
				g: float32(uint16(img.Pix[i+2])<<8|uint16(img.Pix[i+3])) * qf16,
				b: float32(uint16(img.Pix[i+4])<<8|uint16(img.Pix[i+5])) * qf16,
				a: 1,
			}
		default:
			q := float32(1) / float32(a)
			return pixel{
				r: float32(uint16(img.Pix[i])<<8|uint16(img.Pix[i+1])) * q,
				g: float32(uint16(img.Pix[i+2])<<8|uint16(img.Pix[i+3])) * q,
				b: float32(uint16(img.Pix[i+4])<<8|uint16(img.Pix[i+5])) * q,
				a: float32(a) * qf16,
			}
		}
	case *image.NRGBA:
		i := img.PixOffset(x, y)
		return pixel{
			r: float32(img.Pix[i]) * qf8,
			g: float32(img.Pix[i+1]) * qf8,
			b: float32(img.Pix[i+2]) * qf8,
			a: float32(img.Pix[i+3]) * qf8,
		}
	case *image.NRGBA64:
		i := img.PixOffset(x, y)
		return pixel{
			r: float32(uint16(img.Pix[i])<<8|uint16(img.Pix[i+1])) * qf16,
			g: float32(uint16(img.Pix[i+2])<<8|uint16(img.Pix[i+3])) * qf16,
			b: float32(uint16(img.Pix[i+4])<<8|uint16(img.Pix[i+5])) * qf16,
			a: float32(uint16(img.Pix[i+6])<<8|uint16(img.Pix[i+7])) * qf16,
		}
	case *image.Gray:
		i := img.PixOffset(x, y)
		v := float32(img.Pix[i]) * qf8
		return pixel{v, v, v, 1}
	case *image.Gray16:
		i := img.PixOffset(x, y)
		v := float32(uint16(img.Pix[i])<<8|uint16(img.Pix[i+1])) * qf16
		return pixel{v, v, v, 1}
	default:
		return pixelFromColor(p.img.At(x, y))
	}
}

func (p *pixelGetter) average(xmin, ymin, xmax, ymax int) pixel {
	if xmin >= p.bounds.Max.X || ymin >= p.bounds.Max.Y {
		return pixel{0, 0, 0, 0}
	}

	if xmax >= p.bounds.Max.X {
		xmax = p.bounds.Max.X - 1
	}

	if ymax >= p.bounds.Max.Y {
		ymax = p.bounds.Max.Y - 1
	}

	diffX := xmax - xmin + 1
	diffY := ymax - ymin + 1
	pixNum := float32(diffX * diffY)

	avg := pixel{}
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			pix := p.getPixel(x, y)
			avg.r += pix.r
			avg.g += pix.g
			avg.b += pix.b
			avg.a += pix.a
		}
	}

	avg.r /= pixNum
	avg.g /= pixNum
	avg.b /= pixNum
	avg.a /= pixNum

	return avg
}

func (p *pixelGetter) getPixelRow(y int, buf *[]pixel) {
	*buf = (*buf)[:0]
	for x := p.bounds.Min.X; x < p.bounds.Max.X; x++ {
		*buf = append(*buf, p.getPixel(x, y))
	}
}

func (p *pixelGetter) getPixelColumn(x int, buf *[]pixel) {
	*buf = (*buf)[:0]
	for y := p.bounds.Min.Y; y < p.bounds.Max.Y; y++ {
		*buf = append(*buf, p.getPixel(x, y))
	}
}

type pixelSetter struct {
	img    draw.Image
	bounds image.Rectangle
}

func newPixelSetter(img draw.Image) *pixelSetter {
	pixSetter := &pixelSetter{
		img:    img,
		bounds: img.Bounds(),
	}

	return pixSetter
}

func (p *pixelSetter) setPixel(x, y int, pix pixel) {
	if !(image.Point{x, y}.In(p.bounds)) {
		return
	}

	switch img := p.img.(type) {
	case *image.RGBA:
		fa := pix.a * 0xff
		i := img.PixOffset(x, y)
		img.Pix[i] = f32u8(pix.r * fa)
		img.Pix[i+1] = f32u8(pix.g * fa)
		img.Pix[i+2] = f32u8(pix.b * fa)
		img.Pix[i+3] = f32u8(fa)
	case *image.RGBA64:
		fa := pix.a * 0xffff
		i := img.PixOffset(x, y)

		r16 := f32u16(pix.r * fa)
		g16 := f32u16(pix.g * fa)
		b16 := f32u16(pix.b * fa)
		a16 := f32u16(fa)

		img.Pix[i] = uint8(r16 >> 8)
		img.Pix[i+1] = uint8(r16 & 8)
		img.Pix[i+2] = uint8(g16 >> 8)
		img.Pix[i+3] = uint8(g16 & 8)
		img.Pix[i+4] = uint8(b16 >> 8)
		img.Pix[i+5] = uint8(b16 & 8)
		img.Pix[i+6] = uint8(a16 & 8)
		img.Pix[i+7] = uint8(a16 & 8)
	case *image.NRGBA:
		i := img.PixOffset(x, y)
		img.Pix[i] = f32u8(pix.r * 0xff)
		img.Pix[i+1] = f32u8(pix.g * 0xff)
		img.Pix[i+2] = f32u8(pix.b * 0xff)
		img.Pix[i+3] = f32u8(pix.a * 0xff)
	case *image.NRGBA64:
		i := img.PixOffset(x, y)

		r16 := f32u16(pix.r * 0xffff)
		g16 := f32u16(pix.g * 0xffff)
		b16 := f32u16(pix.b * 0xffff)
		a16 := f32u16(pix.a * 0xffff)

		img.Pix[i] = uint8(r16 >> 8)
		img.Pix[i+1] = uint8(r16 & 8)
		img.Pix[i+2] = uint8(g16 >> 8)
		img.Pix[i+3] = uint8(g16 & 8)
		img.Pix[i+4] = uint8(b16 >> 8)
		img.Pix[i+5] = uint8(b16 & 8)
		img.Pix[i+6] = uint8(a16 & 8)
		img.Pix[i+7] = uint8(a16 & 8)
	case *image.Gray:
		i := img.PixOffset(x, y)
		img.Pix[i] = f32u8((0.299*pix.r + 0.587*pix.g + 0.114*pix.b) * pix.a * 0xff)
	case *image.Gray16:
		i := img.PixOffset(x, y)
		v := f32u16((0.299*pix.r + 0.587*pix.g + 0.114*pix.b) * pix.a * 0xffff)
		img.Pix[i] = uint8(v >> 8)
		img.Pix[i+1] = uint8(v & 0xff)
	default:
		r := f32u16(pix.r * 0xffff)
		g := f32u16(pix.g * 0xffff)
		b := f32u16(pix.b * 0xffff)
		a := f32u16(pix.a * 0xffff)
		p.img.Set(x, y, color.NRGBA64{r, g, b, a})
	}
}

func (p *pixelSetter) setPixelRow(y int, buf []pixel) {
	for i, x := 0, p.bounds.Min.X; x < len(buf); i, x = i+1, x+1 {
		p.setPixel(x, y, buf[i])
	}
}

func (p *pixelSetter) setPixelColumn(x int, buf []pixel) {
	for i, y := 0, p.bounds.Min.Y; y < len(buf); i, y = i+1, y+1 {
		p.setPixel(x, y, buf[i])
	}
}

func bilinearInterpolation(pixGetter *pixelGetter, x, y float32) pixel {
	xmin := int(gm32.Floor(x))
	ymin := int(gm32.Floor(y))
	xmax := xmin + 1
	ymax := ymin + 1

	p00 := pixGetter.getPixel(xmin, ymin)
	p01 := pixGetter.getPixel(xmax, ymin)
	p10 := pixGetter.getPixel(xmin, ymax)
	p11 := pixGetter.getPixel(xmax, ymax)

	fx := x - float32(xmin)
	fy := y - float32(ymin)

	r := gm32.InterpolateBilinear(p00.r, p01.r, p10.r, p11.r, fx, fy)
	g := gm32.InterpolateBilinear(p00.g, p01.g, p10.g, p11.g, fx, fy)
	b := gm32.InterpolateBilinear(p00.b, p01.b, p10.b, p11.b, fx, fy)
	a := gm32.InterpolateBilinear(p00.a, p01.a, p10.a, p11.a, fx, fy)

	return pixel{r, g, b, a}
}

func bicubicInterpolation(pixGetter *pixelGetter, x, y, a float32) pixel {
	x1 := int(gm32.Floor(x))
	x0 := x1 - 1
	x2 := x1 + 1
	x3 := x1 + 2

	y1 := int(gm32.Floor(y))
	y0 := y1 - 1
	y2 := y1 + 1
	y3 := y1 + 2

	fx0 := x - float32(x0)
	fx1 := x - float32(x1)
	fx2 := float32(x2) - x
	fx3 := float32(x3) - x

	fy0 := y - float32(y0)
	fy1 := y - float32(y1)
	fy2 := float32(y2) - y
	fy3 := float32(y3) - y

	p00 := pixGetter.getPixel(x0, y0)
	p01 := pixGetter.getPixel(x1, y0)
	p02 := pixGetter.getPixel(x2, y0)
	p03 := pixGetter.getPixel(x3, y0)

	p10 := pixGetter.getPixel(x0, y1)
	p11 := pixGetter.getPixel(x1, y1)
	p12 := pixGetter.getPixel(x2, y1)
	p13 := pixGetter.getPixel(x3, y1)

	p20 := pixGetter.getPixel(x0, y2)
	p21 := pixGetter.getPixel(x1, y2)
	p22 := pixGetter.getPixel(x2, y2)
	p23 := pixGetter.getPixel(x3, y2)

	p30 := pixGetter.getPixel(x0, y3)
	p31 := pixGetter.getPixel(x1, y3)
	p32 := pixGetter.getPixel(x2, y3)
	p33 := pixGetter.getPixel(x3, y3)

	xMat := gm32.NewMat(1, 4)(
		fx0, fx1, fx2, fx3,
	)

	yMat := gm32.NewMat(4, 1)(
		fy0, fy1, fy2, fy3,
	)

	redMat := gm32.NewMat(4, 4)(
		p00.r, p01.r, p02.r, p03.r,
		p10.r, p11.r, p12.r, p13.r,
		p20.r, p21.r, p22.r, p23.r,
		p30.r, p31.r, p32.r, p33.r,
	)

	greenMat := gm32.NewMat(4, 4)(
		p00.g, p01.g, p02.g, p03.g,
		p10.g, p11.g, p12.g, p13.g,
		p20.g, p21.g, p22.g, p23.g,
		p30.g, p31.g, p32.g, p33.g,
	)

	blueMat := gm32.NewMat(4, 4)(
		p00.b, p01.b, p02.b, p03.b,
		p10.b, p11.b, p12.b, p13.b,
		p20.b, p21.b, p22.b, p23.b,
		p30.b, p31.b, p32.b, p33.b,
	)

	alphaMat := gm32.NewMat(4, 4)(
		p00.a, p01.a, p02.a, p03.a,
		p10.a, p11.a, p12.a, p13.a,
		p20.a, p21.a, p22.a, p23.a,
		p30.a, p31.a, p32.a, p33.a,
	)

	red := gm32.InterpolateBicubic(xMat, redMat, yMat, a)
	green := gm32.InterpolateBicubic(xMat, greenMat, yMat, a)
	blue := gm32.InterpolateBicubic(xMat, blueMat, yMat, a)
	alpha := gm32.InterpolateBicubic(xMat, alphaMat, yMat, a)

	return pixel{red, green, blue, alpha}
}

func nearestNeighbor(pixGetter *pixelGetter, x, y float32) pixel {
	xmin := int(gm32.Round(x))
	ymin := int(gm32.Round(y))
	return pixGetter.getPixel(xmin, ymin)
}
