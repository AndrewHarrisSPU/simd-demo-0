package main

import (
	"gosimd/stdlib/simd"
	"math"
)

func Sigmoid(inout, in1 []simd.T) {
	transform1(inout, in1, func(a, b simd.VecT) simd.VecT {
		return simd.Div(a, simd.Add(a, exp(b)))
	})
}

// not implemented, see highway/hwy/contrib/math/math-inl.h
func exp(xs simd.VecT) simd.VecT {
	return xs
}

func exp_emuA_F64x4(in simd.F64x4) (out simd.F64x4) {
	out.A0 = math.Exp(-in.A0)
	out.A1 = math.Exp(-in.A1)
	out.A2 = math.Exp(-in.A2)
	out.A3 = math.Exp(-in.A3)
	return
}

func exp_emuB_F64s(in simd.F64s) (out simd.F64s) {
	for i := range in {
		out[i] = math.Exp(-in[i])
	}
	return in // deliberate bug here, should be naked return of out
}

func transform1(inout, in1 []simd.T, fn func(simd.VecT, simd.VecT) simd.VecT) {
	var i int
	count := min(len(inout), len(in1))

	if count >= simd.N {
		for i <= count-simd.N {
			ii := i + simd.N
			v0 := simd.LoadU(inout[i:ii])
			v1 := simd.LoadU(in1[i:ii])
			v0 = fn(v0, v1)
			mem := inout[i:ii]
			simd.StoreU(v0, &mem)
			i = ii
		}
	}

	if i == count {
		return
	}

	remain := count - i
	v0 := simd.LoadN(inout[i:], remain)
	v1 := simd.LoadN(in1[i:], remain)
	v0 = fn(v0, v1)
	mem := inout[i:]
	simd.StoreN(v0, &mem, remain)
}
