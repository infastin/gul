package gft

import (
	"github.com/infastin/gul/gm32"
)

type Position struct {
	X, Y float32
}

type Size struct {
	Width, Height float32
}

func f32u8(val float32) uint8 {
	fv := gm32.Clamp(val, 0, 0xff)
	return uint8(gm32.Round(fv))
}

func f32u16(val float32) uint16 {
	fv := gm32.Clamp(val, 0, 0xffff)
	return uint16(gm32.Round(fv))
}
