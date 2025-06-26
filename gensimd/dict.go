package gensimd

// set of symbols in stdlib/simd package, but not in Dispatch
var setUnitTypes = map[string]bool{
	"F64": true,
	"F32": true,
}

var setVecTypes = map[string]bool{
	"F64x4": true,
	"F32x8": true,
	"F32s":  true,
	"F64s":  true,
}

var listISAs = []string{
	"emuA", "emuB",
}

var setOps = map[string]bool{
	"Add":    true,
	"Mul":    true,
	"Div":    true,
	"Neg":    true,
	"LoadU":  true,
	"LoadN":  true,
	"StoreU": true,
	"StoreN": true,
}

var dictNativeToUnit = map[string]string{
	"float32": "F32",
	"float64": "F64",
}

var dictUnitsToNative = map[string]string{
	"F32": "float32",
	"F64": "float64",
}

var dictVecTsToNative = map[string]string{
	"F32x8": "float32",
	"F64x4": "float64",
	"F32s":  "float32",
	"F64s":  "float64",
}

var dictVecTsToUnit = map[string]string{
	"F32x8": "F32",
	"F64x4": "F64",
	"F32s":  "F32",
	"F64s":  "F64",
}

var dictISAtoArch = map[string]string{
	"emuA": "native",
	"emuB": "native",
}

var supportedTags tagSet = tagSet{
	"emuA": {
		"F32x8": true,
		"F64x4": true,
	},
	"emuB": {
		"F32s": true,
		"F64s": true,
	},
}
