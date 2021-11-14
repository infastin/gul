package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/giu/gcu"
	"github.com/infastin/gul/gm32"
	"github.com/infastin/gul/tools"
)

type colorchanFilter struct {
	fn     func(float32, interface{}) float32
	lut    []float32
	useLut bool
	data   interface{}
}

func (f *colorchanFilter) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (f *colorchanFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	useLut := false
	if f.useLut {
		lutSize := len(f.lut)
		neededLutSize := 0

		switch pixGetter.img.(type) {
		case *image.RGBA, *image.NRGBA, *image.YCbCr, *image.Gray, *image.CMYK:
			neededLutSize = 0xff + 1
		default:
			neededLutSize = 0xffff + 1
		}

		numCalc := srcb.Dx() * srcb.Dy() * 3
		if numCalc > neededLutSize*2 {
			if lutSize != neededLutSize {
				f.makeLut(neededLutSize)
			}

			useLut = true
		}
	}

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.Y, srcb.Max.Y, 1, func(start, end int) {
		for y := start; y < end; y++ {
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				pix := pixGetter.getPixel(x, y)

				if useLut {
					pix.r = f.getFromLut(pix.r)
					pix.g = f.getFromLut(pix.g)
					pix.b = f.getFromLut(pix.b)
				} else {
					pix.r = f.fn(pix.r, f.data)
					pix.g = f.fn(pix.g, f.data)
					pix.b = f.fn(pix.b, f.data)
				}

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

func (f *invertFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*invertFilter)
	if !ok {
		return false
	}

	f.state ^= filt.state
	f.mergeCount++
	return true
}

func (f *invertFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*invertFilter)
	if !ok {
		return false
	}

	f.state = f.state ^ filt.state
	f.mergeCount--
	return f.mergeCount == 0
}

func (f *invertFilter) Skip() bool {
	return f.state == 0
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
			useLut: false,
		},
		state:      1,
		mergeCount: 1,
	}
}

type contrastFilter struct {
	colorchanFilter
	mergeCount uint
}

func (f *contrastFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*contrastFilter)
	if !ok {
		return false
	}

	r1 := f.data.(float32)
	r2 := filt.data.(float32)
	r := gm32.Clamp(r1+r2, -1, 1)

	f.data = r
	f.mergeCount++

	return true
}

func (f *contrastFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*contrastFilter)
	if !ok {
		return false
	}

	r1 := f.data.(float32)
	r2 := filt.data.(float32)
	r := gm32.Clamp(r1-r2, -1, 1)

	f.data = r
	f.mergeCount--

	return f.mergeCount == 0
}

func (f *contrastFilter) Skip() bool {
	r := f.data.(float32)
	return r == 0
}

func (f *contrastFilter) Copy() Filter {
	return &contrastFilter{
		colorchanFilter: f.colorchanFilter,
		mergeCount:      f.mergeCount,
	}
}

func Contrast(ratio float32) Filter {
	if ratio == 0 {
		return nil
	}

	return &contrastFilter{
		colorchanFilter: colorchanFilter{
			fn: func(f float32, data interface{}) float32 {
				rat := data.(float32)
				alpha := gm32.Clamp(rat, -1, 1) + 1
				c := alpha*(f-0.5) + 0.5

				return gm32.Clamp(c, 0, 1)
			},
			useLut: false,
			data:   ratio,
		},
		mergeCount: 1,
	}
}

type brightnessFilter struct {
	colorchanFilter
	mergeCount uint
}

func (f *brightnessFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*brightnessFilter)
	if !ok {
		return false
	}

	r1 := f.data.(float32)
	r2 := filt.data.(float32)
	r := gm32.Clamp(r1+r2, -1, 1)

	f.data = r
	f.mergeCount++

	return true
}

func (f *brightnessFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*brightnessFilter)
	if !ok {
		return false
	}

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
		colorchanFilter: f.colorchanFilter,
		mergeCount:      f.mergeCount,
	}
}

func Brightness(ratio float32) Filter {
	if ratio == 0 {
		return nil
	}

	return &brightnessFilter{
		colorchanFilter: colorchanFilter{
			fn: func(f float32, data interface{}) float32 {
				rat := data.(float32)
				beta := gm32.Clamp(rat, -1, 1)

				return gm32.Clamp(f+beta, 0, 1)
			},
			useLut: false,
			data:   ratio,
		},
		mergeCount: 1,
	}
}

type brightnessContrastVal struct {
	bratio, cratio float32
}

type brightnessContrastFilter struct {
	colorchanFilter
	mergeCount uint
}

func (f *brightnessContrastFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*brightnessContrastFilter)
	if !ok {
		return false
	}

	v1 := f.data.(brightnessContrastVal)
	v2 := filt.data.(brightnessContrastVal)

	br := gm32.Clamp(v1.bratio+v2.bratio, -1, 1)
	cr := gm32.Clamp(v1.cratio+v2.cratio, -1, 1)

	f.data = brightnessContrastVal{br, cr}
	f.mergeCount++

	return true
}

func (f *brightnessContrastFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*brightnessContrastFilter)
	if !ok {
		return false
	}

	v1 := f.data.(brightnessContrastVal)
	v2 := filt.data.(brightnessContrastVal)

	br := gm32.Clamp(v1.bratio-v2.bratio, -1, 1)
	cr := gm32.Clamp(v1.cratio-v2.cratio, -1, 1)

	f.data = brightnessContrastVal{br, cr}
	f.mergeCount--

	return f.mergeCount == 0
}

func (f *brightnessContrastFilter) Skip() bool {
	v := f.data.(brightnessContrastVal)
	return v.bratio == 0 && v.cratio == 0
}

func (f *brightnessContrastFilter) Copy() Filter {
	return &brightnessContrastFilter{
		colorchanFilter: f.colorchanFilter,
		mergeCount:      f.mergeCount,
	}
}

func BrightnessContrast(bratio, cratio float32) Filter {
	if bratio == 0 && cratio == 0 {
		return nil
	}

	val := brightnessContrastVal{bratio, cratio}

	return &brightnessContrastFilter{
		colorchanFilter: colorchanFilter{
			fn: func(f float32, data interface{}) float32 {
				v := data.(brightnessContrastVal)
				alpha := gm32.Clamp(v.cratio, -1, 1) + 1
				beta := gm32.Clamp(v.bratio, -1, 1)
				c := alpha*(f-0.5) + 0.5 + beta

				return gm32.Clamp(c, 0, 1)
			},
			data:   val,
			useLut: false,
		},
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

func (f *grayscaleFilter) Merge(filter Filter) bool {
	if _, ok := filter.(*grayscaleFilter); !ok {
		return false
	}

	f.mergeCount++
	return true
}

func (f *grayscaleFilter) Undo(filter Filter) bool {
	if _, ok := filter.(*grayscaleFilter); !ok {
		return false
	}

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

type sepiaFilter struct {
	colorFilter
	mergeCount uint
}

func (f *sepiaFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*sepiaFilter)
	if !ok {
		return false
	}

	r1 := f.data.(float32)
	r2 := filt.data.(float32)
	r := gm32.Clamp(r1+r2, 0, 1)

	f.data = r
	f.mergeCount++

	return true
}

func (f *sepiaFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*sepiaFilter)
	if !ok {
		return false
	}

	r1 := f.data.(float32)
	r2 := filt.data.(float32)
	r := gm32.Clamp(r1-r2, 0, 1)

	f.data = r
	f.mergeCount--

	return f.mergeCount == 0
}

func (f *sepiaFilter) Skip() bool {
	r := f.data.(float32)
	return r == 0
}

func (f *sepiaFilter) Copy() Filter {
	return &sepiaFilter{
		colorFilter: f.colorFilter,
		mergeCount:  f.mergeCount,
	}
}

func Sepia(ratio float32) Filter {
	if ratio == 0 {
		return nil
	}

	return &sepiaFilter{
		colorFilter: colorFilter{
			fn: func(pix pixel, data interface{}) pixel {
				rat := data.(float32)
				rat = gm32.Clamp(rat, 0, 1)

				rr := 1 - 0.607*rat
				rg := 0.769 * rat
				rb := 0.189 * rat

				gr := 0.349 * rat
				gg := 1 - 0.314*rat
				gb := 0.168 * rat

				br := 0.272 * rat
				bg := 0.534 * rat
				bb := 1 - 0.869*rat

				r := pix.r*rr + pix.g*rg + pix.b*rb
				g := pix.r*gr + pix.g*gg + pix.b*gb
				b := pix.r*br + pix.g*bg + pix.b*bb

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
	filt, ok := filter.(*hsbFilter)
	if !ok {
		return false
	}

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
	filt, ok := filter.(*hsbFilter)
	if !ok {
		return false
	}

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

type hslDiff struct {
	h, s, l float32
}

type hslFilter struct {
	colorFilter
	mergeCount uint
}

func (f *hslFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*hslFilter)
	if !ok {
		return false
	}

	hsl1 := f.data.(hslDiff)
	hsl2 := filt.data.(hslDiff)

	h := gm32.Clamp(hsl1.h+hsl2.h, -1, 1)
	s := gm32.Clamp(hsl1.s+hsl2.s, -1, 1)
	l := gm32.Clamp(hsl1.l+hsl2.l, -1, 1)

	f.data = hslDiff{h, s, l}
	f.mergeCount++

	return true
}

func (f *hslFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*hslFilter)
	if !ok {
		return false
	}

	hsl1 := f.data.(hslDiff)
	hsl2 := filt.data.(hslDiff)

	h := gm32.Clamp(hsl1.h-hsl2.h, -1, 1)
	s := gm32.Clamp(hsl1.s-hsl2.s, -1, 1)
	l := gm32.Clamp(hsl1.l-hsl2.l, -1, 1)

	f.data = hslDiff{h, s, l}
	f.mergeCount--

	return f.mergeCount == 0
}

func (f *hslFilter) Skip() bool {
	hsl := f.data.(hslDiff)
	return hsl.h == 0 && hsl.s == 0 && hsl.l == 0
}

func (f *hslFilter) Copy() Filter {
	return &hslFilter{
		colorFilter: f.colorFilter,
		mergeCount:  f.mergeCount,
	}
}

func HSL(h, s, l float32) Filter {
	if h == 0 && s == 0 && l == 0 {
		return nil
	}

	hsl := hslDiff{h, s, l}

	return &hslFilter{
		colorFilter: colorFilter{
			fn: func(pix pixel, data interface{}) pixel {
				hsl := data.(hslDiff)
				h1, s1, l1 := gcu.RGBToHSL(pix.r, pix.g, pix.b)

				h2 := gm32.Clamp(h1+hsl.h, 0, 1)
				s2 := gm32.Clamp(s1+hsl.s, 0, 1)
				l2 := gm32.Clamp(l1+hsl.l, 0, 1)

				r, g, b := gcu.HSLToRGB(h2, s2, l2)
				return pixel{r, g, b, pix.a}
			},
			data: hsl,
		},
		mergeCount: 1,
	}
}
