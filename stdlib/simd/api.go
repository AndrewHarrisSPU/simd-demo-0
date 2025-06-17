package simd

// universal ops
// arithmetic ops
func Add(VecT, VecT) VecT    { return VecT{} }
func Mul(VecT, VecT) VecT    { return VecT{} }
func Neg(VecT) VecT          { return VecT{} }
func Div(VecT, VecT) VecT    { return VecT{} }
func MulAdd(VecT, VecT) VecT { return VecT{} }

// logical ops
func Or(VecT, VecT) VecT  { return VecT{} }
func And(VecT, VecT) VecT { return VecT{} }

// load / store ops
func LoadU([]T) VecT         { return VecT{} }
func LoadN([]T, int) VecT    { return VecT{} }
func StoreU(VecT, *[]T)      {}
func StoreN(VecT, *[]T, int) {}

// Dispatch attributes
type T float64
type VecT struct{}

const (
	Arch string = ""
	ISA  string = ""
	N    int    = 0
	D    string = ""
)
