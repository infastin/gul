package gft

import (
	"math"

	"github.com/infastin/gul/gm32"
)

// This is a filter used for combining by using CombineColorhanFilters.
// Must be a pointer.
type ColorchanFilter interface {
	// Returns changed color channel.
	Fn(x float32) float32

	// Returns true, if it is possible to create a lookup table usign Fn.
	// Otherwise, returns false.
	UseLut() bool
}

// This colorhan filter can merge other colorhan filters into itself.
type MergingColorchanFilter interface {
	ColorchanFilter

	// Prepares the filter before calling Fn multiple times.
	Prepare()

	// Returns true, if it is possible to merge one filter into an instance of interface.
	// Otherwise, returns false.
	CanMerge(filter ColorchanFilter) bool

	// Returns true, if it is possible to demerge one filter from an instance of interface.
	// Otherwise, returns false.
	CanUndo(filter ColorchanFilter) bool

	// Merges one filter into an instance of interface.
	Merge(filter ColorchanFilter)

	// Demerges one filter from an instance of interface.
	Undo(filter ColorchanFilter)

	// Returns true, if nothing will change after applying the filter.
	// Otherwise, returns false.
	Skip() bool

	// Returns a copy of the filter.
	Copy() ColorchanFilter
}

type colorchanFilterFunc struct {
	fn     func(x float32) float32
	useLut bool
}

func (f *colorchanFilterFunc) UseLut() bool {
	return f.useLut
}

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
func Invert() MergingColorchanFilter {
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
func Gamma(gamma float32) MergingColorchanFilter {
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
func Contrast(perc float32) MergingColorchanFilter {
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
func Brightness(perc float32) MergingColorchanFilter {
	if perc == 0 {
		return nil
	}

	return &brightnessFilter{
		brightness: perc,
	}
}
