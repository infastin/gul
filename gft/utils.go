package gft

import (
	"github.com/infastin/gul/gm32"
)

const (
	q12 = 0.5
	q13 = 1.0 / 3.0
	q23 = 2.0 / 3.0
	q16 = 1.0 / 6.0
)

func RGBToHSL(r, g, b float32) (h, s, l float32) {
	min := gm32.Min(r, gm32.Min(g, b))
	max := gm32.Max(r, gm32.Max(g, b))

	l = (min + max) * q12
	if min == max {
		return
	}

	d := max - min
	switch max {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	case b:
		h = (r-g)/d + 4
	}

	h *= q16

	switch {
	case l == 0:
		s = 0
	case l <= q12:
		s = d / (2 * l)
	case l < 1:
		s = d / (2 - 2*l)
	}

	return
}

func HSLToRGB(h, s, l float32) (r, g, b float32) {
	var q float32
	if l < q12 {
		q = l * (1 + s)
	} else {
		q = (l + s) - (l * s)
	}

	p := 2*l - q

	t := make([]float32, 3)
	t[0] = h + q13
	t[1] = h
	t[2] = h - q13

	for i := 0; i < 3; i++ {
		switch {
		case t[i] < 0:
			t[i] += 1
		case t[i] > 1:
			t[i] -= 1
		}
	}

	for i := 0; i < 3; i++ {
		switch {
		case t[i] < q16:
			t[i] = p + ((q - p) * 6 * t[i])
		case t[i] < q12:
			t[i] = q
		case t[i] < q23:
			t[i] = p + ((q - p) * (q23 - t[i]) * 6)
		default:
			t[i] = p
		}
	}

	r, g, b = t[0], t[1], t[2]
	return
}

func f32u8(val float32) uint8 {
	fv := gm32.Clamp(val, 0, 0xff)
	return uint8(gm32.Round(fv))
}

func f32u16(val float32) uint16 {
	fv := gm32.Clamp(val, 0, 0xffff)
	return uint16(gm32.Round(fv))
}
