package gcu

import (
	"image/color"

	"github.com/infastin/gul/gm32"
)

const (
	qf16 = 1.0 / 0xffff
	q12  = 0.5
	q13  = 1.0 / 3.0
	q23  = 2.0 / 3.0
	q16  = 1.0 / 6.0
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

func RGBToHSV(r, g, b float32) (h, s, v float32) {
	min := gm32.Min(r, gm32.Min(g, b))
	max := gm32.Max(r, gm32.Max(g, b))

	if max == 0 {
		s = 0
	} else {
		s = 1 - min/max
	}

	v = max

	if min == max {
		return
	}

	d := max - min

	switch {
	case max == r && g >= b:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case max == g:
		h = (b-r)/d + 2
	case max == b:
		h = (r-g)/d + 4
	}

	h *= q16

	return
}

func HSVToRGB(h, s, v float32) (r, g, b float32) {
	ht := h * 6
	c := s * v
	x := c * (1 - gm32.Abs(gm32.Mod(ht, 2)-1))
	m := v - c

	var r1, g1, b1 float32

	switch {
	case ht < 0:
		r1, g1, b1 = 0, 0, 0
	case ht < 1:
		r1, g1, b1 = c, x, 0
	case ht < 2:
		r1, g1, b1 = x, c, 0
	case ht < 3:
		r1, g1, b1 = 0, c, x
	case ht < 4:
		r1, g1, b1 = 0, x, c
	case ht < 5:
		r1, g1, b1 = x, 0, c
	case ht < 6:
		r1, g1, b1 = c, 0, x
	}

	r = r1 + m
	b = b1 + m
	g = g1 + m

	return
}

type HSLA struct {
	H, S, L, A float32
}

func (c HSLA) RGBA() (r, g, b, a uint32) {
	fr, fg, fb := HSLToRGB(c.H, c.S, c.L)
	fa := c.A * 0xffff

	fr = fr * fa
	fg = fg * fa
	fb = fb * fa

	r = uint32(gm32.Round(gm32.Clamp(fr, 0, 0xffff)))
	g = uint32(gm32.Round(gm32.Clamp(fg, 0, 0xffff)))
	b = uint32(gm32.Round(gm32.Clamp(fb, 0, 0xffff)))
	a = uint32(gm32.Round(gm32.Clamp(fa, 0, 0xffff)))

	return
}

func NormalizeRGBA(r, g, b, a uint32) (nr, ng, nb, na float32) {
	switch a {
	case 0xffff:
		nr = float32(r) * qf16
		ng = float32(g) * qf16
		nb = float32(b) * qf16
		na = 1
	default:
		q := 1.0 / float32(a)
		nr = float32(r) * q
		ng = float32(g) * q
		nb = float32(b) * q
		na = float32(a) * qf16
	}

	return
}

var (
	HSLAModel color.Model = color.ModelFunc(hslaModel)
	HSVAModel color.Model = color.ModelFunc(hsvaModel)
)

func hslaModel(c color.Color) color.Color {
	if _, ok := c.(HSLA); ok {
		return c
	}

	r, g, b, a := c.RGBA()
	nr, ng, nb, na := NormalizeRGBA(r, g, b, a)
	h, s, l := RGBToHSL(nr, ng, nb)

	return HSLA{h, s, l, na}
}

type HSVA struct {
	H, S, V, A float32
}

func (c HSVA) RGBA() (r, g, b, a uint32) {
	fr, fg, fb := HSVToRGB(c.H, c.S, c.V)
	fa := c.A * 0xffff

	fr = fr * fa
	fg = fg * fa
	fb = fb * fa

	r = uint32(gm32.Round(gm32.Clamp(fr, 0, 0xffff)))
	g = uint32(gm32.Round(gm32.Clamp(fg, 0, 0xffff)))
	b = uint32(gm32.Round(gm32.Clamp(fb, 0, 0xffff)))
	a = uint32(gm32.Round(gm32.Clamp(fa, 0, 0xffff)))

	return
}

func hsvaModel(c color.Color) color.Color {
	if _, ok := c.(HSVA); ok {
		return c
	}

	r, g, b, a := c.RGBA()
	nr, ng, nb, na := NormalizeRGBA(r, g, b, a)
	h, s, v := RGBToHSV(nr, ng, nb)

	return HSVA{h, s, v, na}
}
