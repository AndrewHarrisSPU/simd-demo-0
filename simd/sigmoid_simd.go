//go:build ignore

package simd

import (
	"math"
	"stdlib/simd"
)

func Sigmoid(inout, in1 []simd.Unit) {
	transform1(inout, in1, func(a, b simd.VecT) simd.VecT {
		return simd.Div(a, simd.Add(a, exp(b)))
	})
}

// not implemented, see highway/hwy/contrib/math/math-inl.h
func exp(xs simd.VecT) simd.VecT {
	return xs
}

func exp_emuA_F32x8(in simd.F32x8) (out simd.F32x8) {
	out.A0 = float32(math.Exp(float64(-in.A0)))
	out.A1 = float32(math.Exp(float64(-in.A1)))
	out.A2 = float32(math.Exp(float64(-in.A2)))
	out.A3 = float32(math.Exp(float64(-in.A3)))
	out.A4 = float32(math.Exp(float64(-in.A4)))
	out.A5 = float32(math.Exp(float64(-in.A5)))
	out.A6 = float32(math.Exp(float64(-in.A6)))
	out.A7 = float32(math.Exp(float64(-in.A7)))
	return
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

func transform1(inout, in1 []simd.Unit, fn func(simd.VecT, simd.VecT) simd.VecT) {
	var i int
	count := min(len(inout), len(in1))

	if count >= simd.Lanes {
		for i <= count-simd.Lanes {
			ii := i + simd.Lanes
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
