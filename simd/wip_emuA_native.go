package simd

import (
	"math"
)

func wip_emuA_F32x8(vs [][]float32) float32 {
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
	var sum F32x8
	const L = 8
	var i int
	for i < l-L {
		ii := i + L
		a := loadU_emuA_F32x8(vs[0][i:ii])
		for j := 1; j < len(vs); j++ {
			a = mul_emuA_F32x8(a, loadU_emuA_F32x8(vs[j][i:ii]))
		}
		sum = add_emuA_F32x8(sum, a)
		i = ii
	}
	if i < l {
		a := loadN_emuA_F32x8(vs[0][i:], l-i)
		for j := 1; j < len(vs); j++ {
			a = mul_emuA_F32x8(a, loadN_emuA_F32x8(vs[j][i:], l-i))
		}
		sum = add_emuA_F32x8(sum, a)
	}
	var tmp []float32
	if l < L {
		tmp = make([]float32, l)
		storeN_emuA_F32x8(sum, &tmp, L-l)
	} else {
		tmp = make([]float32, L)
		storeU_emuA_F32x8(sum, &tmp)
	}
	var x float32
	for _, y := range tmp {
		x += y
	}
	return x
}

func wip_emuA_F64x4(vs [][]float64) float64 {
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
	var sum F64x4
	const L = 4
	var i int
	for i < l-L {
		ii := i + L
		a := loadU_emuA_F64x4(vs[0][i:ii])
		for j := 1; j < len(vs); j++ {
			a = mul_emuA_F64x4(a, loadU_emuA_F64x4(vs[j][i:ii]))
		}
		sum = add_emuA_F64x4(sum, a)
		i = ii
	}
	if i < l {
		a := loadN_emuA_F64x4(vs[0][i:], l-i)
		for j := 1; j < len(vs); j++ {
			a = mul_emuA_F64x4(a, loadN_emuA_F64x4(vs[j][i:], l-i))
		}
		sum = add_emuA_F64x4(sum, a)
	}
	var tmp []float64
	if l < L {
		tmp = make([]float64, l)
		storeN_emuA_F64x4(sum, &tmp, L-l)
	} else {
		tmp = make([]float64, L)
		storeU_emuA_F64x4(sum, &tmp)
	}
	var x float64
	for _, y := range tmp {
		x += y
	}
	return x
}
