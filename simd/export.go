package simd

func Sigmoid_F32(inout, in1 []float32) {
	switch RuntimeD.ISA() {
	case "emuA":
		sigmoid_emuA_F32x8(inout, in1)
	case "emuB":
		sigmoid_emuB_F32s(inout, in1)
	}
}

func Wip_F32(vs [][]float32) float32 {
	switch RuntimeD.ISA() {
	case "emuA":
		return wip_emuA_F32x8(vs)
	case "emuB":
		return wip_emuB_F32s(vs)
	}
	panic("invalid dispatch: " + RuntimeD.ISA())
}

func Sigmoid_F64(inout, in1 []float64) {
	switch RuntimeD.ISA() {
	case "emuA":
		sigmoid_emuA_F64x4(inout, in1)
	case "emuB":
		sigmoid_emuB_F64s(inout, in1)
	}
}

func Wip_F64(vs [][]float64) float64 {
	switch RuntimeD.ISA() {
	case "emuA":
		return wip_emuA_F64x4(vs)
	case "emuB":
		return wip_emuB_F64s(vs)
	}
	panic("invalid dispatch: " + RuntimeD.ISA())
}
