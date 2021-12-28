package gm64

import (
	"fmt"
	"math"
	"strings"
	"text/tabwriter"
)

type Mat struct {
	M, N int
	Data []float64
}

func NewMat(m, n int) func(data ...float64) *Mat {
	if m <= 0 || n <= 0 {
		err := fmt.Errorf("the m and n parameters must be positive (got %d and %d)", m, n)
		panic(err)
	}

	ctor := func(data ...float64) *Mat {
		if len(data) > m*n {
			err := fmt.Errorf("the number of input values must not be greater than m * n (%d * %d)", m, n)
			panic(err)
		}

		o := &Mat{
			M:    m,
			N:    n,
			Data: make([]float64, m*n),
		}

		copy(o.Data, data)
		return o
	}

	return ctor
}

func (m *Mat) Copy() *Mat {
	cp := &Mat{
		M:    m.M,
		N:    m.N,
		Data: make([]float64, m.M*m.N),
	}

	copy(cp.Data, m.Data)

	return cp
}

func (m1 *Mat) Add(m2 *Mat) *Mat {
	if m1.M != m2.M || m1.N != m2.N {
		err := fmt.Errorf(
			"the first and second matrices have different dimensions (got (%dx%d) and (%dx%d))",
			m1.M, m1.N, m2.M, m2.N,
		)
		panic(err)
	}

	o := &Mat{
		M:    m1.M,
		N:    m1.N,
		Data: make([]float64, m1.M*m1.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[j+i*o.N] = m1.Data[j+i*m1.N] + m2.Data[j+i*m2.N]
		}
	}

	return o
}

func (m1 *Mat) Sub(m2 *Mat) *Mat {
	if m1.M != m2.M || m1.N != m2.N {
		err := fmt.Errorf(
			"the first and second matrices have different dimensions (got (%dx%d) and (%dx%d))",
			m1.M, m1.N, m2.M, m2.N,
		)
		panic(err)
	}

	o := &Mat{
		M:    m1.M,
		N:    m1.N,
		Data: make([]float64, m1.M*m1.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[j+i*o.N] = m1.Data[j+i*m1.N] - m2.Data[j+i*m2.N]
		}
	}

	return o
}

func (m *Mat) Mul(c float64) *Mat {
	o := &Mat{
		M:    m.M,
		N:    m.N,
		Data: make([]float64, m.M*m.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[j+i*o.N] = m.Data[j+i*o.N] * c
		}
	}

	return o
}

func (m1 *Mat) MulMat(m2 *Mat) *Mat {
	if m1.N != m2.M {
		err := fmt.Errorf(
			"trying to multiply matrices with different number of columns and rows (got (%dx%d) and (%dx%d))",
			m1.M, m1.N, m2.M, m2.N,
		)
		panic(err)
	}

	o := &Mat{
		M:    m1.M,
		N:    m2.N,
		Data: make([]float64, m1.M*m2.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			for k := 0; k < m1.N; k++ {
				o.Data[j+i*o.N] += m1.Data[k+i*m1.N] * m2.Data[j+k*m2.N]
			}
		}
	}

	return o
}

func (m *Mat) At(i, j int) float64 {
	if i < 0 || j < 0 {
		err := fmt.Errorf("the i and j parameters must be non-negative (got %d and %d)", i, j)
		panic(err)
	}

	if i >= m.M || j >= m.N {
		err := fmt.Errorf(
			"trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (%dx%d))",
			i, j, m.M, m.N,
		)
		panic(err)
	}

	return m.Data[j+i*m.N]
}

func (m *Mat) Set(i, j int, value float64) {
	if i < 0 || j < 0 {
		err := fmt.Errorf("the i and j parameters must be non-negative (got %d and %d)", i, j)
		panic(err)
	}

	if i >= m.M || j >= m.N {
		err := fmt.Errorf(
			"trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (%dx%d))",
			i, j, m.M, m.N,
		)
		panic(err)
	}

	m.Data[j+i*m.N] = value
}

func (m *Mat) Det() float64 {
	if m.M != m.N {
		err := fmt.Errorf(
			"trying to get a determinant of a non-square matrix (matrix size is (%dx%d))",
			m.M, m.N,
		)
		panic(err)
	}

	switch m.M {
	case 1:
		return m.Data[0]
	case 2:
		return m.Data[0]*m.Data[3] - m.Data[1]*m.Data[2]
	case 3:
		return m.Data[0]*m.Data[4]*m.Data[8] - m.Data[0]*m.Data[5]*m.Data[7] - m.Data[1]*m.Data[3]*m.Data[8] +
			m.Data[1]*m.Data[5]*m.Data[6] + m.Data[2]*m.Data[3]*m.Data[7] - m.Data[2]*m.Data[4]*m.Data[6]
	case 4:
		return m.Data[0]*m.Data[5]*m.Data[10]*m.Data[15] - m.Data[0]*m.Data[5]*m.Data[11]*m.Data[14] -
			m.Data[0]*m.Data[6]*m.Data[9]*m.Data[15] + m.Data[0]*m.Data[6]*m.Data[11]*m.Data[13] +
			m.Data[0]*m.Data[7]*m.Data[9]*m.Data[14] - m.Data[0]*m.Data[7]*m.Data[10]*m.Data[13] -
			m.Data[1]*m.Data[4]*m.Data[10]*m.Data[15] + m.Data[1]*m.Data[4]*m.Data[11]*m.Data[14] +
			m.Data[1]*m.Data[6]*m.Data[8]*m.Data[15] - m.Data[1]*m.Data[6]*m.Data[11]*m.Data[12] -
			m.Data[1]*m.Data[7]*m.Data[8]*m.Data[14] + m.Data[1]*m.Data[7]*m.Data[10]*m.Data[12] +
			m.Data[2]*m.Data[4]*m.Data[9]*m.Data[15] - m.Data[2]*m.Data[4]*m.Data[11]*m.Data[13] -
			m.Data[2]*m.Data[5]*m.Data[8]*m.Data[15] + m.Data[2]*m.Data[5]*m.Data[11]*m.Data[12] +
			m.Data[2]*m.Data[7]*m.Data[8]*m.Data[13] - m.Data[2]*m.Data[7]*m.Data[9]*m.Data[12] -
			m.Data[3]*m.Data[4]*m.Data[9]*m.Data[14] + m.Data[3]*m.Data[4]*m.Data[10]*m.Data[13] +
			m.Data[3]*m.Data[5]*m.Data[8]*m.Data[14] - m.Data[3]*m.Data[5]*m.Data[10]*m.Data[12] -
			m.Data[3]*m.Data[6]*m.Data[8]*m.Data[13] + m.Data[3]*m.Data[6]*m.Data[9]*m.Data[12]
	default:
		const EPS = 1e-12

		cp := m.Copy()
		det := float64(1)

		for i := 0; i < cp.M; i++ {
			k := i

			for j := i + 1; j < cp.M; j++ {
				a1 := math.Abs(cp.Data[i+j*cp.N])
				a2 := math.Abs(cp.Data[i+k*cp.N])
				if a1 > a2 {
					k = j
				}
			}

			if math.Abs(cp.Data[i+k*cp.N]) < EPS {
				return 0
			}

			if i != k {
				for j := 0; j < cp.N; j++ {
					tmp := cp.Data[j+i*cp.N]
					cp.Data[j+i*cp.N] = cp.Data[j+k*cp.N]
					cp.Data[j+k*cp.N] = tmp
				}

				det = -det
			}
			det *= cp.Data[i+i*cp.N]

			for j := i + 1; j < cp.M; j++ {
				cp.Data[j+i*cp.N] /= cp.Data[i+i*cp.N]
			}
			cp.Data[i+i*cp.N] = 1

			for j := i + 1; j < cp.M; j++ {
				if math.Abs(cp.Data[i+j*cp.N]) < EPS {
					continue
				}

				tmp := cp.Data[i+j*cp.N]
				for l := i; l < cp.M; l++ {
					cp.Data[l+j*cp.N] -= tmp * cp.Data[l+i*cp.N]
				}
			}
		}

		return det
	}
}

func (m *Mat) Trace() float64 {
	if m.M != m.N {
		err := fmt.Errorf(
			"trying to get a trace of a non-square matrix (matrix size is (%dx%d))",
			m.M, m.N,
		)
		panic(err)
	}

	trace := float64(0)
	for i := 0; i < m.M; i++ {
		trace += m.Data[i+i*m.N]
	}

	return trace
}

func (m *Mat) Transpose() *Mat {
	o := &Mat{
		M:    m.N,
		N:    m.M,
		Data: make([]float64, m.N*m.M),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[i+j*o.N] = m.Data[j+i*m.N]
		}
	}

	return o
}

func (m *Mat) String() string {
	sb := &strings.Builder{}
	w := tabwriter.NewWriter(sb, 4, 4, 1, ' ', 0)

	for i := 0; i < m.M; i++ {
		for j := 0; j < m.N; j++ {
			fmt.Fprintf(w, "%f\t", m.Data[j+i*m.N])
		}

		if i != m.M-1 {
			fmt.Fprintf(w, "\n")
		}
	}

	w.Flush()

	return sb.String()
}
