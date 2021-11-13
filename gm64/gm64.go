package gm64

import (
	"math"
)

func Sinc(x float64) float64 {
	if x == 0 {
		return 1
	}

	return math.Sin(x*math.Pi) / (x * math.Pi)
}

func Clamp(val, min, max float64) float64 {
	return math.Min(math.Max(min, val), max)
}

func RoundN(x float64, n int) float64 {
	return math.Round(x*math.Pow10(n)) / math.Pow10(n)
}
