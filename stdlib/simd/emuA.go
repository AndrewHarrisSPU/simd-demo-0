package simd

func loadU_emuA_F64x4(vs []float64) (reg F64x4) {
	return F64x4{vs[0], vs[1], vs[2], vs[3]}
}

func loadN_emuA_F64x4(vs []float64, remain int) (reg F64x4) {
	switch remain {
	case 1:
		return F64x4{vs[0], 0, 0, 0}
	case 2:
		return F64x4{vs[0], vs[1], 0, 0}
	case 3:
		return F64x4{vs[0], vs[1], vs[2], 0}
	default:
		panic("loadN_emu_F64_s")
	}
}

func storeU_emuA_F64x4(reg F64x4, mem *[]float64) {
	(*mem)[0] = reg.A0
	(*mem)[1] = reg.A1
	(*mem)[2] = reg.A2
	(*mem)[3] = reg.A3
}

func storeN_emuA_F64x4(reg F64x4, mem *[]float64, remain int) {
	switch remain {
	case 1:
		(*mem)[0] = reg.A0
	case 2:
		(*mem)[0] = reg.A0
		(*mem)[1] = reg.A1
	case 3:
		(*mem)[0] = reg.A0
		(*mem)[1] = reg.A1
		(*mem)[2] = reg.A2
	default:
		panic("storeN_emu_F64_s")
	}
}

func add_emuA_F64x4(in1, in2 F64x4) (out F64x4) {
	out.A0 = in1.A0 + in2.A0
	out.A1 = in1.A1 + in2.A1
	out.A2 = in1.A2 + in2.A2
	out.A3 = in1.A3 + in2.A3
	return
}

func mul_emuA_F64x4(in1, in2 F64x4) (out F64x4) {
	out.A0 = in1.A0 * in2.A0
	out.A1 = in1.A1 * in2.A1
	out.A2 = in1.A2 * in2.A2
	out.A3 = in1.A3 * in2.A3
	return
}

func div_emuA_F64x4(in1, in2 F64x4) (out F64x4) {
	out.A0 = in1.A0 / in2.A0
	out.A1 = in1.A1 / in2.A1
	out.A2 = in1.A2 / in2.A2
	out.A3 = in1.A3 / in2.A3
	return
}
