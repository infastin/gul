package gft

import "github.com/infastin/gul/gm32"

type ResamplingFilter interface {
	Kernel(x float32) float32
	Support() float32
}

type resampFilter struct {
	kernel  func(x float32) float32
	support float32
}

func (f *resampFilter) Kernel(x float32) float32 {
	return f.kernel(x)
}

func (f *resampFilter) Support() float32 {
	return f.support
}

func MakeResamplingFilter(kernel func(x float32) float32, support float32) ResamplingFilter {
	return &resampFilter{
		kernel:  kernel,
		support: support,
	}
}

var (
	NearestNeighborResampling ResamplingFilter = MakeResamplingFilter(nil, 0)
	BoxResampling             ResamplingFilter = MakeResamplingFilter(boxKernel, 0.5)
	BilinearResampling        ResamplingFilter = MakeResamplingFilter(bicubic5Kernel, 1)
	Bicubic5Resampling        ResamplingFilter = MakeResamplingFilter(bicubic5Kernel, 2)
	Bicubic75Resampling       ResamplingFilter = MakeResamplingFilter(bicubic75Kernel, 2)
	BSplineResampling         ResamplingFilter = MakeResamplingFilter(bSplineKernel, 2)
	MitchellResampling        ResamplingFilter = MakeResamplingFilter(mitchellKernel, 2)
	CatmullRomResampling      ResamplingFilter = MakeResamplingFilter(catmullRomKernel, 2)
	Lanczos3Resampling        ResamplingFilter = MakeResamplingFilter(lanczos3Kernel, 3)
	Lanczos4Resampling        ResamplingFilter = MakeResamplingFilter(lanczos4Kernel, 4)
	Lanczos6Resampling        ResamplingFilter = MakeResamplingFilter(lanczos4Kernel, 6)
	Lanczos12Resampling       ResamplingFilter = MakeResamplingFilter(lanczos4Kernel, 12)
)

func boxKernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 0.5 {
		return 1
	}

	return 0
}

func bilinearKernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 1 {
		return 1 - x
	}

	return 0
}

func bicubic5Kernel(x float32) float32 {
	abs := gm32.Abs(x)

	switch {
	case abs >= 0 && abs <= 1:
		return 1.5*gm32.Pow(abs, 3) - 2.5*gm32.Pow(abs, 2) + 1
	case abs > 1 && abs <= 2:
		return -0.5*gm32.Pow(abs, 3) + 2.5*gm32.Pow(abs, 2) - 4*abs + 2
	default:
		return 0
	}
}

func bicubic75Kernel(x float32) float32 {
	abs := gm32.Abs(x)

	switch {
	case abs >= 0 && abs <= 1:
		return 1.25*gm32.Pow(abs, 3) - 2.25*gm32.Pow(abs, 2) + 1
	case abs > 1 && abs <= 2:
		return -0.75*gm32.Pow(abs, 3) + 3.75*gm32.Pow(abs, 2) - 6*abs + 3
	default:
		return 0
	}
}

func bSplineKernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	switch {
	case x < 1:
		return (0.5 * gm32.Pow(x, 3)) - gm32.Pow(x, 2) + (2.0 / 3.0)
	case x < 2:
		x = 2 - x
		return (1.0 / 6.0) * gm32.Pow(x, 3)
	default:
		return 0
	}
}

func lanczos3Kernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 3 {
		return gm32.Sinc(x) * gm32.Sinc(x/3)
	}

	return 0
}

func lanczos4Kernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 4 {
		return gm32.Sinc(x) * gm32.Sinc(x/4)
	}

	return 0
}

func lanczos6Kernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 6 {
		return gm32.Sinc(x) * gm32.Sinc(x/6)
	}

	return 0
}

func lanczos12Kernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 12 {
		return gm32.Sinc(x) * gm32.Sinc(x/12)
	}

	return 0
}

func mitchell(x, b, c float32) float32 {
	if x < 0 {
		x = -x
	}

	switch {
	case x < 1:
		return ((12-9*b-6*c)*x*x*x + (-18+12*b+6*c)*x*x + (6 - 2*b)) / 6
	case x < 2:
		return ((-b-6*c)*x*x*x + (6*b+30*c)*x*x + (-12*b-48*c)*x + (8*b + 24*c)) / 6
	default:
		return 0
	}
}

func mitchellKernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 2 {
		return mitchell(x, 1.0/3.0, 1.0/3.0)
	}

	return 0
}

func catmullRomKernel(x float32) float32 {
	if x < 0 {
		x = -x
	}

	if x < 2 {
		return mitchell(x, 0, 0.5)
	}

	return 0
}
