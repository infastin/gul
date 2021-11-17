package gft

import (
	"image"
	"image/draw"
	"math"

	"github.com/infastin/gul/gm32"
	"github.com/infastin/gul/tools"
)

// This is a filter used for combining by using CombineColorhanFilters.
// Must be a pointer.
type ColorchanFilter interface {
	// Returns changed color channel.
	Fn(x float32) float32

	// Prepares the filter before calling Fn multiple times.
	Prepare()

	// Returns true, if it is possible to create a lookup table usign Fn.
	// Otherwise, returns false.
	UseLut() bool

	// Returns true, if it is possible to combine two filters.
	// Otherwise, returns false.
	CanMerge(filter ColorchanFilter) bool

	// Returns true, if it is possible to decombine two filters.
	// Otherwise, returns false.
	CanUndo(filter ColorchanFilter) bool

	// Combines two filters, if CanMerge has returned true.
	Merge(filter ColorchanFilter)

	// Decombines two filters, if CanUndo has returned true.
	Undo(filter ColorchanFilter)

	// Returns true, if nothing will change after applying the filter.
	// Otherwise, returns false.
	Skip() bool

	// Returns a copy of the filter.
	Copy() ColorchanFilter
}

type combineColorchanFilter struct {
	filters    []ColorchanFilter
	luts       [][]float32
	mergeCount uint
}

func (f *combineColorchanFilter) Merge(filter Filter) bool {
	filt, ok := filter.(*combineColorchanFilter)
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

func (f *combineColorchanFilter) Undo(filter Filter) bool {
	filt, ok := filter.(*combineColorchanFilter)
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

func (f *combineColorchanFilter) Skip() bool {
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

func (f *combineColorchanFilter) Copy() Filter {
	r := &combineColorchanFilter{
		mergeCount: f.mergeCount,
	}

	r.filters = make([]ColorchanFilter, len(f.filters))
	for i := 0; i < len(f.filters); i++ {
		if f.filters[i] == nil {
			r.filters[i] = nil
			continue
		}

		r.filters[i] = f.filters[i].Copy()
	}

	return r
}

func (f *combineColorchanFilter) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (f *combineColorchanFilter) makeLut(lutSize int, index int) {
	lutLen := len(f.luts[index])
	start := 0

	if lutLen == 0 {
		f.luts[index] = make([]float32, lutSize)
	} else if lutLen < lutSize {
		newLut := make([]float32, lutSize)
		copy(newLut[:lutLen], f.luts[index])
		f.luts[index] = newLut
		start = lutLen
	}

	q := float32(1) / float32(lutSize-1)
	for i := start; i < lutSize; i++ {
		v := float32(i) * q
		f.luts[index][i] = f.filters[index].Fn(v)
	}
}

func (f *combineColorchanFilter) getFromLut(x float32, index int) float32 {
	i := int(gm32.Round(x * float32(len(f.luts[index])-1)))
	return f.luts[index][i]
}

func (f *combineColorchanFilter) Apply(dst draw.Image, src image.Image, parallel bool) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	useLut := make([]bool, len(f.filters))

	for i, filt := range f.filters {
		if filt == nil {
			continue
		}

		filt.Prepare()

		if filt.UseLut() {
			lutSize := len(f.luts[i])
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
					f.makeLut(neededLutSize, i)
				}

				useLut[i] = true
			}
		} else {
			useLut[i] = false
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

				for i, filt := range f.filters {
					if filt == nil {
						continue
					}

					if useLut[i] {
						pix.r = f.getFromLut(pix.r, i)
						pix.g = f.getFromLut(pix.g, i)
						pix.b = f.getFromLut(pix.b, i)
					} else {
						pix.r = filt.Fn(pix.r)
						pix.g = filt.Fn(pix.g)
						pix.b = filt.Fn(pix.b)
					}
				}

				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pix)
			}
		}
	})
}

// Creates combination of colorchan filters and returns filter.
func CombineColorchanFilters(filters ...ColorchanFilter) Filter {
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

	return &combineColorchanFilter{
		filters:    filters,
		luts:       make([][]float32, len(filters)),
		mergeCount: 1,
	}
}

type colorchanFilterFunc struct {
	fn     func(x float32) float32
	useLut bool
}

func (f *colorchanFilterFunc) CanMerge(ColorchanFilter) bool {
	return false
}

func (f *colorchanFilterFunc) Merge(ColorchanFilter) {}

func (f *colorchanFilterFunc) CanUndo(ColorchanFilter) bool {
	return false
}

func (f *colorchanFilterFunc) Undo(ColorchanFilter) {}

func (f *colorchanFilterFunc) Skip() bool {
	return false
}

func (f *colorchanFilterFunc) Copy() ColorchanFilter {
	return &colorchanFilterFunc{
		fn:     f.fn,
		useLut: f.useLut,
	}
}

func (f *colorchanFilterFunc) UseLut() bool {
	return f.useLut
}

func (f *colorchanFilterFunc) Prepare() {}

func (f *colorchanFilterFunc) Fn(x float32) float32 {
	return f.fn(x)
}

func ColorchanFilterFunc(fn func(x float32) float32, useLut bool) ColorchanFilter {
	return &colorchanFilterFunc{
		fn:     fn,
		useLut: useLut,
	}
}

type invertFilter struct {
	state byte
}

func (f *invertFilter) CanMerge(filter ColorchanFilter) bool {
	if _, ok := filter.(*invertFilter); ok {
		return true
	}

	return false
}

func (f *invertFilter) Merge(filter ColorchanFilter) {
	filt := filter.(*invertFilter)
	f.state ^= filt.state
}

func (f *invertFilter) Undo(filter ColorchanFilter) {
	filt := filter.(*invertFilter)
	f.state = f.state ^ filt.state
}

func (f *invertFilter) CanUndo(filter ColorchanFilter) bool {
	if _, ok := filter.(*invertFilter); ok {
		return true
	}

	return false
}

func (f *invertFilter) Skip() bool {
	return f.state == 0
}

func (f *invertFilter) UseLut() bool {
	return true
}

func (f *invertFilter) Copy() ColorchanFilter {
	return &invertFilter{
		state: f.state,
	}
}

func (f *invertFilter) Prepare() {}

func (f *invertFilter) Fn(x float32) float32 {
	return 1 - x
}

// Negates the colors of an image.
func Invert() ColorchanFilter {
	return &invertFilter{
		state: 1,
	}
}

type gammaFilter struct {
	gamma float32
}

func (f *gammaFilter) CanMerge(filter ColorchanFilter) bool {
	if _, ok := filter.(*gammaFilter); ok {
		return true
	}

	return false
}

func (f *gammaFilter) Merge(filter ColorchanFilter) {
	filt := filter.(*gammaFilter)
	f.gamma = gm32.Max(1.0e-5, f.gamma+filt.gamma)
}

func (f *gammaFilter) CanUndo(filter ColorchanFilter) bool {
	if _, ok := filter.(*gammaFilter); ok {
		return true
	}

	return false
}

func (f *gammaFilter) Undo(filter ColorchanFilter) {
	filt := filter.(*gammaFilter)
	f.gamma = gm32.Max(1.0e-5, f.gamma-filt.gamma)
}

func (f *gammaFilter) Skip() bool {
	return f.gamma == 1
}

func (f *gammaFilter) Copy() ColorchanFilter {
	return &gammaFilter{
		gamma: f.gamma,
	}
}

func (f *gammaFilter) Prepare() {
	f.gamma = gm32.Max(1.0e-5, f.gamma)
}

func (f *gammaFilter) Fn(x float32) float32 {
	e := 1 / f.gamma
	return gm32.Pow(x, e)
}

func (f *gammaFilter) UseLut() bool {
	return true
}

// Gamma creates a filter that performs a gamma correction on an image.
// The gamma parameter must be positive. Gamma = 1 gives the original image.
// Gamma less than 1 darkens the image and gamma greater than 1 lightens it.
func Gamma(gamma float32) ColorchanFilter {
	if gamma == 1 {
		return nil
	}

	return &gammaFilter{
		gamma: gamma,
	}
}

type contrastFilter struct {
	contrast float32
}

func (f *contrastFilter) CanMerge(filter ColorchanFilter) bool {
	if _, ok := filter.(*contrastFilter); ok {
		return true
	}

	return false
}

func (f *contrastFilter) Merge(filter ColorchanFilter) {
	filt := filter.(*contrastFilter)
	f.contrast = gm32.Clamp(f.contrast+filt.contrast, -100, 100)
}

func (f *contrastFilter) CanUndo(filter ColorchanFilter) bool {
	if _, ok := filter.(*contrastFilter); ok {
		return true
	}

	return false
}

func (f *contrastFilter) Undo(filter ColorchanFilter) {
	filt := filter.(*contrastFilter)
	f.contrast = gm32.Clamp(f.contrast-filt.contrast, -100, 100)
}

func (f *contrastFilter) Skip() bool {
	return f.contrast == 0
}

func (f *contrastFilter) Copy() ColorchanFilter {
	return &contrastFilter{
		contrast: f.contrast,
	}
}

func (f *contrastFilter) Prepare() {
	f.contrast = gm32.Clamp(f.contrast, -100, 100)
}

func (f *contrastFilter) Fn(x float32) float32 {
	alpha := (f.contrast / 100) + 1
	alpha = gm32.Tan(alpha * (math.Pi / 4))

	c := (x-0.5)*alpha + 0.5
	return gm32.Clamp(c, 0, 1)
}

func (f *contrastFilter) UseLut() bool {
	return false
}

// Changes contrast of an image.
// The percentage parameter must be in the range [-100, 100].
// It can have any value for merging purposes.
// The percentage = -100 gives solid gray image. The percentage = 100 gives an overcontrasted image.
func Contrast(perc float32) ColorchanFilter {
	if perc == 0 {
		return nil
	}

	return &contrastFilter{
		contrast: perc,
	}
}

type brightnessFilter struct {
	brightness float32
}

func (f *brightnessFilter) CanMerge(filter ColorchanFilter) bool {
	if _, ok := filter.(*brightnessFilter); ok {
		return true
	}

	return false
}

func (f *brightnessFilter) Merge(filter ColorchanFilter) {
	filt := filter.(*brightnessFilter)
	f.brightness = gm32.Clamp(f.brightness+filt.brightness, -100, 100)
}

func (f *brightnessFilter) CanUndo(filter ColorchanFilter) bool {
	if _, ok := filter.(*brightnessFilter); ok {
		return true
	}

	return false
}

func (f *brightnessFilter) Undo(filter ColorchanFilter) {
	filt := filter.(*brightnessFilter)
	f.brightness = gm32.Clamp(f.brightness-filt.brightness, -100, 100)
}

func (f *brightnessFilter) Skip() bool {
	return f.brightness == 0
}

func (f *brightnessFilter) Copy() ColorchanFilter {
	return &brightnessFilter{
		brightness: f.brightness,
	}
}

func (f *brightnessFilter) Prepare() {
	f.brightness = gm32.Clamp(f.brightness, -100, 100)
}

func (f *brightnessFilter) Fn(x float32) float32 {
	beta := f.brightness / 100

	if beta < 0 {
		x *= (1 + beta)
	} else {
		x += (1 - x) * beta
	}

	return gm32.Clamp(x, 0, 1)
}

func (f *brightnessFilter) UseLut() bool {
	return false
}

// Changes brightness of an image.
// The percentage parameter must be in the range [-100, 100].
// It can have any value for merging purposes.
// The percentage = -100 gives solid black image. The percentage = 100 gives solid white image.
func Brightness(perc float32) ColorchanFilter {
	if perc == 0 {
		return nil
	}

	return &brightnessFilter{
		brightness: perc,
	}
}
