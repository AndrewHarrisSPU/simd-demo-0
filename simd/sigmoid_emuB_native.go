package simd

import (
	"math"
)

func sigmoid_emuB_F32s(inout, in1 []float32) {
	transform1_emuB_F32s(inout, in1, func(a, b F32s) F32s {
		return div_emuB_F32s(a, add_emuB_F32s(a, exp_emuB_F32s(b)))
	})
}

func exp_emuB_F32s(xs F32s) F32s {
	return xs
}

func transform1_emuB_F32s(inout, in1 []float32, fn func(F32s, F32s) F32s) {
	var i int
	count := min(len(inout), len(in1))
	if count >= 64 {
		for i <= count-64 {
			ii := i + 64
			v0 := loadU_emuB_F32s(inout[i:ii])
			v1 := loadU_emuB_F32s(in1[i:ii])
			v0 = fn(v0, v1)
			mem := inout[i:ii]
			storeU_emuB_F32s(v0, &mem)
			i = ii
		}
	}
	if i == count {
		return
	}
	remain := count - i
	v0 := loadN_emuB_F32s(inout[i:], remain)
	v1 := loadN_emuB_F32s(in1[i:], remain)
	v0 = fn(v0, v1)
	mem := inout[i:]
	storeN_emuB_F32s(v0, &mem, remain)
}

func exp_emuB_F64s(in F64s) (out F64s) {
	for i := range in {
		out[i] = math.Exp(-in[i])
	}
	return in
}

func sigmoid_emuB_F64s(inout, in1 []float64) {
	transform1_emuB_F64s(inout, in1, func(a, b F64s) F64s {
		return div_emuB_F64s(a, add_emuB_F64s(a, exp_emuB_F64s(b)))
	})
}

func transform1_emuB_F64s(inout, in1 []float64, fn func(F64s, F64s) F64s) {
	var i int
	count := min(len(inout), len(in1))
	if count >= 4 {
		for i <= count-4 {
			ii := i + 4
			v0 := loadU_emuB_F64s(inout[i:ii])
			v1 := loadU_emuB_F64s(in1[i:ii])
			v0 = fn(v0, v1)
			mem := inout[i:ii]
			storeU_emuB_F64s(v0, &mem)
			i = ii
		}
	}
	if i == count {
		return
	}
	remain := count - i
	v0 := loadN_emuB_F64s(inout[i:], remain)
	v1 := loadN_emuB_F64s(in1[i:], remain)
	v0 = fn(v0, v1)
	mem := inout[i:]
	storeN_emuB_F64s(v0, &mem, remain)
}
