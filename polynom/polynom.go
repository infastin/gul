package polynom

import (
	"container/list"
	"fmt"
	"math"
)

type Monomial struct {
	Coef   float64
	Degree int
}

type Polynomial struct {
	monoms *list.List
	degree int
}

func New(monoms ...Monomial) *Polynomial {
	p := &Polynomial{
		monoms: list.New(),
	}

	for _, m := range monoms {
		if m.Coef == 0 {
			continue
		}

		p.addMonom(m)
	}

	return p
}

func (p *Polynomial) addMonom(monom Monomial) {
	for it := p.monoms.Front(); it != nil; it = it.Next() {
		m := it.Value.(Monomial)

		switch {
		case monom.Degree > m.Degree:
			continue
		case monom.Degree < m.Degree:
			p.monoms.InsertBefore(monom, it)
		case monom.Degree == m.Degree:
			coef := m.Coef + monom.Coef
			if coef == 0 {
				p.monoms.Remove(it)
			} else {
				it.Value = Monomial{
					Coef:   coef,
					Degree: monom.Degree,
				}
			}

			p.degree = p.monoms.Back().Value.(Monomial).Degree
		}

		return
	}

	p.monoms.PushBack(monom)
	p.degree = monom.Degree
}

func (p *Polynomial) subMonom(monom Monomial) {
	p.addMonom(Monomial{
		Coef:   -monom.Coef,
		Degree: monom.Degree,
	})
}

func (p *Polynomial) AddMonom(monom Monomial) *Polynomial {
	o := New()

	if p.monoms.Len() == 0 && monom.Coef == 0 {
		return o
	}

	for it := p.monoms.Front(); it != nil; it = it.Next() {
		m := it.Value.(Monomial)

		coef := m.Coef + monom.Coef
		if coef == 0 {
			continue
		}

		o.monoms.PushBack(Monomial{
			Coef:   coef,
			Degree: m.Degree,
		})
	}

	if m := o.monoms.Back(); m != nil {
		o.degree = m.Value.(Monomial).Degree
	}

	return o
}

func (p *Polynomial) SubMonom(monom Monomial) *Polynomial {
	return p.AddMonom(Monomial{
		Coef:   -monom.Coef,
		Degree: monom.Degree,
	})
}

func (p *Polynomial) MulMonom(monom Monomial) *Polynomial {
	o := New()

	if monom.Coef == 0 || p.monoms.Len() == 0 {
		return o
	}

	for it := p.monoms.Front(); it != nil; it = it.Next() {
		m := it.Value.(Monomial)
		degree := m.Degree + monom.Degree

		if o.monoms.Len() == 0 || degree > o.degree {
			o.degree = degree
		}

		o.monoms.PushBack(Monomial{
			Coef:   m.Coef * monom.Coef,
			Degree: degree,
		})
	}

	o.degree = o.monoms.Back().Value.(Monomial).Degree
	return o
}

func (p1 *Polynomial) Add(p2 *Polynomial) *Polynomial {
	o := New()

	if p1.monoms.Len() == 0 && p2.monoms.Len() == 0 {
		return o
	}

	it1, it2 := p1.monoms.Front(), p2.monoms.Front()
loop:
	for it1 != nil || it2 != nil {
		var m Monomial

		switch {
		case it1 != nil && it2 != nil:
			m1 := it1.Value.(Monomial)
			m2 := it2.Value.(Monomial)

			switch {
			case m1.Degree == m2.Degree:
				it1 = it1.Next()
				it2 = it2.Next()

				coef := m1.Coef + m2.Coef
				if coef == 0 {
					continue loop
				}

				m = Monomial{
					Coef:   coef,
					Degree: m1.Degree,
				}
			case m1.Degree < m2.Degree:
				it1 = it1.Next()
				m = m1
			default:
				it2 = it2.Next()
				m = m2
			}
		case it1 != nil:
			m = it1.Value.(Monomial)
			it1 = it1.Next()
		case it2 != nil:
			m = it2.Value.(Monomial)
			it2 = it2.Next()
		}

		o.monoms.PushBack(m)
	}

	if m := o.monoms.Back(); m != nil {
		o.degree = m.Value.(Monomial).Degree
	}

	return o
}

func (p1 *Polynomial) Sub(p2 *Polynomial) *Polynomial {
	o := New()

	if p1.monoms.Len() == 0 && p2.monoms.Len() == 0 {
		return o
	}

	it1, it2 := p1.monoms.Front(), p2.monoms.Front()
loop:
	for it1 != nil || it2 != nil {
		var m Monomial

		switch {
		case it1 != nil && it2 != nil:
			m1 := it1.Value.(Monomial)
			m2 := it2.Value.(Monomial)

			switch {
			case m1.Degree == m2.Degree:
				it1 = it1.Next()
				it2 = it2.Next()

				coef := m1.Coef - m2.Coef
				if coef == 0 {
					continue loop
				}

				m = Monomial{
					Coef:   coef,
					Degree: m1.Degree,
				}
			case m1.Degree < m2.Degree:
				it1 = it1.Next()
				m = m1
			default:
				it2 = it2.Next()
				m = m2
				m.Coef = -m.Coef
			}
		case it1 != nil:
			m = it1.Value.(Monomial)
			it1 = it1.Next()
		case it2 != nil:
			m = it2.Value.(Monomial)
			m.Coef = -m.Coef
			it2 = it2.Next()
		}

		o.monoms.PushBack(m)
	}

	if m := o.monoms.Back(); m != nil {
		o.degree = m.Value.(Monomial).Degree
	}

	return o
}

func (p1 *Polynomial) Mul(p2 *Polynomial) *Polynomial {
	o := New()

	for it1 := p1.monoms.Front(); it1 != nil; it1 = it1.Next() {
		m1 := it1.Value.(Monomial)
		for it2 := p2.monoms.Front(); it2 != nil; it2 = it2.Next() {
			m2 := it2.Value.(Monomial)
			o.addMonom(Monomial{
				Coef:   m1.Coef * m2.Coef,
				Degree: m1.Degree + m2.Degree,
			})
		}
	}

	return o
}

func (p1 *Polynomial) DivMod(p2 *Polynomial) (q, r *Polynomial) {
	q = New()

	r = p1
	d := p2

	for !r.IsZero() && r.degree >= d.degree {
		m1 := r.monoms.Back().Value.(Monomial)
		m2 := d.monoms.Back().Value.(Monomial)

		t := Monomial{
			Coef:   m1.Coef / m2.Coef,
			Degree: m1.Degree - m2.Degree,
		}

		q.addMonom(t)
		r = r.Sub(d.MulMonom(t))
	}

	return
}

func (p1 *Polynomial) Div(p2 *Polynomial) *Polynomial {
	q, _ := p1.DivMod(p2)
	return q
}

func (p1 *Polynomial) Mod(p2 *Polynomial) *Polynomial {
	_, r := p1.DivMod(p2)
	return r
}

func (p1 *Polynomial) Euclidean(p2 *Polynomial) *Polynomial {
	a, b := p1, p2
	var gcd *Polynomial

	for !b.IsZero() {
		gcd = b
		_, r := a.DivMod(b)
		a = b
		b = r
	}

	if gcd == nil {
		gcd = New()
	}

	return gcd
}

func (p *Polynomial) IsZero() bool {
	return p.monoms.Len() == 0
}

func (p *Polynomial) Degree() int {
	return p.degree
}

func (p *Polynomial) Calc(x float64) float64 {
	var fx float64
	for it := p.monoms.Front(); it != nil; it = it.Next() {
		m := it.Value.(Monomial)
		fx += m.Coef * math.Pow(x, float64(m.Degree))
	}
	return fx
}

func (p *Polynomial) String() string {
	if p.monoms.Len() == 0 {
		return "0"
	}

	var res string
	for it := p.monoms.Back(); it != nil; it = it.Prev() {
		m := it.Value.(Monomial)

		if it != p.monoms.Back() {
			if m.Coef > 0 {
				res += " + "
			} else {
				res += " - "
			}
		} else if m.Coef < 0 {
			res += "-"
		}

		if math.Abs(m.Coef) != 1 || m.Degree == 0 {
			res += fmt.Sprintf("%.2f", math.Abs(m.Coef))
		}

		if m.Degree != 0 {
			if m.Degree == 1 {
				res += "x"
			} else {
				res += fmt.Sprintf("x^%d", m.Degree)
			}
		}
	}

	return res
}
