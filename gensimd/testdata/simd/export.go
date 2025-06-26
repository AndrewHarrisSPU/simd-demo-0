package simd

func Foo_F64(in float64) (out float64) {
	switch RuntimeD.ISA() {
	case "emuA":
		return foo_emuA_F64x4(in)
	}
	panic("invalid dispatch: " + RuntimeD.ISA())
}
