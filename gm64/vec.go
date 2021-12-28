package gm64

import (
	"fmt"
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

type Vec struct {
	N    int
	Data []float64
}

func NewVec(n int) func(data ...float64) *Vec {
	if n <= 0 {
		err := fmt.Errorf("the n parameter must be positive (got %d)", n)
		panic(err)
	}

	ctor := func(data ...float64) *Vec {
		if len(data) > n {
			err := fmt.Errorf("the number of input values must not be greater than n (%d)", n)
			panic(err)
		}

		o := &Vec{
			N:    n,
			Data: make([]float64, n),
		}

		copy(o.Data, data)
		return o
	}

	return ctor
}

func (v *Vec) Copy() *Vec {
	cp := &Vec{
		N:    v.N,
		Data: make([]float64, v.N),
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
		Data: make([]float64, v1.N),
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
		Data: make([]float64, v1.N),
	}

	for i := 0; i < o.N; i++ {
		o.Data[i] = v1.Data[i] - v2.Data[i]
	}

	return o
}

func (v *Vec) Mul(c float64) *Vec {
	o := &Vec{
		N:    v.N,
		Data: make([]float64, v.N),
	}

	for i := 0; i < o.N; i++ {
		o.Data[i] = v.Data[i] * c
	}

	return o
}

func (v *Vec) Len() float64 {
	sum := float64(0)
	for i := 0; i < v.N; i++ {
		sum += v.Data[i] * v.Data[i]
	}
	return math.Sqrt(sum)
}

func (v *Vec) Normalize() *Vec {
	o := &Vec{
		N:    v.N,
		Data: make([]float64, v.N),
	}

	l := 1.0 / v.Len()

	for i := 0; i < o.N; i++ {
		o.Data[i] = v.Data[i] * l
	}

	return o
}

func (v1 *Vec) Dot(v2 *Vec) float64 {
	if v1.N != v2.N {
		err := fmt.Errorf(
			"the first and second vectors have different sized (got %d and %d)",
			v1.N, v2.N,
		)
		panic(err)
	}

	sum := float64(0)
	for i := 0; i < v1.N; i++ {
		sum += v1.Data[i] * v2.Data[i]
	}

	return sum
}

func (m *Mat) Row(i int) *Vec {
	if i < 0 {
		err := fmt.Errorf("the i parameter must be non-negative (got %d)", i)
		panic(err)
	}

	if i >= m.M {
		err := fmt.Errorf("trying to get a row out of matrix bounds (got row index %d, while matrix has only %d rows)", i, m.M)
		panic(err)
	}

	return NewVec(m.N)(m.Data[i*m.N : (i+1)*m.N]...)
}

func (m *Mat) Col(j int) *Vec {
	if j < 0 {
		err := fmt.Errorf("the j parameter must be non-negative (got %d)", j)
		panic(err)
	}

	if j >= m.N {
		err := fmt.Errorf("trying to get a column out of matrix bounds (got column index %d, while matrix has only %d columns)", j, m.N)
		panic(err)
	}

	col := make([]float64, m.M)
	for i := 0; i < m.M; i++ {
		col[i] = m.Data[j+i*m.N]
	}

	return NewVec(m.M)(col...)
}

func (v *Vec) String() string {
	return fmt.Sprint(v.Data)
}
