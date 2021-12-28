package gm32

type Vec2 [2]float32

func (v Vec2) Len() float32 {
	return Hypot(v[0], v[1])
}

func (v Vec2) Normalize() Vec2 {
	l := 1.0 / v.Len()
	return Vec2{v[0] * l, v[1] * l}
}

func (v1 Vec2) Dot(v2 Vec2) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func (v1 Vec2) Cross(v2 Vec2) float32 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1[0] + v2[0], v1[1] + v2[1]}
}

func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v1[0] - v2[0], v1[1] - v2[1]}
}

func (v1 Vec2) Mul(v2 Vec2) Vec2 {
	return Vec2{v1[0] * v2[0], v1[1] * v2[1]}
}

func (v Vec2) Elem() (x, y float32) {
	return v[0], v[1]
}
