package gm32

import "github.com/infastin/gul/pkg/matrix"

func InterpolateBilinear(v00, v01, v10, v11, fx, fy float32) float32 {
	tmp := (1-fy)*((1-fx)*v00+fx*v01) +
		fy*((1-fx)*v10+fx*v11)

	return tmp
}

func bicubicKernel(x, a float32) float32 {
	abs := Abs(x)

	switch {
	case abs >= 0 && abs <= 1:
		return (a+2)*Pow(abs, 3) - (a+3)*Pow(abs, 2) + 1
	case abs > 1 && abs <= 2:
		return a*Pow(abs, 3) - (5*a)*Pow(abs, 2) + (8*a)*abs - 4*a
	default:
		return 0
	}
}

func InterpolateBicubic(left, mid, right *matrix.Matrix32, a float32) float32 {
	for i := 0; i < left.N; i++ {
		left.Data[i] = bicubicKernel(left.Data[i], a)
	}

	for i := 0; i < right.M; i++ {
		right.Data[i] = bicubicKernel(right.Data[i], a)
	}

	return left.Mul(mid).Mul(right).Data[0]
}
