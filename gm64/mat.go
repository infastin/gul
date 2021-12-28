package gm64

import "fmt"

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

func (m1 *Mat) Add(m2 *Mat) *Mat {
	if m1.M != m2.M || m1.N != m2.N {
		err := fmt.Errorf("the first and second matrices have different dimensions (got (%dx%d) and (%dx%d))", m1.M, m1.N, m2.M, m2.N)
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
		err := fmt.Errorf("the first and second matrices have different dimensions (got (%dx%d) and (%dx%d))", m1.M, m1.N, m2.M, m2.N)
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

func (m1 *Mat) Mul(m2 *Mat) *Mat {
	if m1.N != m2.M {
		err := fmt.Errorf("trying to multiply matrices with different number of columns and rows (got (%dx%d) and (%dx%d))", m1.M, m1.N, m2.M, m2.N)
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

func (m *Mat) Get(i, j int) float64 {
	if i >= m.M || j >= m.N {
		err := fmt.Errorf("trying to get a value out of matrix bounds (got position (%d, %d) while matrix size is (%dx%d))", i, j, m.M, m.N)
		panic(err)
	}

	return m.Data[j+i*m.N]
}

func (m *Mat) Set(i, j int, value float64) {
	if i >= m.M || j >= m.N {
		err := fmt.Errorf("trying to set a value out of matrix bounds (got position (%d, %d) while matrix size is (%dx%d))", i, j, m.M, m.N)
		panic(err)
	}

	m.Data[j+i*m.N] = value
}
