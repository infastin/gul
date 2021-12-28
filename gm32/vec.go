package gm32

import "fmt"

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

func (v Vec2) Mul(c float32) Vec2 {
	return Vec2{v[0] * c, v[1] * c}
}

func (v Vec2) Elem() (x, y float32) {
	return v[0], v[1]
}

type Vec3 [3]float32

func (v Vec3) Len() float32 {
	return Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v Vec3) Normalize() Vec3 {
	l := 1.0 / v.Len()
	return Vec3{v[0] * l, v[1] * l, v[2] * l}
}

func (v1 Vec3) Dot(v2 Vec3) float32 {
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

func (v Vec3) Mul(c float32) Vec3 {
	return Vec3{v[0] * c, v[1] * c, v[2] * c}
}

func (v Vec3) Elem() (x, y, z float32) {
	return v[0], v[1], v[2]
}

type Vec4 [4]float32

func (v Vec4) Len() float32 {
	return Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])
}

func (v Vec4) Normalize() Vec4 {
	l := 1.0 / v.Len()
	return Vec4{v[0] * l, v[1] * l, v[2] * l, v[3] * l}
}

func (v1 Vec4) Dot(v2 Vec4) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2] + v1[3]*v2[3]
}

func (v1 Vec4) Add(v2 Vec4) Vec4 {
	return Vec4{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2], v1[3] + v2[3]}
}

func (v1 Vec4) Sub(v2 Vec4) Vec4 {
	return Vec4{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2], v1[3] - v2[3]}
}

func (v Vec4) Mul(c float32) Vec4 {
	return Vec4{v[0] * c, v[1] * c, v[2] * c, v[3] * c}
}

func (v Vec4) Elem() (x, y, z, w float32) {
	return v[0], v[1], v[2], v[3]
}

type Vec struct {
	N    int
	Data []float32
}

func NewVec(n int) func(data ...float32) *Vec {
	if n <= 0 {
		err := fmt.Errorf("the n parameter must be positive (got %d)", n)
		panic(err)
	}

	ctor := func(data ...float32) *Vec {
		if len(data) > n {
			err := fmt.Errorf("the number of input values must not be greater than n (%d)", n)
			panic(err)
		}

		o := &Vec{
			N:    n,
			Data: make([]float32, n),
		}

		copy(o.Data, data)
		return o
	}

	return ctor
}

func (v *Vec) Copy() *Vec {
	cp := &Vec{
		N:    v.N,
		Data: make([]float32, v.N),
	}

	copy(cp.Data, v.Data)

	return cp
}

func (v1 *Vec) Add(v2 *Vec) *Vec {
	if v1.N != v2.N {
		err := fmt.Errorf(
			"the first and second vectors have different sized (got %d and %d)",
			v1.N, v2.N,
		)
		panic(err)
	}

	o := &Vec{
		N:    v1.N,
		Data: make([]float32, v1.N),
	}

	for i := 0; i < o.N; i++ {
		o.Data[i] = v1.Data[i] + v2.Data[i]
	}

	return o
}

func (v1 *Vec) Sub(v2 *Vec) *Vec {
	if v1.N != v2.N {
		err := fmt.Errorf(
			"the first and second vectors have different sized (got %d and %d)",
			v1.N, v2.N,
		)
		panic(err)
	}

	o := &Vec{
		N:    v1.N,
		Data: make([]float32, v1.N),
	}

	for i := 0; i < o.N; i++ {
		o.Data[i] = v1.Data[i] - v2.Data[i]
	}

	return o
}

func (v *Vec) Mul(c float32) *Vec {
	o := &Vec{
		N:    v.N,
		Data: make([]float32, v.N),
	}

	for i := 0; i < o.N; i++ {
		o.Data[i] = v.Data[i] * c
	}

	return o
}

func (v *Vec) Len() float32 {
	sum := float32(0)
	for i := 0; i < v.N; i++ {
		sum += v.Data[i] * v.Data[i]
	}
	return Sqrt(sum)
}

func (v *Vec) Normalize() *Vec {
	o := &Vec{
		N:    v.N,
		Data: make([]float32, v.N),
	}

	l := 1.0 / v.Len()

	for i := 0; i < o.N; i++ {
		o.Data[i] = v.Data[i] * l
	}

	return o
}

func (v1 *Vec) Dot(v2 *Vec) float32 {
	if v1.N != v2.N {
		err := fmt.Errorf(
			"the first and second vectors have different sized (got %d and %d)",
			v1.N, v2.N,
		)
		panic(err)
	}

	sum := float32(0)
	for i := 0; i < v1.N; i++ {
		sum += v1.Data[i] * v2.Data[i]
	}

	return sum
}

func (v *Vec) String() string {
	return fmt.Sprint(v.Data)
}
