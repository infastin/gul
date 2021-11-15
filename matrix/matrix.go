package matrix

type Matrix64 struct {
	M, N int
	Data []float64
}

func New64(m, n int) func(data ...float64) *Matrix64 {
	ctor := func(data ...float64) *Matrix64 {
		o := &Matrix64{
			M:    m,
			N:    n,
			Data: make([]float64, m*n),
		}

		copy(o.Data, data)
		return o
	}

	return ctor
}

func (m1 *Matrix64) Add(m2 *Matrix64) *Matrix64 {
	if m1.M != m2.M || m1.N != m2.N {
		panic("the first and second matrices have different dimensions")
	}

	o := &Matrix64{
		M:    m1.M,
		N:    m1.N,
		Data: make([]float64, m1.M*m1.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[i+j*o.M] = m1.Data[i+j*m1.M] + m2.Data[i+j*m2.M]
		}
	}

	return o
}

func (m1 *Matrix64) Sub(m2 *Matrix64) *Matrix64 {
	if m1.M != m2.M || m1.N != m2.N {
		panic("the first and second matrices have different dimensions")
	}

	o := &Matrix64{
		M:    m1.M,
		N:    m1.N,
		Data: make([]float64, m1.M*m1.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[i+j*o.M] = m1.Data[i+j*m1.M] - m2.Data[i+j*m2.M]
		}
	}

	return o
}

func (m1 *Matrix64) Mul(m2 *Matrix64) *Matrix64 {
	if m1.M != m2.N {
		panic("trying to multiply matrices with different columns and rows number")
	}

	o := &Matrix64{
		M:    m1.M,
		N:    m2.N,
		Data: make([]float64, m1.M*m2.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			for k := 0; k < m1.M; k++ {
				o.Data[i+j*o.M] += m1.Data[i+k*m1.M] * m2.Data[k+j*m2.M]
			}
		}
	}

	return o
}

func (m *Matrix64) Index(i, j int) float64 {
	return m.Data[i+j*m.M]
}

type Matrix32 struct {
	M, N int
	Data []float32
}

func New32(m, n int) func(data ...float32) *Matrix32 {
	ctor := func(data ...float32) *Matrix32 {
		o := &Matrix32{
			M:    m,
			N:    n,
			Data: make([]float32, m*n),
		}

		copy(o.Data, data)
		return o
	}

	return ctor
}

func (m1 *Matrix32) Add(m2 *Matrix32) *Matrix32 {
	if m1.M != m2.M || m1.N != m2.N {
		panic("the first and second matrices have different dimensions")
	}

	o := &Matrix32{
		M:    m1.M,
		N:    m1.N,
		Data: make([]float32, m1.M*m1.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[i+j*o.M] = m1.Data[i+j*m1.M] + m2.Data[i+j*m2.M]
		}
	}

	return o
}

func (m1 *Matrix32) Sub(m2 *Matrix32) *Matrix32 {
	if m1.M != m2.M || m1.N != m2.N {
		panic("the first and second matrices have different dimensions")
	}

	o := &Matrix32{
		M:    m1.M,
		N:    m1.N,
		Data: make([]float32, m1.M*m1.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			o.Data[i+j*o.M] = m1.Data[i+j*m1.M] - m2.Data[i+j*m2.M]
		}
	}

	return o
}

func (m1 *Matrix32) Mul(m2 *Matrix32) *Matrix32 {
	if m1.N != m2.M {
		panic("trying to multiply matrices with different columns and rows number")
	}

	o := &Matrix32{
		M:    m1.M,
		N:    m2.N,
		Data: make([]float32, m1.M*m2.N),
	}

	for i := 0; i < o.M; i++ {
		for j := 0; j < o.N; j++ {
			for k := 0; k < m1.N; k++ {
				o.Data[i+j*o.M] += m1.Data[i+k*m1.M] * m2.Data[k+j*m2.M]
			}
		}
	}

	return o
}

func (m *Matrix32) Index(i, j int) float32 {
	return m.Data[i+j*m.M]
}
