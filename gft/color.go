package gft

import (
	"image"
	"image/draw"

	"github.com/infastin/gul/giu/gcu"
	"github.com/infastin/gul/gm32"
	"github.com/infastin/gul/tools"
)

type ColorFilter interface {
	Fn(pix pixel) pixel
	Merge(filter ColorFilter)
	Undo(filter ColorFilter) bool
	Skip() bool
	Copy() ColorFilter
	CanMerge(filter ColorFilter) bool
	CanUndo(filter ColorFilter) bool
}

type combineColorFilter struct {
	filters []ColorFilter
}

func (f *combineColorFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*combineColorFilter)
	if !ok {
		return false
	}

	if len(f.filters) != len(filt.filters) {
		return false
	}

	for i := 0; i < len(f.filters); i++ {
		if f.filters[i] == nil || filt.filters[i] == nil {
			continue
		}

		if !f.filters[i].CanMerge(filt.filters[i]) {
			return false
		}
	}

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		if f.filters[i] == nil {
			f.filters[i] = filt.filters[i]
			continue
		}

		f.filters[i].Merge(filt.filters[i])
	}

	return true
}

func (f *combineColorFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*combineColorFilter)
	if !ok {
		return false
	}

	if len(f.filters) != len(filt.filters) {
		return false
	}

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		if f.filters[i] == nil || !f.filters[i].CanUndo(filt.filters[i]) {
			return false
		}
	}

	numTrue := 0
	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		if f.filters[i].Undo(filt.filters[i]) {
			numTrue++
		}
	}

	return numTrue == len(f.filters)
}

func (f *combineColorFilter) Skip() bool {
	for _, filt := range f.filters {
		if filt == nil {
			continue
		}

		if !filt.Skip() {
			return false
		}
	}

	return true
}

func (f *combineColorFilter) Copy() Filter {
	r := &combineColorFilter{}

	r.filters = make([]ColorFilter, len(f.filters))
	for i := 0; i < len(f.filters); i++ {
		if f.filters[i] == nil {
			r.filters[i] = nil
			continue
		}

		r.filters[i] = f.filters[i].Copy()
	}

	return r
}

func (f *combineColorFilter) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (f *combineColorFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
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

				for _, filt := range f.filters {
					pix = filt.Fn(pix)
				}

				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pix)
			}
		}
	})
}

// Creates combination of color filters and returns filter.
func CombineColorFilters(filters ...ColorFilter) Filter {
	if len(filters) == 0 {
		return nil
	}

	numNil := 0
	for _, filt := range filters {
		if filt == nil {
			numNil++
		}
	}

	if numNil == len(filters) {
		return nil
	}

	return &combineColorFilter{
		filters: filters,
	}
}

type colorFilterFunc struct {
	fn func(pix pixel) pixel
}

func (f *colorFilterFunc) CanMerge(ColorFilter) bool {
	return false
}

func (f *colorFilterFunc) Merge(ColorFilter) {}

func (f *colorFilterFunc) CanUndo(ColorFilter) bool {
	return false
}

func (f *colorFilterFunc) Undo(ColorFilter) bool {
	return false
}

func (f *colorFilterFunc) Skip() bool {
	return false
}

func (f *colorFilterFunc) Copy() ColorFilter {
	return &colorFilterFunc{
		fn: f.fn,
	}
}

func (f *colorFilterFunc) Fn(pix pixel) pixel {
	return f.fn(pix)
}

func ColorFilterFunc(fn func(pix pixel) pixel) ColorFilter {
	return &colorFilterFunc{
		fn: fn,
	}
}

type grayscaleFilter struct {
	mergeCount uint
}

func (f *grayscaleFilter) CanMerge(filter ColorFilter) bool {
	if _, ok := filter.(*grayscaleFilter); ok {
		return true
	}

	return false
}

func (f *grayscaleFilter) Merge(filter ColorFilter) {
	f.mergeCount++
}

func (f *grayscaleFilter) Undo(filter ColorFilter) bool {
	f.mergeCount--
	return f.mergeCount == 0
}

func (f *grayscaleFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*grayscaleFilter); ok {
		return true
	}

	return false
}

func (f *grayscaleFilter) Skip() bool {
	return false
}

func (f *grayscaleFilter) Copy() ColorFilter {
	return &grayscaleFilter{
		mergeCount: f.mergeCount,
	}
}

func (f *grayscaleFilter) Fn(pix pixel) pixel {
	v := 0.299*pix.r + 0.587*pix.g + 0.114*pix.b
	return pixel{v, v, v, pix.a}
}

// Grayscales an image.
func Grayscale() ColorFilter {
	return &grayscaleFilter{
		mergeCount: 1,
	}
}

type sepiaFilter struct {
	percentage float32
	mergeCount uint
}

func (f *sepiaFilter) CanMerge(filter ColorFilter) bool {
	if _, ok := filter.(*sepiaFilter); ok {
		return true
	}

	return false
}

func (f *sepiaFilter) Merge(filter ColorFilter) {
	filt := filter.(*sepiaFilter)
	f.percentage = gm32.Clamp(f.percentage+filt.percentage, 0, 100)
	f.mergeCount++
}

func (f *sepiaFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*sepiaFilter); ok {
		return true
	}

	return false
}

func (f *sepiaFilter) Undo(filter ColorFilter) bool {
	filt := filter.(*sepiaFilter)

	f.percentage = gm32.Clamp(f.percentage-filt.percentage, 0, 100)
	f.mergeCount--

	return f.mergeCount == 0
}

func (f *sepiaFilter) Skip() bool {
	return f.percentage == 0
}

func (f *sepiaFilter) Fn(pix pixel) pixel {
	rat := gm32.Clamp(f.percentage, 0, 100) / 100

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
}

func (f *sepiaFilter) Copy() ColorFilter {
	return &sepiaFilter{
		percentage: f.percentage,
		mergeCount: f.mergeCount,
	}
}

// Creates sepia-toned version of an image.
// The percentage parameter specifies how much the image should be adjusted.
// It must be in the range [0, 100].
// It can be any value for merging purposes.
func Sepia(perc float32) ColorFilter {
	if perc == 0 {
		return nil
	}

	return &sepiaFilter{
		percentage: perc,
		mergeCount: 1,
	}
}

type hsbFilter struct {
	h, s, b    float32
	mergeCount uint
}

func (f *hsbFilter) CanMerge(filter ColorFilter) bool {
	if _, ok := filter.(*hsbFilter); ok {
		return true
	}

	return false
}

func (f *hsbFilter) Merge(filter ColorFilter) {
	filt := filter.(*hsbFilter)

	f.h = gm32.Clamp(f.h+filt.h, -360, 360)
	f.s = gm32.Clamp(f.s+filt.s, -100, 100)
	f.b = gm32.Clamp(f.b+filt.b, -100, 100)

	f.mergeCount++
}

func (f *hsbFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*hsbFilter); ok {
		return true
	}

	return false
}

func (f *hsbFilter) Undo(filter ColorFilter) bool {
	filt := filter.(*hsbFilter)

	f.h = gm32.Clamp(f.h-filt.h, -360, 360)
	f.s = gm32.Clamp(f.s-filt.s, -100, 100)
	f.b = gm32.Clamp(f.b-filt.b, -100, 100)

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *hsbFilter) Skip() bool {
	return f.h == 0 && f.s == 0 && f.b == 0
}

func (f *hsbFilter) Copy() ColorFilter {
	return &hsbFilter{
		h:          f.h,
		s:          f.s,
		b:          f.b,
		mergeCount: f.mergeCount,
	}
}

func (f *hsbFilter) Fn(pix pixel) pixel {
	h0 := gm32.Clamp(f.h, -360, 360) / 360
	s0 := gm32.Clamp(f.s, -100, 100) / 100
	b0 := gm32.Clamp(f.b, -100, 100) / 100

	h1, s1, v1 := gcu.RGBToHSV(pix.r, pix.g, pix.b)

	h2 := gm32.Clamp(h1+h0, 0, 1)
	s2 := gm32.Clamp(s1+s0, 0, 1)
	v2 := gm32.Clamp(v1+b0, 0, 1)

	r, g, b := gcu.HSVToRGB(h2, s2, v2)
	return pixel{r, g, b, pix.a}
}

// Changes HSB of each color in the image.
// The hue parameter must be in the range [-360, 360].
// The saturation and brightness parameters must be in the range [-100, 100].
// Each parameter can have any value for merging purposes.
func HSB(h, s, b float32) ColorFilter {
	if h == 0 && s == 0 && b == 0 {
		return nil
	}

	return &hsbFilter{
		h:          h,
		s:          s,
		b:          b,
		mergeCount: 1,
	}
}

type hslFilter struct {
	h, s, l    float32
	mergeCount uint
}

func (f *hslFilter) CanMerge(filter ColorFilter) bool {
	if _, ok := filter.(*hslFilter); ok {
		return true
	}

	return false
}

func (f *hslFilter) Merge(filter ColorFilter) {
	filt := filter.(*hslFilter)

	f.h = gm32.Clamp(f.h+filt.h, -360, 360)
	f.s = gm32.Clamp(f.s+filt.s, -100, 100)
	f.l = gm32.Clamp(f.l+filt.l, -100, 100)

	f.mergeCount++
}

func (f *hslFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*hslFilter); ok {
		return true
	}

	return false
}

func (f *hslFilter) Undo(filter ColorFilter) bool {
	filt := filter.(*hslFilter)

	f.h = gm32.Clamp(f.h-filt.h, -360, 360)
	f.s = gm32.Clamp(f.s-filt.s, -100, 100)
	f.l = gm32.Clamp(f.l-filt.l, -100, 100)

	f.mergeCount--

	return f.mergeCount == 0
}

func (f *hslFilter) Skip() bool {
	return f.h == 0 && f.s == 0 && f.l == 0
}

func (f *hslFilter) Copy() ColorFilter {
	return &hslFilter{
		h:          f.h,
		s:          f.s,
		l:          f.l,
		mergeCount: f.mergeCount,
	}
}

func (f *hslFilter) Fn(pix pixel) pixel {
	h0 := gm32.Clamp(f.h, -360, 360) / 360
	s0 := gm32.Clamp(f.s, -100, 100) / 100
	l0 := gm32.Clamp(f.l, -100, 100) / 100

	h1, s1, l1 := gcu.RGBToHSL(pix.r, pix.g, pix.b)

	h2 := gm32.Clamp(h1+h0, 0, 1)
	s2 := gm32.Clamp(s1+s0, 0, 1)
	l2 := gm32.Clamp(l1+l0, 0, 1)

	r, g, b := gcu.HSLToRGB(h2, s2, l2)
	return pixel{r, g, b, pix.a}
}

// Changes HSL of each color in the image.
// The hue parameter must be in the range [-360, 360].
// The saturation and lightness parameters must be in the range [-100, 100].
// Each parameter can have any value for merging purposes.
func HSL(h, s, l float32) ColorFilter {
	if h == 0 && s == 0 && l == 0 {
		return nil
	}

	return &hslFilter{
		h:          h,
		s:          s,
		l:          l,
		mergeCount: 1,
	}
}
