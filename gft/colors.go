package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/giu/gcu"
	"github.com/infastin/gul/gm32"
	"github.com/infastin/gul/tools"
)

type colorchanFilter struct {
	fn   func(float32, interface{}) float32
	lut  []float32
	data interface{}
}

func (f *colorchanFilter) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (f *colorchanFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	lutSize := len(f.lut)
	neededLutSize := 0

	switch pixGetter.img.(type) {
	case *image.RGBA, *image.NRGBA, *image.YCbCr, *image.Gray, *image.CMYK:
		neededLutSize = 0xff + 1
	default:
		neededLutSize = 0xffff + 1
	}

	if lutSize != neededLutSize {
		f.makeLut(neededLutSize)
	}

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.Y, srcb.Max.Y, 1, func(start, end int) {
		for y := start; y < end; y++ {
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				pix := pixGetter.getPixel(x, y)
				pix.r = f.getFromLut(pix.r)
				pix.g = f.getFromLut(pix.g)
				pix.b = f.getFromLut(pix.b)
				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pix)
			}
		}
	})
}

func (f *colorchanFilter) Merge(Filter) bool {
	return false
}

func (f *colorchanFilter) Undo(Filter) bool {
	return false
}

func (f *colorchanFilter) Skip() bool {
	return false
}

func (f *colorchanFilter) Copy() Filter {
	lut := append([]float32{}, f.lut...)
	cpy := &colorchanFilter{
		fn:   f.fn,
		lut:  lut,
		data: f.data,
	}
	return cpy
}

func (f *colorchanFilter) makeLut(lutSize int) {
	lutLen := len(f.lut)
	start := 0

	if lutLen == 0 {
		f.lut = make([]float32, lutSize)
	} else if lutLen < lutSize {
		newLut := make([]float32, lutSize)
		copy(newLut[:lutLen], f.lut)
		f.lut = newLut
		start = lutLen
	}

	q := float32(1) / float32(lutSize-1)
	for i := start; i < lutSize; i++ {
		v := float32(i) * q
		f.lut[i] = f.fn(v, f.data)
	}
}

func (f *colorchanFilter) getFromLut(x float32) float32 {
	i := int(gm32.Round(x * float32(len(f.lut)-1)))
	return f.lut[i]
}

type invertFilter struct {
	colorchanFilter
	mergeCount uint
	state      byte
}

func (f *invertFilter) Skip() bool {
	return f.state == 0
}

func (f *invertFilter) Merge(filter Filter) bool {
	filt := filter.(*invertFilter)
	f.state ^= filt.state
	f.mergeCount++
	return true
}

func (f *invertFilter) Undo(filter Filter) bool {
	filt := filter.(*invertFilter)
	f.state = f.state ^ filt.state
	f.mergeCount--
	return f.mergeCount == 0
}

func (f *invertFilter) Copy() Filter {
	return &invertFilter{
		colorchanFilter: f.colorchanFilter,
		state:           f.state,
		mergeCount:      f.mergeCount,
	}
}

func Invert() Filter {
	return &invertFilter{
		colorchanFilter: colorchanFilter{
			fn: func(f float32, _ interface{}) float32 {
				return 1 - f
			},
		},
		state:      1,
		mergeCount: 1,
	}
}

type colorFilter struct {
	fn   func(pixel, interface{}) pixel
	data interface{}
}

func (f *colorFilter) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (f *colorFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.Y, srcb.Max.Y, 1, func(start, end int) {
		for y := start; y < end; y++ {
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				pix := pixGetter.getPixel(x, y)
				pix = f.fn(pix, f.data)
				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pix)
			}
		}
	})
}

func (f *colorFilter) Merge(Filter) bool {
	return false
}

func (f *colorFilter) Undo(Filter) bool {
	return false
}

func (f *colorFilter) Skip() bool {
	return false
}

func (f *colorFilter) Copy() Filter {
	return &colorFilter{
		fn:   f.fn,
		data: f.data,
	}
}

type grayscaleFilter struct {
	colorFilter
	mergeCount uint
}

func (f *grayscaleFilter) Merge(Filter) bool {
	f.mergeCount++
	return true
}

func (f *grayscaleFilter) Undo(Filter) bool {
	f.mergeCount--
	return f.mergeCount == 0
}

func (f *grayscaleFilter) Copy() Filter {
	return &grayscaleFilter{
		colorFilter: f.colorFilter,
		mergeCount:  f.mergeCount,
	}
}

func Grayscale() Filter {
	return &grayscaleFilter{
		colorFilter: colorFilter{
			fn: func(pix pixel, _ interface{}) pixel {
				v := 0.299*pix.r + 0.587*pix.g + 0.114*pix.b
				return pixel{v, v, v, pix.a}
			},
		},
		mergeCount: 1,
	}
}

type brightnessFilter struct {
	colorFilter
	mergeCount uint
}

func (f *brightnessFilter) Merge(filter Filter) bool {
	filt := filter.(*brightnessFilter)

	r1 := f.data.(float32)
	r2 := filt.data.(float32)

	r := gm32.Clamp(r1+r2, -1, 1)
	f.data = r

	f.mergeCount++

	return true
}

func (f *brightnessFilter) Undo(filter Filter) bool {
	filt := filter.(*brightnessFilter)

	r1 := f.data.(float32)
	r2 := filt.data.(float32)

	r := gm32.Clamp(r1-r2, -1, 1)
	f.data = r

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *brightnessFilter) Skip() bool {
	r := f.data.(float32)
	return r == 0
}

func (f *brightnessFilter) Copy() Filter {
	return &brightnessFilter{
		colorFilter: f.colorFilter,
		mergeCount:  f.mergeCount,
	}
}

func Brightness(ratio float32) Filter {
	if ratio == 0 {
		return nil
	}

	return &brightnessFilter{
		colorFilter: colorFilter{
			fn: func(pix pixel, data interface{}) pixel {
				rat := data.(float32)
				h, s, v := gcu.RGBToHSV(pix.r, pix.g, pix.b)
				v = gm32.Clamp(v+rat, 0, 1)
				r, g, b := gcu.HSVToRGB(h, s, v)
				return pixel{r, g, b, pix.a}
			},
			data: ratio,
		},
		mergeCount: 1,
	}
}

type saturationFilter struct {
	colorFilter
	mergeCount uint
}

func (f *saturationFilter) Merge(filter Filter) bool {
	filt := filter.(*saturationFilter)

	r1 := f.data.(float32)
	r2 := filt.data.(float32)

	r := gm32.Clamp(r1+r2, -1, 1)
	f.data = r

	f.mergeCount++

	return true
}

func (f *saturationFilter) Undo(filter Filter) bool {
	filt := filter.(*saturationFilter)

	r1 := f.data.(float32)
	r2 := filt.data.(float32)

	r := gm32.Clamp(r1-r2, -1, 1)
	f.data = r

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *saturationFilter) Skip() bool {
	r := f.data.(float32)
	return r == 0
}

func (f *saturationFilter) Copy() Filter {
	return &saturationFilter{
		colorFilter: f.colorFilter,
		mergeCount:  f.mergeCount,
	}
}

func Saturation(ratio float32) Filter {
	if ratio == 0 {
		return nil
	}

	return &saturationFilter{
		colorFilter: colorFilter{
			fn: func(pix pixel, data interface{}) pixel {
				rat := data.(float32)
				h, s, v := gcu.RGBToHSV(pix.r, pix.g, pix.b)
				s = gm32.Clamp(s+rat, 0, 1)
				r, g, b := gcu.HSVToRGB(h, s, v)
				return pixel{r, g, b, pix.a}
			},
			data: ratio,
		},
		mergeCount: 1,
	}
}

type hueFilter struct {
	colorFilter
	mergeCount uint
}

func (f *hueFilter) Merge(filter Filter) bool {
	filt := filter.(*hueFilter)

	r1 := f.data.(float32)
	r2 := filt.data.(float32)

	r := gm32.Clamp(r1+r2, -1, 1)
	f.data = r

	f.mergeCount++

	return true
}

func (f *hueFilter) Undo(filter Filter) bool {
	filt := filter.(*hueFilter)

	r1 := f.data.(float32)
	r2 := filt.data.(float32)

	r := gm32.Clamp(r1-r2, -1, 1)
	f.data = r

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *hueFilter) Skip() bool {
	r := f.data.(float32)
	return r == 0
}

func (f *hueFilter) Copy() Filter {
	return &hueFilter{
		colorFilter: f.colorFilter,
		mergeCount:  f.mergeCount,
	}
}

func Hue(ratio float32) Filter {
	if ratio == 0 {
		return nil
	}

	return &hueFilter{
		colorFilter: colorFilter{
			fn: func(pix pixel, data interface{}) pixel {
				rat := data.(float32)
				h, s, v := gcu.RGBToHSV(pix.r, pix.g, pix.b)
				h = gm32.Clamp(h+rat, 0, 1)
				r, g, b := gcu.HSVToRGB(h, s, v)
				return pixel{r, g, b, pix.a}
			},
			data: ratio,
		},
		mergeCount: 1,
	}
}

type hsbDiff struct {
	h, s, b float32
}

type hsbFilter struct {
	colorFilter
	mergeCount uint
}

func (f *hsbFilter) Merge(filter Filter) bool {
	filt := filter.(*hsbFilter)

	hsb1 := f.data.(hsbDiff)
	hsb2 := filt.data.(hsbDiff)

	h := gm32.Clamp(hsb1.h+hsb2.h, -1, 1)
	s := gm32.Clamp(hsb1.s+hsb2.s, -1, 1)
	b := gm32.Clamp(hsb1.b+hsb2.b, -1, 1)

	f.data = hsbDiff{h, s, b}
	f.mergeCount++

	return true
}

func (f *hsbFilter) Undo(filter Filter) bool {
	filt := filter.(*hsbFilter)

	hsb1 := f.data.(hsbDiff)
	hsb2 := filt.data.(hsbDiff)

	h := gm32.Clamp(hsb1.h-hsb2.h, -1, 1)
	s := gm32.Clamp(hsb1.s-hsb2.s, -1, 1)
	b := gm32.Clamp(hsb1.b-hsb2.b, -1, 1)

	f.data = hsbDiff{h, s, b}
	f.mergeCount--

	return f.mergeCount == 0
}

func (f *hsbFilter) Skip() bool {
	hsb := f.data.(hsbDiff)
	return hsb.h == 0 && hsb.s == 0 && hsb.b == 0
}

func (f *hsbFilter) Copy() Filter {
	return &hsbFilter{
		colorFilter: f.colorFilter,
		mergeCount:  f.mergeCount,
	}
}

func HSB(h, s, b float32) Filter {
	if h == 0 && s == 0 && b == 0 {
		return nil
	}

	hsb := hsbDiff{h, s, b}

	return &hsbFilter{
		colorFilter: colorFilter{
			fn: func(pix pixel, data interface{}) pixel {
				hsb := data.(hsbDiff)
				h1, s1, v1 := gcu.RGBToHSV(pix.r, pix.g, pix.b)

				h2 := gm32.Clamp(h1+hsb.h, 0, 1)
				s2 := gm32.Clamp(s1+hsb.s, 0, 1)
				v2 := gm32.Clamp(v1+hsb.b, 0, 1)

				r, g, b := gcu.HSVToRGB(h2, s2, v2)
				return pixel{r, g, b, pix.a}
			},
			data: hsb,
		},
		mergeCount: 1,
	}
}
