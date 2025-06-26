package simd

import (
	"math"
)

func wip_emuB_F32s(vs [][]float32) float32 {
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
	var sum F32s
	const L = 64
	var i int
	for i < l-L {
		ii := i + L
		a := loadU_emuB_F32s(vs[0][i:ii])
		for j := 1; j < len(vs); j++ {
			a = mul_emuB_F32s(a, loadU_emuB_F32s(vs[j][i:ii]))
		}
		sum = add_emuB_F32s(sum, a)
		i = ii
	}
	if i < l {
		a := loadN_emuB_F32s(vs[0][i:], l-i)
		for j := 1; j < len(vs); j++ {
			a = mul_emuB_F32s(a, loadN_emuB_F32s(vs[j][i:], l-i))
		}
		sum = add_emuB_F32s(sum, a)
	}
	var tmp []float32
	if l < L {
		tmp = make([]float32, l)
		storeN_emuB_F32s(sum, &tmp, L-l)
	} else {
		tmp = make([]float32, L)
		storeU_emuB_F32s(sum, &tmp)
	}
	var x float32
	for _, y := range tmp {
		x += y
	}
	return x
}

func wip_emuB_F64s(vs [][]float64) float64 {
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
	var sum F64s
	const L = 4
	var i int
	for i < l-L {
		ii := i + L
		a := loadU_emuB_F64s(vs[0][i:ii])
		for j := 1; j < len(vs); j++ {
			a = mul_emuB_F64s(a, loadU_emuB_F64s(vs[j][i:ii]))
		}
		sum = add_emuB_F64s(sum, a)
		i = ii
	}
	if i < l {
		a := loadN_emuB_F64s(vs[0][i:], l-i)
		for j := 1; j < len(vs); j++ {
			a = mul_emuB_F64s(a, loadN_emuB_F64s(vs[j][i:], l-i))
		}
		sum = add_emuB_F64s(sum, a)
	}
	var tmp []float64
	if l < L {
		tmp = make([]float64, l)
		storeN_emuB_F64s(sum, &tmp, L-l)
	} else {
		tmp = make([]float64, L)
		storeU_emuB_F64s(sum, &tmp)
	}
	var x float64
	for _, y := range tmp {
		x += y
	}
	return x
}
