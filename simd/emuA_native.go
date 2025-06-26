package simd

func loadU_emuA_F32x8(vs []float32) (reg F32x8) {
	return F32x8{vs[0], vs[1], vs[2], vs[3], vs[4], vs[5], vs[6], vs[7]}
}

func loadN_emuA_F32x8(vs []float32, remain int) (reg F32x8) {
	switch remain {
	case 1:
		return F32x8{vs[0], 0, 0, 0, 0, 0, 0, 0}
	case 2:
		return F32x8{vs[0], vs[1], 0, 0, 0, 0, 0, 0}
	case 3:
		return F32x8{vs[0], vs[1], vs[2], 0, 0, 0, 0, 0}
	case 4:
		return F32x8{vs[0], vs[1], vs[2], vs[3], 0, 0, 0, 0}
	case 5:
		return F32x8{vs[0], vs[1], vs[2], vs[3], vs[4], 0, 0, 0}
	case 6:
		return F32x8{vs[0], vs[1], vs[2], vs[3], vs[4], vs[5], 0, 0}
	case 7:
		return F32x8{vs[0], vs[1], vs[2], vs[3], vs[4], vs[5], vs[6], 0}
	default:
		panic("loadN_emu_F32x8")
	}
}

func storeU_emuA_F32x8(reg F32x8, mem *[]float32) {
	(*mem)[0] = reg.A0
	(*mem)[1] = reg.A1
	(*mem)[2] = reg.A2
	(*mem)[3] = reg.A3
	(*mem)[4] = reg.A4
	(*mem)[5] = reg.A5
	(*mem)[6] = reg.A6
	(*mem)[7] = reg.A7
}

func storeN_emuA_F32x8(reg F32x8, mem *[]float32, remain int) {
	switch remain {
	case 1:
		(*mem)[0] = reg.A0
	case 2:
		(*mem)[0], (*mem)[1] = reg.A0, reg.A1
	case 3:
		(*mem)[0], (*mem)[1] = reg.A0, reg.A1
		(*mem)[2] = reg.A2
	case 4:
		(*mem)[0], (*mem)[1] = reg.A0, reg.A1
		(*mem)[2], (*mem)[3] = reg.A2, reg.A3
	case 5:
		(*mem)[0], (*mem)[1] = reg.A0, reg.A1
		(*mem)[2], (*mem)[3] = reg.A2, reg.A3
		(*mem)[4] = reg.A4
	case 6:
		(*mem)[0], (*mem)[1] = reg.A0, reg.A1
		(*mem)[2], (*mem)[3] = reg.A2, reg.A3
		(*mem)[4], (*mem)[5] = reg.A4, reg.A5
	case 7:
		(*mem)[0], (*mem)[1] = reg.A0, reg.A1
		(*mem)[2], (*mem)[3] = reg.A2, reg.A3
		(*mem)[4], (*mem)[5] = reg.A4, reg.A5
		(*mem)[6] = reg.A6
	}
}

func add_emuA_F32x8(in1, in2 F32x8) (out F32x8) {
	out.A0 = in1.A0 + in2.A0
	out.A1 = in1.A1 + in2.A1
	out.A2 = in1.A2 + in2.A2
	out.A3 = in1.A3 + in2.A3
	out.A4 = in1.A4 + in2.A4
	out.A5 = in1.A5 + in2.A5
	out.A6 = in1.A6 + in2.A6
	out.A7 = in1.A7 + in2.A7
	return
}

func mul_emuA_F32x8(in1, in2 F32x8) (out F32x8) {
	out.A0 = in1.A0 * in2.A0
	out.A1 = in1.A1 * in2.A1
	out.A2 = in1.A2 * in2.A2
	out.A3 = in1.A3 * in2.A3
	out.A4 = in1.A4 * in2.A4
	out.A5 = in1.A5 * in2.A5
	out.A6 = in1.A6 * in2.A6
	out.A7 = in1.A7 * in2.A7
	return
}

func div_emuA_F32x8(in1, in2 F32x8) (out F32x8) {
	out.A0 = in1.A0 / in2.A0
	out.A1 = in1.A1 / in2.A1
	out.A2 = in1.A2 / in2.A2
	out.A3 = in1.A3 / in2.A3
	out.A4 = in1.A4 / in2.A4
	out.A5 = in1.A5 / in2.A5
	out.A6 = in1.A6 / in2.A6
	out.A7 = in1.A7 / in2.A7
	return
}

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
		panic("loadN_emu_F64x4")
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
