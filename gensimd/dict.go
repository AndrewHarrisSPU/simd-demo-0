package gensimd

import (
	"unicode"
	"unicode/utf8"

	"gosimd/stdlib/simd"
)

// list of symbols in stdlib/simd package, but not in Dispatch
var simdTypes = []string{
	"F64",
}

var simdVecTypes = []string{
	"F64x4", "F64s",
}

var simdISAs = []string{
	"emuA", "emuB",
}

var simdOps = []string{
	"Add", "Mul", "Div", "LoadU", "LoadN", "StoreU", "StoreN",
}

func makeSymbolMap(d simd.Dispatch) map[string]string {
	m := make(map[string]string)

	// merge Dispatch symbols
	m["ISA"] = d.ISA()
	m["T"] = d.T()
	m["N"] = d.N()
	m["tag"] = d.Tag()
	m["VecT"] = d.VecT()
	m["Arch"] = d.Arch()

	// preserve simd VecTypes
	for _, vect := range simdVecTypes {
		m[vect] = vect
	}

	// merge simd types
	for _, t := range simdTypes {
		m[t] = unexport(t)
	}

	// merge simd ops symbols
	for _, op := range simdOps {
		m[op] = unexport(op) + `_` + m["tag"]
	}

	return m
}

func unexport(symbol string) string {
	if len(symbol) == 0 {
		panic("empty symbol")
	}

	r, n := utf8.DecodeRuneInString(symbol)
	if !unicode.IsUpper(r) {
		return symbol
	}

	return string(unicode.ToLower(r)) + symbol[n:]
}
