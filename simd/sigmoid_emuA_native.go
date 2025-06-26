package simd

import (
	"math"
)

func sigmoid_emuA_F32x8(inout, in1 []float32) {
	transform1_emuA_F32x8(inout, in1, func(a, b F32x8) F32x8 {
		return div_emuA_F32x8(a, add_emuA_F32x8(a, exp_emuA_F32x8(b)))
	})
}

func exp_emuA_F32x8(xs F32x8) F32x8 {
	return xs
}

func transform1_emuA_F32x8(inout, in1 []float32, fn func(F32x8, F32x8) F32x8) {
	var i int
	count := min(len(inout), len(in1))
	if count >= 8 {
		for i <= count-8 {
			ii := i + 8
			v0 := loadU_emuA_F32x8(inout[i:ii])
			v1 := loadU_emuA_F32x8(in1[i:ii])
			v0 = fn(v0, v1)
			mem := inout[i:ii]
			storeU_emuA_F32x8(v0, &mem)
			i = ii
		}
	}
	if i == count {
		return
	}
	remain := count - i
	v0 := loadN_emuA_F32x8(inout[i:], remain)
	v1 := loadN_emuA_F32x8(in1[i:], remain)
	v0 = fn(v0, v1)
	mem := inout[i:]
	storeN_emuA_F32x8(v0, &mem, remain)
}

func exp_emuA_F64x4(in F64x4) (out F64x4) {
	out.A0 = math.Exp(-in.A0)
	out.A1 = math.Exp(-in.A1)
	out.A2 = math.Exp(-in.A2)
	out.A3 = math.Exp(-in.A3)
	return
}

func sigmoid_emuA_F64x4(inout, in1 []float64) {
	transform1_emuA_F64x4(inout, in1, func(a, b F64x4) F64x4 {
		return div_emuA_F64x4(a, add_emuA_F64x4(a, exp_emuA_F64x4(b)))
	})
}

func transform1_emuA_F64x4(inout, in1 []float64, fn func(F64x4, F64x4) F64x4) {
	var i int
	count := min(len(inout), len(in1))
	if count >= 4 {
		for i <= count-4 {
			ii := i + 4
			v0 := loadU_emuA_F64x4(inout[i:ii])
			v1 := loadU_emuA_F64x4(in1[i:ii])
			v0 = fn(v0, v1)
			mem := inout[i:ii]
			storeU_emuA_F64x4(v0, &mem)
			i = ii
		}
	}
	if i == count {
		return
	}
	remain := count - i
	v0 := loadN_emuA_F64x4(inout[i:], remain)
	v1 := loadN_emuA_F64x4(in1[i:], remain)
	v0 = fn(v0, v1)
	mem := inout[i:]
	storeN_emuA_F64x4(v0, &mem, remain)
}
