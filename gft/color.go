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
	Prepare()
	Merge(filter ColorFilter)
	Undo(filter ColorFilter)
	Skip() bool
	Copy() ColorFilter
	CanMerge(filter ColorFilter) bool
	CanUndo(filter ColorFilter) bool
}

type combineColorFilter struct {
	filters    []ColorFilter
	mergeCount uint
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

	f.mergeCount++

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

	for i := 0; i < len(f.filters); i++ {
		if filt.filters[i] == nil {
			continue
		}

		f.filters[i].Undo(filt.filters[i])
	}

	f.mergeCount--

	return f.mergeCount == 0
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
	r := &combineColorFilter{
		mergeCount: f.mergeCount,
	}

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

	for _, filt := range f.filters {
		if filt == nil {
			continue
		}

		filt.Prepare()
	}

	procs := 1
	if parallel {
		procs = 0
	}

	tools.Parallelize(procs, srcb.Min.Y, srcb.Max.Y, 1, func(start, end int) {
		for y := start; y < end; y++ {
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				pix := pixGetter.getPixel(x, y)

				for _, filt := range f.filters {
					if filt == nil {
						continue
					}

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
		filters:    filters,
		mergeCount: 1,
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

func (f *colorFilterFunc) Undo(ColorFilter) {}

func (f *colorFilterFunc) Skip() bool {
	return false
}

func (f *colorFilterFunc) Copy() ColorFilter {
	return &colorFilterFunc{
		fn: f.fn,
	}
}

func (f *colorFilterFunc) Prepare() {}

func (f *colorFilterFunc) Fn(pix pixel) pixel {
	return f.fn(pix)
}

func ColorFilterFunc(fn func(pix pixel) pixel) ColorFilter {
	return &colorFilterFunc{
		fn: fn,
	}
}

type sepiaFilter struct {
	percentage float32
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
}

func (f *sepiaFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*sepiaFilter); ok {
		return true
	}

	return false
}

func (f *sepiaFilter) Undo(filter ColorFilter) {
	filt := filter.(*sepiaFilter)
	f.percentage = gm32.Clamp(f.percentage-filt.percentage, 0, 100)
}

func (f *sepiaFilter) Skip() bool {
	return f.percentage == 0
}

func (f *sepiaFilter) Prepare() {
	f.percentage = gm32.Clamp(f.percentage, 0, 100)
}

func (f *sepiaFilter) Fn(pix pixel) pixel {
	rat := f.percentage / 100

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
	}
}

type hsbFilter struct {
	h, s, b float32
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
}

func (f *hsbFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*hsbFilter); ok {
		return true
	}

	return false
}

func (f *hsbFilter) Undo(filter ColorFilter) {
	filt := filter.(*hsbFilter)

	f.h = gm32.Clamp(f.h-filt.h, -360, 360)
	f.s = gm32.Clamp(f.s-filt.s, -100, 100)
	f.b = gm32.Clamp(f.b-filt.b, -100, 100)
}

func (f *hsbFilter) Skip() bool {
	return f.h == 0 && f.s == 0 && f.b == 0
}

func (f *hsbFilter) Copy() ColorFilter {
	return &hsbFilter{
		h: f.h,
		s: f.s,
		b: f.b,
	}
}

func (f *hsbFilter) Prepare() {
	f.h = gm32.Clamp(f.h, -360, 360)
	f.s = gm32.Clamp(f.s, -100, 100)
	f.b = gm32.Clamp(f.b, -100, 100)
}

func (f *hsbFilter) Fn(pix pixel) pixel {
	h0 := f.h / 360
	s0 := f.s / 100
	b0 := f.b / 100

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
		h: h,
		s: s,
		b: b,
	}
}

type hslFilter struct {
	h, s, l float32
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
}

func (f *hslFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*hslFilter); ok {
		return true
	}

	return false
}

func (f *hslFilter) Undo(filter ColorFilter) {
	filt := filter.(*hslFilter)

	f.h = gm32.Clamp(f.h-filt.h, -360, 360)
	f.s = gm32.Clamp(f.s-filt.s, -100, 100)
	f.l = gm32.Clamp(f.l-filt.l, -100, 100)
}

func (f *hslFilter) Skip() bool {
	return f.h == 0 && f.s == 0 && f.l == 0
}

func (f *hslFilter) Copy() ColorFilter {
	return &hslFilter{
		h: f.h,
		s: f.s,
		l: f.l,
	}
}

func (f *hslFilter) Prepare() {
	f.h = gm32.Clamp(f.h, -360, 360)
	f.s = gm32.Clamp(f.s, -100, 100)
	f.l = gm32.Clamp(f.l, -100, 100)
}

func (f *hslFilter) Fn(pix pixel) pixel {
	h0 := f.h / 360
	s0 := f.s / 100
	l0 := f.l / 100

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
		h: h,
		s: s,
		l: l,
	}
}

type ColorLevels struct {
	CyanRed      float32
	MagentaGreen float32
	YellowBlue   float32
}

func (cl1 ColorLevels) Add(cl2 ColorLevels) ColorLevels {
	return ColorLevels{
		CyanRed:      gm32.Clamp(cl1.CyanRed+cl2.CyanRed, -100, 100),
		MagentaGreen: gm32.Clamp(cl1.MagentaGreen+cl2.MagentaGreen, -100, 100),
		YellowBlue:   gm32.Clamp(cl1.YellowBlue+cl2.YellowBlue, -100, 100),
	}
}

func (cl1 ColorLevels) Sub(cl2 ColorLevels) ColorLevels {
	return ColorLevels{
		CyanRed:      gm32.Clamp(cl1.CyanRed-cl2.CyanRed, -100, 100),
		MagentaGreen: gm32.Clamp(cl1.MagentaGreen-cl2.MagentaGreen, -100, 100),
		YellowBlue:   gm32.Clamp(cl1.YellowBlue-cl2.YellowBlue, -100, 100),
	}
}

func (cl ColorLevels) Clamp() ColorLevels {
	return ColorLevels{
		CyanRed:      gm32.Clamp(cl.CyanRed, -100, 100),
		MagentaGreen: gm32.Clamp(cl.MagentaGreen, -100, 100),
		YellowBlue:   gm32.Clamp(cl.YellowBlue, -100, 100),
	}
}

func (cl ColorLevels) IsZero() bool {
	return cl.CyanRed == 0 && cl.MagentaGreen == 0 && cl.YellowBlue == 0
}

type colorBalanceFilter struct {
	shadows            ColorLevels
	midtones           ColorLevels
	highlights         ColorLevels
	preserveLuminosity bool
}

func (f *colorBalanceFilter) CanMerge(filter ColorFilter) bool {
	if _, ok := filter.(*colorBalanceFilter); ok {
		return true
	}

	return false
}

func (f *colorBalanceFilter) Merge(filter ColorFilter) {
	filt := filter.(*colorBalanceFilter)

	f.shadows = f.shadows.Add(filt.shadows)
	f.midtones = f.midtones.Add(filt.midtones)
	f.highlights = f.highlights.Add(filt.highlights)
}

func (f *colorBalanceFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*colorBalanceFilter); ok {
		return true
	}

	return false
}

func (f *colorBalanceFilter) Undo(filter ColorFilter) {
	filt := filter.(*colorBalanceFilter)

	f.shadows = f.shadows.Sub(filt.shadows)
	f.midtones = f.midtones.Sub(filt.midtones)
	f.highlights = f.highlights.Sub(filt.highlights)
}

func (f *colorBalanceFilter) Skip() bool {
	return f.shadows.IsZero() && f.midtones.IsZero() && f.highlights.IsZero()
}

func (f *colorBalanceFilter) Copy() ColorFilter {
	return &colorBalanceFilter{
		shadows:            f.shadows,
		midtones:           f.midtones,
		highlights:         f.highlights,
		preserveLuminosity: f.preserveLuminosity,
	}
}

func (f *colorBalanceFilter) mask(val, lightness, shadows, midtones, highlights float32) float32 {
	const (
		a     = 0.25
		b     = 0.333
		scale = 0.7
	)

	shadows /= 100
	midtones /= 100
	highlights /= 100

	shadows *= gm32.Clamp(((lightness-b)/-a)+0.5, 0, 1) * scale
	midtones *= gm32.Clamp(((lightness-b)/a)+0.5, 0, 1) *
		gm32.Clamp(((lightness+b-1)/-a)+0.5, 0, 1) * scale
	highlights *= gm32.Clamp(((lightness+b-1)/a)+0.5, 0, 1) * scale

	val += shadows
	val += midtones
	val += highlights
	val = gm32.Clamp(val, 0, 1)

	return val
}

func (f *colorBalanceFilter) Prepare() {
	f.shadows = f.shadows.Clamp()
	f.midtones = f.midtones.Clamp()
	f.highlights = f.highlights.Clamp()
}

func (f *colorBalanceFilter) Fn(pix pixel) pixel {
	_, _, l := gcu.RGBToHSL(pix.r, pix.g, pix.b)

	rn := f.mask(pix.r, l, f.shadows.CyanRed, f.midtones.CyanRed, f.highlights.CyanRed)
	gn := f.mask(pix.g, l, f.shadows.MagentaGreen, f.midtones.MagentaGreen, f.highlights.MagentaGreen)
	bn := f.mask(pix.b, l, f.shadows.YellowBlue, f.midtones.YellowBlue, f.highlights.YellowBlue)

	if f.preserveLuminosity {
		h2, s2, _ := gcu.RGBToHSL(rn, gn, bn)
		rn, gn, bn = gcu.HSLToRGB(h2, s2, l)
	}

	return pixel{rn, gn, bn, pix.a}
}

// Adjusts color distribution in an image.
// Each color level  must be in the range [-100, 100].
// The color levels can have any value for merging purposes.
func ColorBalance(shadows, midtones, highlights ColorLevels, preserveLuminosity bool) ColorFilter {
	if shadows.IsZero() && midtones.IsZero() && highlights.IsZero() {
		return nil
	}

	return &colorBalanceFilter{
		shadows:            shadows,
		midtones:           midtones,
		highlights:         highlights,
		preserveLuminosity: preserveLuminosity,
	}
}

type colorizeFilter struct {
	h, s, l float32
}

func (f *colorizeFilter) CanMerge(filter ColorFilter) bool {
	if _, ok := filter.(*colorizeFilter); ok {
		return true
	}

	return false
}

func (f *colorizeFilter) Merge(filter ColorFilter) {
	filt := filter.(*colorizeFilter)

	f.h = gm32.Clamp(f.h+filt.h, 0, 360)
	f.s = gm32.Clamp(f.s+filt.s, 0, 100)
	f.l = gm32.Clamp(f.l+filt.l, -100, 100)
}

func (f *colorizeFilter) CanUndo(filter ColorFilter) bool {
	if _, ok := filter.(*colorizeFilter); ok {
		return true
	}

	return false
}

func (f *colorizeFilter) Undo(filter ColorFilter) {
	filt := filter.(*colorizeFilter)

	f.h = gm32.Clamp(f.h-filt.h, 0, 360)
	f.s = gm32.Clamp(f.s-filt.s, 0, 100)
	f.l = gm32.Clamp(f.l-filt.l, -100, 100)
}

func (f *colorizeFilter) Skip() bool {
	return false
}

func (f *colorizeFilter) Copy() ColorFilter {
	return &colorizeFilter{
		h: f.h,
		s: f.s,
		l: f.l,
	}
}

func (f *colorizeFilter) Prepare() {
	f.h = gm32.Clamp(f.h, 0, 360)
	f.s = gm32.Clamp(f.s, 0, 100)
	f.l = gm32.Clamp(f.l, -100, 100)
}

// https://github.com/GNOME/gimp/blob/708f075f804caa5cbd11cae2f85f3726456aedcb/app/operations/gimpoperationcolorize.c#L222
func (f *colorizeFilter) Fn(pix pixel) pixel {
	h := f.h / 360
	s := f.s / 100
	l := f.l / 100

	lum := gcu.RGBLuminance(pix.r, pix.g, pix.b)

	switch {
	case l > 0:
		lum = lum * (1 - l)
		lum += l
	case l < 0:
		lum = lum * (l + 1)
	}

	r, g, b := gcu.HSLToRGB(h, s, lum)
	return pixel{r, g, b, pix.a}
}

// Colorizes an image.
// The hue parameter must be in the range [0, 360].
// The saturation parameter must be in the range [0, 100].
// The lightness parameter must be in the range [-100, 100].
func Colorize(h, s, l float32) ColorFilter {
	return &colorizeFilter{
		h: h,
		s: s,
		l: l,
	}
}

var (
	// Grayscales an image.
	Grayscale ColorFilter = ColorFilterFunc(grayscale)
)

func grayscale(pix pixel) pixel {
	v := 0.299*pix.r + 0.587*pix.g + 0.114*pix.b
	return pixel{v, v, v, pix.a}
}
