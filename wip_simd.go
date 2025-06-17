package main

import (
	"gosimd/stdlib/simd"
	"math"
)

func Wip(vs [][]simd.T) simd.T {
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
	const L = simd.N
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

	var tmp []simd.T

	if l < L {
		tmp = make([]simd.T, l)
		simd.StoreN(sum, &tmp, L-l)
	} else {
		tmp = make([]simd.T, L)
		simd.StoreU(sum, &tmp)
	}

	var x simd.T
	for _, y := range tmp {
		x += y
	}
	return x
}
