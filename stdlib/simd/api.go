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
func LoadU([]Unit) VecT         { return VecT{} }
func LoadN([]Unit, int) VecT    { return VecT{} }
func StoreU(VecT, *[]Unit)      {}
func StoreN(VecT, *[]Unit, int) {}

// Dispatch attributes
type Unit float64
type VecT struct{}

const (
	Arch  string = ""
	ISA   string = ""
	Lanes int    = 0
	D     string = ""
)
