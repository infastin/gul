// This file is generated by codegen.go; DO NOT EDIT!
package gm64

import "fmt"

type Mat2 [4]float64

func (m1 Mat2) Add(m2 Mat2) Mat2 {
	return Mat2{
		m1[0] + m2[0], m1[1] + m2[1],
		m1[2] + m2[2], m1[3] + m2[3],
	}
}

func (m1 Mat2) Sub(m2 Mat2) Mat2 {
	return Mat2{
		m1[0] - m2[0], m1[1] - m2[1],
		m1[2] - m2[2], m1[3] - m2[3],
	}
}

func (m Mat2) Mul(c float64) Mat2 {
	return Mat2{
		m[0] * c, m[1] * c,
		m[2] * c, m[3] * c,
	}
}

func (m1 Mat2) MulMat2x1(m2 Vec2) Vec2 {
	return Vec2{
		m1[0]*m2[0] + m1[1]*m2[1],
		m1[2]*m2[0] + m1[3]*m2[1],
	}
}

func (m1 Mat2) MulMat2(m2 Mat2) Mat2 {
	return Mat2{
		m1[0]*m2[0] + m1[1]*m2[2], m1[0]*m2[1] + m1[1]*m2[3],
		m1[2]*m2[0] + m1[3]*m2[2], m1[2]*m2[1] + m1[3]*m2[3],
	}
}

func (m1 Mat2) MulMat2x3(m2 Mat2x3) Mat2x3 {
	return Mat2x3{
		m1[0]*m2[0] + m1[1]*m2[3], m1[0]*m2[1] + m1[1]*m2[4], m1[0]*m2[2] + m1[1]*m2[5],
		m1[2]*m2[0] + m1[3]*m2[3], m1[2]*m2[1] + m1[3]*m2[4], m1[2]*m2[2] + m1[3]*m2[5],
	}
}

func (m1 Mat2) MulMat2x4(m2 Mat2x4) Mat2x4 {
	return Mat2x4{
		m1[0]*m2[0] + m1[1]*m2[4], m1[0]*m2[1] + m1[1]*m2[5], m1[0]*m2[2] + m1[1]*m2[6], m1[0]*m2[3] + m1[1]*m2[7],
		m1[2]*m2[0] + m1[3]*m2[4], m1[2]*m2[1] + m1[3]*m2[5], m1[2]*m2[2] + m1[3]*m2[6], m1[2]*m2[3] + m1[3]*m2[7],
	}
}

func (m Mat2) Trace() float64 {
	return m[0] + m[3]
}

func (m Mat2) Det() float64 {
	return m[0]*m[3] - m[1]*m[2]
}

func (m Mat2) At(i, j int) float64 {
	if i >= 2 || j >= 2 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (2x2))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*2]
}

func (m Mat2) Set(i, j int, value float64) {
	if i >= 2 || j >= 2 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (2x2))",
			i, j,
		)
		panic(err)
	}

	m[j+i*2] = value
}

type Mat2x3 [6]float64

func (m1 Mat2x3) Add(m2 Mat2x3) Mat2x3 {
	return Mat2x3{
		m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2],
		m1[3] + m2[3], m1[4] + m2[4], m1[5] + m2[5],
	}
}

func (m1 Mat2x3) Sub(m2 Mat2x3) Mat2x3 {
	return Mat2x3{
		m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2],
		m1[3] - m2[3], m1[4] - m2[4], m1[5] - m2[5],
	}
}

func (m Mat2x3) Mul(c float64) Mat2x3 {
	return Mat2x3{
		m[0] * c, m[1] * c, m[2] * c,
		m[3] * c, m[4] * c, m[5] * c,
	}
}

func (m1 Mat2x3) MulMat3x1(m2 Vec3) Vec2 {
	return Vec2{
		m1[0]*m2[0] + m1[1]*m2[1] + m1[2]*m2[2],
		m1[3]*m2[0] + m1[4]*m2[1] + m1[5]*m2[2],
	}
}

func (m1 Mat2x3) MulMat3x2(m2 Mat3x2) Mat2 {
	return Mat2{
		m1[0]*m2[0] + m1[1]*m2[2] + m1[2]*m2[4], m1[0]*m2[1] + m1[1]*m2[3] + m1[2]*m2[5],
		m1[3]*m2[0] + m1[4]*m2[2] + m1[5]*m2[4], m1[3]*m2[1] + m1[4]*m2[3] + m1[5]*m2[5],
	}
}

func (m1 Mat2x3) MulMat3(m2 Mat3) Mat2x3 {
	return Mat2x3{
		m1[0]*m2[0] + m1[1]*m2[3] + m1[2]*m2[6], m1[0]*m2[1] + m1[1]*m2[4] + m1[2]*m2[7], m1[0]*m2[2] + m1[1]*m2[5] + m1[2]*m2[8],
		m1[3]*m2[0] + m1[4]*m2[3] + m1[5]*m2[6], m1[3]*m2[1] + m1[4]*m2[4] + m1[5]*m2[7], m1[3]*m2[2] + m1[4]*m2[5] + m1[5]*m2[8],
	}
}

func (m1 Mat2x3) MulMat3x4(m2 Mat3x4) Mat2x4 {
	return Mat2x4{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8], m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9], m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10], m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11],
		m1[3]*m2[0] + m1[4]*m2[4] + m1[5]*m2[8], m1[3]*m2[1] + m1[4]*m2[5] + m1[5]*m2[9], m1[3]*m2[2] + m1[4]*m2[6] + m1[5]*m2[10], m1[3]*m2[3] + m1[4]*m2[7] + m1[5]*m2[11],
	}
}

func (m Mat2x3) At(i, j int) float64 {
	if i >= 2 || j >= 3 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (2x3))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*3]
}

func (m Mat2x3) Set(i, j int, value float64) {
	if i >= 2 || j >= 3 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (2x3))",
			i, j,
		)
		panic(err)
	}

	m[j+i*3] = value
}

type Mat2x4 [8]float64

func (m1 Mat2x4) Add(m2 Mat2x4) Mat2x4 {
	return Mat2x4{
		m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2], m1[3] + m2[3],
		m1[4] + m2[4], m1[5] + m2[5], m1[6] + m2[6], m1[7] + m2[7],
	}
}

func (m1 Mat2x4) Sub(m2 Mat2x4) Mat2x4 {
	return Mat2x4{
		m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2], m1[3] - m2[3],
		m1[4] - m2[4], m1[5] - m2[5], m1[6] - m2[6], m1[7] - m2[7],
	}
}

func (m Mat2x4) Mul(c float64) Mat2x4 {
	return Mat2x4{
		m[0] * c, m[1] * c, m[2] * c, m[3] * c,
		m[4] * c, m[5] * c, m[6] * c, m[7] * c,
	}
}

func (m1 Mat2x4) MulMat4x1(m2 Vec4) Vec2 {
	return Vec2{
		m1[0]*m2[0] + m1[1]*m2[1] + m1[2]*m2[2] + m1[3]*m2[3],
		m1[4]*m2[0] + m1[5]*m2[1] + m1[6]*m2[2] + m1[7]*m2[3],
	}
}

func (m1 Mat2x4) MulMat4x2(m2 Mat4x2) Mat2 {
	return Mat2{
		m1[0]*m2[0] + m1[1]*m2[2] + m1[2]*m2[4] + m1[3]*m2[6], m1[0]*m2[1] + m1[1]*m2[3] + m1[2]*m2[5] + m1[3]*m2[7],
		m1[4]*m2[0] + m1[5]*m2[2] + m1[6]*m2[4] + m1[7]*m2[6], m1[4]*m2[1] + m1[5]*m2[3] + m1[6]*m2[5] + m1[7]*m2[7],
	}
}

func (m1 Mat2x4) MulMat4x3(m2 Mat4x3) Mat2x3 {
	return Mat2x3{
		m1[0]*m2[0] + m1[1]*m2[3] + m1[2]*m2[6] + m1[3]*m2[9], m1[0]*m2[1] + m1[1]*m2[4] + m1[2]*m2[7] + m1[3]*m2[10], m1[0]*m2[2] + m1[1]*m2[5] + m1[2]*m2[8] + m1[3]*m2[11],
		m1[4]*m2[0] + m1[5]*m2[3] + m1[6]*m2[6] + m1[7]*m2[9], m1[4]*m2[1] + m1[5]*m2[4] + m1[6]*m2[7] + m1[7]*m2[10], m1[4]*m2[2] + m1[5]*m2[5] + m1[6]*m2[8] + m1[7]*m2[11],
	}
}

func (m1 Mat2x4) MulMat4(m2 Mat4) Mat2x4 {
	return Mat2x4{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8] + m1[3]*m2[12], m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9] + m1[3]*m2[13], m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10] + m1[3]*m2[14], m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11] + m1[3]*m2[15],
		m1[4]*m2[0] + m1[5]*m2[4] + m1[6]*m2[8] + m1[7]*m2[12], m1[4]*m2[1] + m1[5]*m2[5] + m1[6]*m2[9] + m1[7]*m2[13], m1[4]*m2[2] + m1[5]*m2[6] + m1[6]*m2[10] + m1[7]*m2[14], m1[4]*m2[3] + m1[5]*m2[7] + m1[6]*m2[11] + m1[7]*m2[15],
	}
}

func (m Mat2x4) At(i, j int) float64 {
	if i >= 2 || j >= 4 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (2x4))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*4]
}

func (m Mat2x4) Set(i, j int, value float64) {
	if i >= 2 || j >= 4 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (2x4))",
			i, j,
		)
		panic(err)
	}

	m[j+i*4] = value
}

type Mat3x2 [6]float64

func (m1 Mat3x2) Add(m2 Mat3x2) Mat3x2 {
	return Mat3x2{
		m1[0] + m2[0], m1[1] + m2[1],
		m1[2] + m2[2], m1[3] + m2[3],
		m1[4] + m2[4], m1[5] + m2[5],
	}
}

func (m1 Mat3x2) Sub(m2 Mat3x2) Mat3x2 {
	return Mat3x2{
		m1[0] - m2[0], m1[1] - m2[1],
		m1[2] - m2[2], m1[3] - m2[3],
		m1[4] - m2[4], m1[5] - m2[5],
	}
}

func (m Mat3x2) Mul(c float64) Mat3x2 {
	return Mat3x2{
		m[0] * c, m[1] * c,
		m[2] * c, m[3] * c,
		m[4] * c, m[5] * c,
	}
}

func (m1 Mat3x2) MulMat2x1(m2 Vec2) Vec3 {
	return Vec3{
		m1[0]*m2[0] + m1[1]*m2[1],
		m1[2]*m2[0] + m1[3]*m2[1],
		m1[4]*m2[0] + m1[5]*m2[1],
	}
}

func (m1 Mat3x2) MulMat2(m2 Mat2) Mat3x2 {
	return Mat3x2{
		m1[0]*m2[0] + m1[1]*m2[2], m1[0]*m2[1] + m1[1]*m2[3],
		m1[2]*m2[0] + m1[3]*m2[2], m1[2]*m2[1] + m1[3]*m2[3],
		m1[4]*m2[0] + m1[5]*m2[2], m1[4]*m2[1] + m1[5]*m2[3],
	}
}

func (m1 Mat3x2) MulMat2x3(m2 Mat2x3) Mat3 {
	return Mat3{
		m1[0]*m2[0] + m1[1]*m2[3], m1[0]*m2[1] + m1[1]*m2[4], m1[0]*m2[2] + m1[1]*m2[5],
		m1[2]*m2[0] + m1[3]*m2[3], m1[2]*m2[1] + m1[3]*m2[4], m1[2]*m2[2] + m1[3]*m2[5],
		m1[4]*m2[0] + m1[5]*m2[3], m1[4]*m2[1] + m1[5]*m2[4], m1[4]*m2[2] + m1[5]*m2[5],
	}
}

func (m1 Mat3x2) MulMat2x4(m2 Mat2x4) Mat3x4 {
	return Mat3x4{
		m1[0]*m2[0] + m1[1]*m2[4], m1[0]*m2[1] + m1[1]*m2[5], m1[0]*m2[2] + m1[1]*m2[6], m1[0]*m2[3] + m1[1]*m2[7],
		m1[2]*m2[0] + m1[3]*m2[4], m1[2]*m2[1] + m1[3]*m2[5], m1[2]*m2[2] + m1[3]*m2[6], m1[2]*m2[3] + m1[3]*m2[7],
		m1[4]*m2[0] + m1[5]*m2[4], m1[4]*m2[1] + m1[5]*m2[5], m1[4]*m2[2] + m1[5]*m2[6], m1[4]*m2[3] + m1[5]*m2[7],
	}
}

func (m Mat3x2) At(i, j int) float64 {
	if i >= 3 || j >= 2 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (3x2))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*2]
}

func (m Mat3x2) Set(i, j int, value float64) {
	if i >= 3 || j >= 2 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (3x2))",
			i, j,
		)
		panic(err)
	}

	m[j+i*2] = value
}

type Mat3 [9]float64

func (m1 Mat3) Add(m2 Mat3) Mat3 {
	return Mat3{
		m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2],
		m1[3] + m2[3], m1[4] + m2[4], m1[5] + m2[5],
		m1[6] + m2[6], m1[7] + m2[7], m1[8] + m2[8],
	}
}

func (m1 Mat3) Sub(m2 Mat3) Mat3 {
	return Mat3{
		m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2],
		m1[3] - m2[3], m1[4] - m2[4], m1[5] - m2[5],
		m1[6] - m2[6], m1[7] - m2[7], m1[8] - m2[8],
	}
}

func (m Mat3) Mul(c float64) Mat3 {
	return Mat3{
		m[0] * c, m[1] * c, m[2] * c,
		m[3] * c, m[4] * c, m[5] * c,
		m[6] * c, m[7] * c, m[8] * c,
	}
}

func (m1 Mat3) MulMat3x1(m2 Vec3) Vec3 {
	return Vec3{
		m1[0]*m2[0] + m1[1]*m2[1] + m1[2]*m2[2],
		m1[3]*m2[0] + m1[4]*m2[1] + m1[5]*m2[2],
		m1[6]*m2[0] + m1[7]*m2[1] + m1[8]*m2[2],
	}
}

func (m1 Mat3) MulMat3x2(m2 Mat3x2) Mat3x2 {
	return Mat3x2{
		m1[0]*m2[0] + m1[1]*m2[2] + m1[2]*m2[4], m1[0]*m2[1] + m1[1]*m2[3] + m1[2]*m2[5],
		m1[3]*m2[0] + m1[4]*m2[2] + m1[5]*m2[4], m1[3]*m2[1] + m1[4]*m2[3] + m1[5]*m2[5],
		m1[6]*m2[0] + m1[7]*m2[2] + m1[8]*m2[4], m1[6]*m2[1] + m1[7]*m2[3] + m1[8]*m2[5],
	}
}

func (m1 Mat3) MulMat3(m2 Mat3) Mat3 {
	return Mat3{
		m1[0]*m2[0] + m1[1]*m2[3] + m1[2]*m2[6], m1[0]*m2[1] + m1[1]*m2[4] + m1[2]*m2[7], m1[0]*m2[2] + m1[1]*m2[5] + m1[2]*m2[8],
		m1[3]*m2[0] + m1[4]*m2[3] + m1[5]*m2[6], m1[3]*m2[1] + m1[4]*m2[4] + m1[5]*m2[7], m1[3]*m2[2] + m1[4]*m2[5] + m1[5]*m2[8],
		m1[6]*m2[0] + m1[7]*m2[3] + m1[8]*m2[6], m1[6]*m2[1] + m1[7]*m2[4] + m1[8]*m2[7], m1[6]*m2[2] + m1[7]*m2[5] + m1[8]*m2[8],
	}
}

func (m1 Mat3) MulMat3x4(m2 Mat3x4) Mat3x4 {
	return Mat3x4{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8], m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9], m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10], m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11],
		m1[3]*m2[0] + m1[4]*m2[4] + m1[5]*m2[8], m1[3]*m2[1] + m1[4]*m2[5] + m1[5]*m2[9], m1[3]*m2[2] + m1[4]*m2[6] + m1[5]*m2[10], m1[3]*m2[3] + m1[4]*m2[7] + m1[5]*m2[11],
		m1[6]*m2[0] + m1[7]*m2[4] + m1[8]*m2[8], m1[6]*m2[1] + m1[7]*m2[5] + m1[8]*m2[9], m1[6]*m2[2] + m1[7]*m2[6] + m1[8]*m2[10], m1[6]*m2[3] + m1[7]*m2[7] + m1[8]*m2[11],
	}
}

func (m Mat3) Trace() float64 {
	return m[0] + m[4] + m[8]
}

func (m Mat3) Det() float64 {
	return m[0]*m[4]*m[8] - m[0]*m[5]*m[7] - m[1]*m[3]*m[8] + m[1]*m[5]*m[6] + m[2]*m[3]*m[7] - m[2]*m[4]*m[6]
}

func (m Mat3) At(i, j int) float64 {
	if i >= 3 || j >= 3 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (3x3))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*3]
}

func (m Mat3) Set(i, j int, value float64) {
	if i >= 3 || j >= 3 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (3x3))",
			i, j,
		)
		panic(err)
	}

	m[j+i*3] = value
}

type Mat3x4 [12]float64

func (m1 Mat3x4) Add(m2 Mat3x4) Mat3x4 {
	return Mat3x4{
		m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2], m1[3] + m2[3],
		m1[4] + m2[4], m1[5] + m2[5], m1[6] + m2[6], m1[7] + m2[7],
		m1[8] + m2[8], m1[9] + m2[9], m1[10] + m2[10], m1[11] + m2[11],
	}
}

func (m1 Mat3x4) Sub(m2 Mat3x4) Mat3x4 {
	return Mat3x4{
		m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2], m1[3] - m2[3],
		m1[4] - m2[4], m1[5] - m2[5], m1[6] - m2[6], m1[7] - m2[7],
		m1[8] - m2[8], m1[9] - m2[9], m1[10] - m2[10], m1[11] - m2[11],
	}
}

func (m Mat3x4) Mul(c float64) Mat3x4 {
	return Mat3x4{
		m[0] * c, m[1] * c, m[2] * c, m[3] * c,
		m[4] * c, m[5] * c, m[6] * c, m[7] * c,
		m[8] * c, m[9] * c, m[10] * c, m[11] * c,
	}
}

func (m1 Mat3x4) MulMat4x1(m2 Vec4) Vec3 {
	return Vec3{
		m1[0]*m2[0] + m1[1]*m2[1] + m1[2]*m2[2] + m1[3]*m2[3],
		m1[4]*m2[0] + m1[5]*m2[1] + m1[6]*m2[2] + m1[7]*m2[3],
		m1[8]*m2[0] + m1[9]*m2[1] + m1[10]*m2[2] + m1[11]*m2[3],
	}
}

func (m1 Mat3x4) MulMat4x2(m2 Mat4x2) Mat3x2 {
	return Mat3x2{
		m1[0]*m2[0] + m1[1]*m2[2] + m1[2]*m2[4] + m1[3]*m2[6], m1[0]*m2[1] + m1[1]*m2[3] + m1[2]*m2[5] + m1[3]*m2[7],
		m1[4]*m2[0] + m1[5]*m2[2] + m1[6]*m2[4] + m1[7]*m2[6], m1[4]*m2[1] + m1[5]*m2[3] + m1[6]*m2[5] + m1[7]*m2[7],
		m1[8]*m2[0] + m1[9]*m2[2] + m1[10]*m2[4] + m1[11]*m2[6], m1[8]*m2[1] + m1[9]*m2[3] + m1[10]*m2[5] + m1[11]*m2[7],
	}
}

func (m1 Mat3x4) MulMat4x3(m2 Mat4x3) Mat3 {
	return Mat3{
		m1[0]*m2[0] + m1[1]*m2[3] + m1[2]*m2[6] + m1[3]*m2[9], m1[0]*m2[1] + m1[1]*m2[4] + m1[2]*m2[7] + m1[3]*m2[10], m1[0]*m2[2] + m1[1]*m2[5] + m1[2]*m2[8] + m1[3]*m2[11],
		m1[4]*m2[0] + m1[5]*m2[3] + m1[6]*m2[6] + m1[7]*m2[9], m1[4]*m2[1] + m1[5]*m2[4] + m1[6]*m2[7] + m1[7]*m2[10], m1[4]*m2[2] + m1[5]*m2[5] + m1[6]*m2[8] + m1[7]*m2[11],
		m1[8]*m2[0] + m1[9]*m2[3] + m1[10]*m2[6] + m1[11]*m2[9], m1[8]*m2[1] + m1[9]*m2[4] + m1[10]*m2[7] + m1[11]*m2[10], m1[8]*m2[2] + m1[9]*m2[5] + m1[10]*m2[8] + m1[11]*m2[11],
	}
}

func (m1 Mat3x4) MulMat4(m2 Mat4) Mat3x4 {
	return Mat3x4{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8] + m1[3]*m2[12], m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9] + m1[3]*m2[13], m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10] + m1[3]*m2[14], m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11] + m1[3]*m2[15],
		m1[4]*m2[0] + m1[5]*m2[4] + m1[6]*m2[8] + m1[7]*m2[12], m1[4]*m2[1] + m1[5]*m2[5] + m1[6]*m2[9] + m1[7]*m2[13], m1[4]*m2[2] + m1[5]*m2[6] + m1[6]*m2[10] + m1[7]*m2[14], m1[4]*m2[3] + m1[5]*m2[7] + m1[6]*m2[11] + m1[7]*m2[15],
		m1[8]*m2[0] + m1[9]*m2[4] + m1[10]*m2[8] + m1[11]*m2[12], m1[8]*m2[1] + m1[9]*m2[5] + m1[10]*m2[9] + m1[11]*m2[13], m1[8]*m2[2] + m1[9]*m2[6] + m1[10]*m2[10] + m1[11]*m2[14], m1[8]*m2[3] + m1[9]*m2[7] + m1[10]*m2[11] + m1[11]*m2[15],
	}
}

func (m Mat3x4) At(i, j int) float64 {
	if i >= 3 || j >= 4 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (3x4))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*4]
}

func (m Mat3x4) Set(i, j int, value float64) {
	if i >= 3 || j >= 4 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (3x4))",
			i, j,
		)
		panic(err)
	}

	m[j+i*4] = value
}

type Mat4x2 [8]float64

func (m1 Mat4x2) Add(m2 Mat4x2) Mat4x2 {
	return Mat4x2{
		m1[0] + m2[0], m1[1] + m2[1],
		m1[2] + m2[2], m1[3] + m2[3],
		m1[4] + m2[4], m1[5] + m2[5],
		m1[6] + m2[6], m1[7] + m2[7],
	}
}

func (m1 Mat4x2) Sub(m2 Mat4x2) Mat4x2 {
	return Mat4x2{
		m1[0] - m2[0], m1[1] - m2[1],
		m1[2] - m2[2], m1[3] - m2[3],
		m1[4] - m2[4], m1[5] - m2[5],
		m1[6] - m2[6], m1[7] - m2[7],
	}
}

func (m Mat4x2) Mul(c float64) Mat4x2 {
	return Mat4x2{
		m[0] * c, m[1] * c,
		m[2] * c, m[3] * c,
		m[4] * c, m[5] * c,
		m[6] * c, m[7] * c,
	}
}

func (m1 Mat4x2) MulMat2x1(m2 Vec2) Vec4 {
	return Vec4{
		m1[0]*m2[0] + m1[1]*m2[1],
		m1[2]*m2[0] + m1[3]*m2[1],
		m1[4]*m2[0] + m1[5]*m2[1],
		m1[6]*m2[0] + m1[7]*m2[1],
	}
}

func (m1 Mat4x2) MulMat2(m2 Mat2) Mat4x2 {
	return Mat4x2{
		m1[0]*m2[0] + m1[1]*m2[2], m1[0]*m2[1] + m1[1]*m2[3],
		m1[2]*m2[0] + m1[3]*m2[2], m1[2]*m2[1] + m1[3]*m2[3],
		m1[4]*m2[0] + m1[5]*m2[2], m1[4]*m2[1] + m1[5]*m2[3],
		m1[6]*m2[0] + m1[7]*m2[2], m1[6]*m2[1] + m1[7]*m2[3],
	}
}

func (m1 Mat4x2) MulMat2x3(m2 Mat2x3) Mat4x3 {
	return Mat4x3{
		m1[0]*m2[0] + m1[1]*m2[3], m1[0]*m2[1] + m1[1]*m2[4], m1[0]*m2[2] + m1[1]*m2[5],
		m1[2]*m2[0] + m1[3]*m2[3], m1[2]*m2[1] + m1[3]*m2[4], m1[2]*m2[2] + m1[3]*m2[5],
		m1[4]*m2[0] + m1[5]*m2[3], m1[4]*m2[1] + m1[5]*m2[4], m1[4]*m2[2] + m1[5]*m2[5],
		m1[6]*m2[0] + m1[7]*m2[3], m1[6]*m2[1] + m1[7]*m2[4], m1[6]*m2[2] + m1[7]*m2[5],
	}
}

func (m1 Mat4x2) MulMat2x4(m2 Mat2x4) Mat4 {
	return Mat4{
		m1[0]*m2[0] + m1[1]*m2[4], m1[0]*m2[1] + m1[1]*m2[5], m1[0]*m2[2] + m1[1]*m2[6], m1[0]*m2[3] + m1[1]*m2[7],
		m1[2]*m2[0] + m1[3]*m2[4], m1[2]*m2[1] + m1[3]*m2[5], m1[2]*m2[2] + m1[3]*m2[6], m1[2]*m2[3] + m1[3]*m2[7],
		m1[4]*m2[0] + m1[5]*m2[4], m1[4]*m2[1] + m1[5]*m2[5], m1[4]*m2[2] + m1[5]*m2[6], m1[4]*m2[3] + m1[5]*m2[7],
		m1[6]*m2[0] + m1[7]*m2[4], m1[6]*m2[1] + m1[7]*m2[5], m1[6]*m2[2] + m1[7]*m2[6], m1[6]*m2[3] + m1[7]*m2[7],
	}
}

func (m Mat4x2) At(i, j int) float64 {
	if i >= 4 || j >= 2 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (4x2))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*2]
}

func (m Mat4x2) Set(i, j int, value float64) {
	if i >= 4 || j >= 2 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (4x2))",
			i, j,
		)
		panic(err)
	}

	m[j+i*2] = value
}

type Mat4x3 [12]float64

func (m1 Mat4x3) Add(m2 Mat4x3) Mat4x3 {
	return Mat4x3{
		m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2],
		m1[3] + m2[3], m1[4] + m2[4], m1[5] + m2[5],
		m1[6] + m2[6], m1[7] + m2[7], m1[8] + m2[8],
		m1[9] + m2[9], m1[10] + m2[10], m1[11] + m2[11],
	}
}

func (m1 Mat4x3) Sub(m2 Mat4x3) Mat4x3 {
	return Mat4x3{
		m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2],
		m1[3] - m2[3], m1[4] - m2[4], m1[5] - m2[5],
		m1[6] - m2[6], m1[7] - m2[7], m1[8] - m2[8],
		m1[9] - m2[9], m1[10] - m2[10], m1[11] - m2[11],
	}
}

func (m Mat4x3) Mul(c float64) Mat4x3 {
	return Mat4x3{
		m[0] * c, m[1] * c, m[2] * c,
		m[3] * c, m[4] * c, m[5] * c,
		m[6] * c, m[7] * c, m[8] * c,
		m[9] * c, m[10] * c, m[11] * c,
	}
}

func (m1 Mat4x3) MulMat3x1(m2 Vec3) Vec4 {
	return Vec4{
		m1[0]*m2[0] + m1[1]*m2[1] + m1[2]*m2[2],
		m1[3]*m2[0] + m1[4]*m2[1] + m1[5]*m2[2],
		m1[6]*m2[0] + m1[7]*m2[1] + m1[8]*m2[2],
		m1[9]*m2[0] + m1[10]*m2[1] + m1[11]*m2[2],
	}
}

func (m1 Mat4x3) MulMat3x2(m2 Mat3x2) Mat4x2 {
	return Mat4x2{
		m1[0]*m2[0] + m1[1]*m2[2] + m1[2]*m2[4], m1[0]*m2[1] + m1[1]*m2[3] + m1[2]*m2[5],
		m1[3]*m2[0] + m1[4]*m2[2] + m1[5]*m2[4], m1[3]*m2[1] + m1[4]*m2[3] + m1[5]*m2[5],
		m1[6]*m2[0] + m1[7]*m2[2] + m1[8]*m2[4], m1[6]*m2[1] + m1[7]*m2[3] + m1[8]*m2[5],
		m1[9]*m2[0] + m1[10]*m2[2] + m1[11]*m2[4], m1[9]*m2[1] + m1[10]*m2[3] + m1[11]*m2[5],
	}
}

func (m1 Mat4x3) MulMat3(m2 Mat3) Mat4x3 {
	return Mat4x3{
		m1[0]*m2[0] + m1[1]*m2[3] + m1[2]*m2[6], m1[0]*m2[1] + m1[1]*m2[4] + m1[2]*m2[7], m1[0]*m2[2] + m1[1]*m2[5] + m1[2]*m2[8],
		m1[3]*m2[0] + m1[4]*m2[3] + m1[5]*m2[6], m1[3]*m2[1] + m1[4]*m2[4] + m1[5]*m2[7], m1[3]*m2[2] + m1[4]*m2[5] + m1[5]*m2[8],
		m1[6]*m2[0] + m1[7]*m2[3] + m1[8]*m2[6], m1[6]*m2[1] + m1[7]*m2[4] + m1[8]*m2[7], m1[6]*m2[2] + m1[7]*m2[5] + m1[8]*m2[8],
		m1[9]*m2[0] + m1[10]*m2[3] + m1[11]*m2[6], m1[9]*m2[1] + m1[10]*m2[4] + m1[11]*m2[7], m1[9]*m2[2] + m1[10]*m2[5] + m1[11]*m2[8],
	}
}

func (m1 Mat4x3) MulMat3x4(m2 Mat3x4) Mat4 {
	return Mat4{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8], m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9], m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10], m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11],
		m1[3]*m2[0] + m1[4]*m2[4] + m1[5]*m2[8], m1[3]*m2[1] + m1[4]*m2[5] + m1[5]*m2[9], m1[3]*m2[2] + m1[4]*m2[6] + m1[5]*m2[10], m1[3]*m2[3] + m1[4]*m2[7] + m1[5]*m2[11],
		m1[6]*m2[0] + m1[7]*m2[4] + m1[8]*m2[8], m1[6]*m2[1] + m1[7]*m2[5] + m1[8]*m2[9], m1[6]*m2[2] + m1[7]*m2[6] + m1[8]*m2[10], m1[6]*m2[3] + m1[7]*m2[7] + m1[8]*m2[11],
		m1[9]*m2[0] + m1[10]*m2[4] + m1[11]*m2[8], m1[9]*m2[1] + m1[10]*m2[5] + m1[11]*m2[9], m1[9]*m2[2] + m1[10]*m2[6] + m1[11]*m2[10], m1[9]*m2[3] + m1[10]*m2[7] + m1[11]*m2[11],
	}
}

func (m Mat4x3) At(i, j int) float64 {
	if i >= 4 || j >= 3 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (4x3))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*3]
}

func (m Mat4x3) Set(i, j int, value float64) {
	if i >= 4 || j >= 3 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (4x3))",
			i, j,
		)
		panic(err)
	}

	m[j+i*3] = value
}

type Mat4 [16]float64

func (m1 Mat4) Add(m2 Mat4) Mat4 {
	return Mat4{
		m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2], m1[3] + m2[3],
		m1[4] + m2[4], m1[5] + m2[5], m1[6] + m2[6], m1[7] + m2[7],
		m1[8] + m2[8], m1[9] + m2[9], m1[10] + m2[10], m1[11] + m2[11],
		m1[12] + m2[12], m1[13] + m2[13], m1[14] + m2[14], m1[15] + m2[15],
	}
}

func (m1 Mat4) Sub(m2 Mat4) Mat4 {
	return Mat4{
		m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2], m1[3] - m2[3],
		m1[4] - m2[4], m1[5] - m2[5], m1[6] - m2[6], m1[7] - m2[7],
		m1[8] - m2[8], m1[9] - m2[9], m1[10] - m2[10], m1[11] - m2[11],
		m1[12] - m2[12], m1[13] - m2[13], m1[14] - m2[14], m1[15] - m2[15],
	}
}

func (m Mat4) Mul(c float64) Mat4 {
	return Mat4{
		m[0] * c, m[1] * c, m[2] * c, m[3] * c,
		m[4] * c, m[5] * c, m[6] * c, m[7] * c,
		m[8] * c, m[9] * c, m[10] * c, m[11] * c,
		m[12] * c, m[13] * c, m[14] * c, m[15] * c,
	}
}

func (m1 Mat4) MulMat4x1(m2 Vec4) Vec4 {
	return Vec4{
		m1[0]*m2[0] + m1[1]*m2[1] + m1[2]*m2[2] + m1[3]*m2[3],
		m1[4]*m2[0] + m1[5]*m2[1] + m1[6]*m2[2] + m1[7]*m2[3],
		m1[8]*m2[0] + m1[9]*m2[1] + m1[10]*m2[2] + m1[11]*m2[3],
		m1[12]*m2[0] + m1[13]*m2[1] + m1[14]*m2[2] + m1[15]*m2[3],
	}
}

func (m1 Mat4) MulMat4x2(m2 Mat4x2) Mat4x2 {
	return Mat4x2{
		m1[0]*m2[0] + m1[1]*m2[2] + m1[2]*m2[4] + m1[3]*m2[6], m1[0]*m2[1] + m1[1]*m2[3] + m1[2]*m2[5] + m1[3]*m2[7],
		m1[4]*m2[0] + m1[5]*m2[2] + m1[6]*m2[4] + m1[7]*m2[6], m1[4]*m2[1] + m1[5]*m2[3] + m1[6]*m2[5] + m1[7]*m2[7],
		m1[8]*m2[0] + m1[9]*m2[2] + m1[10]*m2[4] + m1[11]*m2[6], m1[8]*m2[1] + m1[9]*m2[3] + m1[10]*m2[5] + m1[11]*m2[7],
		m1[12]*m2[0] + m1[13]*m2[2] + m1[14]*m2[4] + m1[15]*m2[6], m1[12]*m2[1] + m1[13]*m2[3] + m1[14]*m2[5] + m1[15]*m2[7],
	}
}

func (m1 Mat4) MulMat4x3(m2 Mat4x3) Mat4x3 {
	return Mat4x3{
		m1[0]*m2[0] + m1[1]*m2[3] + m1[2]*m2[6] + m1[3]*m2[9], m1[0]*m2[1] + m1[1]*m2[4] + m1[2]*m2[7] + m1[3]*m2[10], m1[0]*m2[2] + m1[1]*m2[5] + m1[2]*m2[8] + m1[3]*m2[11],
		m1[4]*m2[0] + m1[5]*m2[3] + m1[6]*m2[6] + m1[7]*m2[9], m1[4]*m2[1] + m1[5]*m2[4] + m1[6]*m2[7] + m1[7]*m2[10], m1[4]*m2[2] + m1[5]*m2[5] + m1[6]*m2[8] + m1[7]*m2[11],
		m1[8]*m2[0] + m1[9]*m2[3] + m1[10]*m2[6] + m1[11]*m2[9], m1[8]*m2[1] + m1[9]*m2[4] + m1[10]*m2[7] + m1[11]*m2[10], m1[8]*m2[2] + m1[9]*m2[5] + m1[10]*m2[8] + m1[11]*m2[11],
		m1[12]*m2[0] + m1[13]*m2[3] + m1[14]*m2[6] + m1[15]*m2[9], m1[12]*m2[1] + m1[13]*m2[4] + m1[14]*m2[7] + m1[15]*m2[10], m1[12]*m2[2] + m1[13]*m2[5] + m1[14]*m2[8] + m1[15]*m2[11],
	}
}

func (m1 Mat4) MulMat4(m2 Mat4) Mat4 {
	return Mat4{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8] + m1[3]*m2[12], m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9] + m1[3]*m2[13], m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10] + m1[3]*m2[14], m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11] + m1[3]*m2[15],
		m1[4]*m2[0] + m1[5]*m2[4] + m1[6]*m2[8] + m1[7]*m2[12], m1[4]*m2[1] + m1[5]*m2[5] + m1[6]*m2[9] + m1[7]*m2[13], m1[4]*m2[2] + m1[5]*m2[6] + m1[6]*m2[10] + m1[7]*m2[14], m1[4]*m2[3] + m1[5]*m2[7] + m1[6]*m2[11] + m1[7]*m2[15],
		m1[8]*m2[0] + m1[9]*m2[4] + m1[10]*m2[8] + m1[11]*m2[12], m1[8]*m2[1] + m1[9]*m2[5] + m1[10]*m2[9] + m1[11]*m2[13], m1[8]*m2[2] + m1[9]*m2[6] + m1[10]*m2[10] + m1[11]*m2[14], m1[8]*m2[3] + m1[9]*m2[7] + m1[10]*m2[11] + m1[11]*m2[15],
		m1[12]*m2[0] + m1[13]*m2[4] + m1[14]*m2[8] + m1[15]*m2[12], m1[12]*m2[1] + m1[13]*m2[5] + m1[14]*m2[9] + m1[15]*m2[13], m1[12]*m2[2] + m1[13]*m2[6] + m1[14]*m2[10] + m1[15]*m2[14], m1[12]*m2[3] + m1[13]*m2[7] + m1[14]*m2[11] + m1[15]*m2[15],
	}
}

func (m Mat4) Trace() float64 {
	return m[0] + m[5] + m[10] + m[15]
}

func (m Mat4) Det() float64 {
	return m[0]*m[5]*m[10]*m[15] - m[0]*m[5]*m[11]*m[14] - m[0]*m[6]*m[9]*m[15] + m[0]*m[6]*m[11]*m[13] +
		m[0]*m[7]*m[9]*m[14] - m[0]*m[7]*m[10]*m[13] - m[1]*m[4]*m[10]*m[15] + m[1]*m[4]*m[11]*m[14] +
		m[1]*m[6]*m[8]*m[15] - m[1]*m[6]*m[11]*m[12] - m[1]*m[7]*m[8]*m[14] + m[1]*m[7]*m[10]*m[12] +
		m[2]*m[4]*m[9]*m[15] - m[2]*m[4]*m[11]*m[13] - m[2]*m[5]*m[8]*m[15] + m[2]*m[5]*m[11]*m[12] +
		m[2]*m[7]*m[8]*m[13] - m[2]*m[7]*m[9]*m[12] - m[3]*m[4]*m[9]*m[14] + m[3]*m[4]*m[10]*m[13] +
		m[3]*m[5]*m[8]*m[14] - m[3]*m[5]*m[10]*m[12] - m[3]*m[6]*m[8]*m[13] + m[3]*m[6]*m[9]*m[12]
}

func (m Mat4) At(i, j int) float64 {
	if i >= 4 || j >= 4 {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (4x4))",
			i, j,
		)
		panic(err)
	}

	return m[j+i*4]
}

func (m Mat4) Set(i, j int, value float64) {
	if i >= 4 || j >= 4 {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (4x4))",
			i, j,
		)
		panic(err)
	}

	m[j+i*4] = value
}
