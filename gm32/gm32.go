package gm32

import (
	"github.com/infastin/gul/gm64"
	"math"
)

func Sincos(x float32) (float32, float32) {
	sin, cos := math.Sincos(float64(x))
	return float32(sin), float32(cos)
}

func Sin(x float32) float32 {
	sin := math.Sin(float64(x))
	return float32(sin)
}

func Cos(x float32) float32 {
	cos := math.Cos(float64(x))
	return float32(cos)
}

func Tan(x float32) float32 {
	tan := math.Tan(float64(x))
	return float32(tan)
}

func Sinc(x float32) float32 {
	f := gm64.Sinc(float64(x))
	return float32(f)
}

func Abs(x float32) float32 {
	if x < 0 {
		return -x
	}

	return x
}

func Mod(x, y float32) float32 {
	mod := math.Mod(float64(x), float64(y))
	return float32(mod)
}

func Pow(x, y float32) float32 {
	pog := math.Pow(float64(x), float64(y))
	return float32(pog)
}

func Pow10(n int) float32 {
	p := math.Pow10(n)
	return float32(p)
}

func Min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

func Round(x float32) float32 {
	r := math.Round(float64(x))
	return float32(r)
}

func RoundN(x float32, n int) float32 {
	return Round(x*Pow10(n)) / Pow10(n)
}

func Floor(x float32) float32 {
	f := math.Floor(float64(x))
	return float32(f)
}

func Ceil(x float32) float32 {
	f := math.Ceil(float64(x))
	return float32(f)
}

func Clamp(val, min, max float32) float32 {
	return Min(Max(min, val), max)
}
