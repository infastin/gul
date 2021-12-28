package gm64

import (
	"math"
)

type Vec2 [2]float64

func (v Vec2) Len() float64 {
	return math.Hypot(v[0], v[1])
}

func (v Vec2) Normalize() Vec2 {
	l := 1.0 / v.Len()
	return Vec2{v[0] * l, v[1] * l}
}

func (v1 Vec2) Dot(v2 Vec2) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func (v1 Vec2) Cross(v2 Vec2) float64 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1[0] + v2[0], v1[1] + v2[1]}
}

func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v1[0] - v2[0], v1[1] - v2[1]}
}

func (v Vec2) Mul(c float64) Vec2 {
	return Vec2{v[0] * c, v[1] * c}
}

func (v Vec2) Elem() (x, y float64) {
	return v[0], v[1]
}

type Vec3 [3]float64

func (v Vec3) Len() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v Vec3) Normalize() Vec3 {
	l := 1.0 / v.Len()
	return Vec3{v[0] * l, v[1] * l, v[2] * l}
}

func (v1 Vec3) Dot(v2 Vec3) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{v1[1]*v2[2] - v1[2]*v2[1], v1[2]*v2[0] - v1[0]*v2[2], v1[0]*v2[1] - v2[1]*v1[0]}
}

func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

func (v Vec3) Mul(c float64) Vec3 {
	return Vec3{v[0] * c, v[1] * c, v[2] * c}
}

func (v Vec3) Elem() (x, y, z float64) {
	return v[0], v[1], v[2]
}

type Vec4 [4]float64

func (v Vec4) Len() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])
}

func (v Vec4) Normalize() Vec4 {
	l := 1.0 / v.Len()
	return Vec4{v[0] * l, v[1] * l, v[2] * l, v[3] * l}
}

func (v1 Vec4) Dot(v2 Vec4) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2] + v1[3]*v2[3]
}

func (v1 Vec4) Add(v2 Vec4) Vec4 {
	return Vec4{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2], v1[3] + v2[3]}
}

func (v1 Vec4) Sub(v2 Vec4) Vec4 {
	return Vec4{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2], v1[3] - v2[3]}
}

func (v Vec4) Mul(c float64) Vec4 {
	return Vec4{v[0] * c, v[1] * c, v[2] * c, v[3] * c}
}

func (v Vec4) Elem() (x, y, z, w float64) {
	return v[0], v[1], v[2], v[3]
}
