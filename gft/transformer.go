package gft

import (
	"image"
)

type Transformer interface {
	// Returns bounds of the image after transformation.
	Bounds(src image.Rectangle) image.Rectangle

	// Returns the new point after transformation.
	// Last two booleans stand for opposite direction of the x and y coordinates
	// Example: if x = 10 and image width = 100, then
	// setting oppX to true will mean that x = 89.
	Transform(sx, sy int) (dx, dy int, oppX, oppY bool)

	// If possible, merges two transformers
	// and returns a new transformer with true as the second value.
	// May return nil transformer if got nothing after merge.
	// If merging isn't possible, returns false as the second value.
	// And the first returned value is ignored.
	Merge(in Transformer) (Transformer, bool)

	// If out := t.Merge(in) then
	// in := t.Recreate(out)
	Recreate(in Transformer) (Transformer, bool)
}

var (
	FlipHTransformer      Transformer = &fliphTransformer{}
	FlipVTransformer      Transformer = &flipvTransformer{}
	TransposeTransformer  Transformer = &transposeTransformer{}
	TransverseTransformer Transformer = &transverseTransformer{}
	Rotate90Transformer   Transformer = &rotate90Transformer{}
	Rotate180Transformer  Transformer = &rotate180Transformer{}
	Rotate270Transformer  Transformer = &rotate270Transformer{}
)

type fliphTransformer struct{}

func (t *fliphTransformer) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (t *fliphTransformer) Transform(x, y int) (int, int, bool, bool) {
	return x, y, true, false
}

func (t *fliphTransformer) Merge(in Transformer) (Transformer, bool) {
	if _, ok := in.(*fliphTransformer); ok {
		return nil, true
	}

	return nil, false
}

func (t *fliphTransformer) Recreate(in Transformer) (Transformer, bool) {
	if in == nil {
		return FlipHTransformer, true
	}

	if _, ok := in.(*fliphTransformer); ok {
		return nil, true
	}

	return nil, false
}

type flipvTransformer struct{}

func (t *flipvTransformer) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (t *flipvTransformer) Transform(x, y int) (int, int, bool, bool) {
	return x, y, false, true
}

func (t *flipvTransformer) Merge(in Transformer) (Transformer, bool) {
	if _, ok := in.(*flipvTransformer); ok {
		return nil, true
	}

	return nil, false
}

func (t *flipvTransformer) Recreate(in Transformer) (Transformer, bool) {
	if in == nil {
		return FlipVTransformer, true
	}

	if _, ok := in.(*flipvTransformer); ok {
		return nil, true
	}

	return nil, false
}

type transposeTransformer struct{}

func (t *transposeTransformer) Bounds(src image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, src.Dy(), src.Dx())
}

func (t *transposeTransformer) Transform(x, y int) (int, int, bool, bool) {
	return y, x, false, false
}

func (t *transposeTransformer) Merge(in Transformer) (Transformer, bool) {
	if _, ok := in.(*transposeTransformer); ok {
		return nil, true
	}

	return nil, false
}

func (t *transposeTransformer) Recreate(in Transformer) (Transformer, bool) {
	if in == nil {
		return TransposeTransformer, true
	}

	if _, ok := in.(*transposeTransformer); ok {
		return nil, true
	}

	return nil, false
}

type transverseTransformer struct{}

func (t *transverseTransformer) Bounds(src image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, src.Dy(), src.Dx())
}

func (t *transverseTransformer) Transform(x, y int) (int, int, bool, bool) {
	return y, x, true, true
}

func (t *transverseTransformer) Merge(in Transformer) (Transformer, bool) {
	if _, ok := in.(*transverseTransformer); ok {
		return nil, true
	}

	return nil, false
}

func (t *transverseTransformer) Recreate(in Transformer) (Transformer, bool) {
	if in == nil {
		return TransverseTransformer, true
	}

	if _, ok := in.(*transverseTransformer); ok {
		return nil, true
	}

	return nil, false
}

type rotate90Transformer struct{}

func (t *rotate90Transformer) Bounds(src image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, src.Dy(), src.Dx())
}

func (t *rotate90Transformer) Transform(x, y int) (int, int, bool, bool) {
	return y, x, true, false
}

func (t *rotate90Transformer) Merge(in Transformer) (Transformer, bool) {
	switch in.(type) {
	case *rotate90Transformer:
		return Rotate180Transformer, true
	case *rotate180Transformer:
		return Rotate270Transformer, true
	case *rotate270Transformer:
		return nil, true
	default:
		return nil, false
	}
}

func (t *rotate90Transformer) Recreate(in Transformer) (Transformer, bool) {
	switch in.(type) {
	case nil:
		return Rotate270Transformer, true
	case *rotate90Transformer:
		return nil, true
	case *rotate180Transformer:
		return Rotate90Transformer, true
	case *rotate270Transformer:
		return Rotate180Transformer, true
	default:
		return nil, false
	}
}

type rotate180Transformer struct{}

func (t *rotate180Transformer) Bounds(src image.Rectangle) image.Rectangle {
	return src
}

func (t *rotate180Transformer) Transform(x, y int) (int, int, bool, bool) {
	return x, y, true, true
}

func (t *rotate180Transformer) Merge(in Transformer) (Transformer, bool) {
	switch in.(type) {
	case *rotate90Transformer:
		return Rotate270Transformer, true
	case *rotate180Transformer:
		return nil, true
	case *rotate270Transformer:
		return Rotate90Transformer, true
	default:
		return nil, false
	}
}

func (t *rotate180Transformer) Recreate(in Transformer) (Transformer, bool) {
	switch in.(type) {
	case nil:
		return Rotate180Transformer, true
	case *rotate90Transformer:
		return Rotate270Transformer, true
	case *rotate180Transformer:
		return nil, true
	case *rotate270Transformer:
		return Rotate90Transformer, true
	default:
		return nil, false
	}
}

type rotate270Transformer struct{}

func (t *rotate270Transformer) Bounds(src image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, src.Dy(), src.Dx())
}

func (t *rotate270Transformer) Transform(x, y int) (int, int, bool, bool) {
	return y, x, false, true
}

func (t *rotate270Transformer) Merge(in Transformer) (Transformer, bool) {
	switch in.(type) {
	case *rotate90Transformer:
		return nil, true
	case *rotate180Transformer:
		return Rotate90Transformer, true
	case *rotate270Transformer:
		return Rotate180Transformer, true
	default:
		return nil, false
	}
}

func (t *rotate270Transformer) Recreate(in Transformer) (Transformer, bool) {
	switch in.(type) {
	case nil:
		return Rotate90Transformer, true
	case *rotate90Transformer:
		return Rotate180Transformer, true
	case *rotate180Transformer:
		return Rotate270Transformer, true
	case *rotate270Transformer:
		return nil, true
	default:
		return nil, false
	}
}
