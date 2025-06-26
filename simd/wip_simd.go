//go:build ignore

package simd

import (
	"math"
	"stdlib/simd"
)

func Wip(vs [][]simd.Unit) simd.Unit {
	if len(vs) == 0 {
		return 0
	}

	l := math.MaxInt
	for _, v := range vs {
		l = min(l, len(v))
	}
	if l == 0 {
		return 0
	}

	var sum simd.VecT
	const L = simd.Lanes
	var i int

	for i < l-L {
		ii := i + L
		a := simd.LoadU(vs[0][i:ii])
		for j := 1; j < len(vs); j++ {
			a = simd.Mul(a, simd.LoadU(vs[j][i:ii]))
		}
		sum = simd.Add(sum, a)
		i = ii
	}

	if i < l {
		a := simd.LoadN(vs[0][i:], l-i)
		for j := 1; j < len(vs); j++ {
			a = simd.Mul(a, simd.LoadN(vs[j][i:], l-i))
		}
		sum = simd.Add(sum, a)
	}

	var tmp []simd.Unit

	if l < L {
		tmp = make([]simd.Unit, l)
		simd.StoreN(sum, &tmp, L-l)
	} else {
		tmp = make([]simd.Unit, L)
		simd.StoreU(sum, &tmp)
	}

	var x simd.Unit
	for _, y := range tmp {
		x += y
	}
	return x
}
