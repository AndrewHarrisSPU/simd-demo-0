package simd

// 32s

func loadU_emuB_F32s(vs []float32) (reg F32s) {
	for i := range vs {
		reg[i] = vs[i]
	}
	return
}

func loadN_emuB_F32s(vs []float32, remain int) (reg F32s) {
	for i := range remain {
		reg[i] = vs[i]
	}
	return
}

func storeU_emuB_F32s(reg F32s, mem *[]float32) {
	for i := range reg {
		(*mem)[i] = reg[i]
	}
}

func storeN_emuB_F32s(reg F32s, mem *[]float32, remain int) {
	for i := range remain {
		(*mem)[i] = reg[i]
	}
}

func add_emuB_F32s(in1, in2 F32s) (out F32s) {
	for i := range in1 {
		out[i] = in1[i] + in2[i]
	}
	return
}

func sub_emuB_F32s(in1, in2 F32s) (out F32s) {
	for i := range in1 {
		out[i] = in1[i] - in2[i]
	}
	return
}

func mul_emuB_F32s(in1, in2 F32s) (out F32s) {
	for i := range in1 {
		out[i] = in1[i] * in2[i]
	}
	return
}

func div_emuB_F32s(in1, in2 F32s) (out F32s) {
	for i := range in1 {
		out[i] = in1[i] / in2[i]
	}
	return
}

// 64s

func loadU_emuB_F64s(vs []float64) (reg F64s) {
	for i := range vs {
		reg[i] = vs[i]
	}
	return
}

func loadN_emuB_F64s(vs []float64, remain int) (reg F64s) {
	for i := range remain {
		reg[i] = vs[i]
	}
	return
}

func storeU_emuB_F64s(reg F64s, mem *[]float64) {
	for i := range reg {
		(*mem)[i] = reg[i]
	}
}

func storeN_emuB_F64s(reg F64s, mem *[]float64, remain int) {
	for i := range remain {
		(*mem)[i] = reg[i]
	}
}

func add_emuB_F64s(in1, in2 F64s) (out F64s) {
	for i := range in1 {
		out[i] = in1[i] + in2[i]
	}
	return
}

func mul_emuB_F64s(in1, in2 F64s) (out F64s) {
	for i := range in1 {
		out[i] = in1[i] * in2[i]
	}
	return
}

func div_emuB_F64s(in1, in2 F64s) (out F64s) {
	for i := range in1 {
		out[i] = in1[i] / in2[i]
	}
	return
}
