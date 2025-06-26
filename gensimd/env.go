package gensimd

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"gosimd/stdlib/simd"

	"golang.org/x/tools/go/ast/astutil"
)

type genEnv struct {
	// simd mapping
	targets    []simd.Dispatch
	simd       map[string]string // substitution dictionary
	virtualOps map[string]bool   // maintains set of virtual API operations observed in sources

	// configures filesystem paths
	srcDir string // directory containing contents of scanned sources
	apiDir string // direcotry containing virtual and standard API
	dstDir string // directory to write outputs to

	// in-memory file tracking
	src map[string]*source          // file name -> source text
	dst map[string]*strings.Builder // source -> rewrite
	isa map[string]*strings.Builder // standard isa imports
	api map[string]*strings.Builder // virtual api

	// tracking FuncDecls
	export   map[string]string // clones of exported FuncDecl node text
	overload map[symbol]bool   // observes FuncDecls that are tagged like overloads
	local    map[string]bool   // observes FuncDecls that are not tagged and not exported
}

type source struct {
	importDecls []ast.Decl // imports
	text        string     // clone of file contents
}

func newGenEnv(targets []simd.Dispatch, apiDir, srcDir, dstDir string) *genEnv {
	env := &genEnv{
		targets:    targets,
		simd:       make(map[string]string),
		virtualOps: make(map[string]bool),

		apiDir: apiDir,
		srcDir: srcDir,
		dstDir: dstDir,

		src: make(map[string]*source),
		dst: make(map[string]*strings.Builder),
		isa: make(map[string]*strings.Builder),
		api: make(map[string]*strings.Builder),

		export:   make(map[string]string),
		overload: make(map[symbol]bool),
		local:    make(map[string]bool),
	}

	for _, target := range targets {
		isa := target.ISA()
		if _, found := env.isa[isa]; found {
			continue
		}
		env.initISA(isa)
	}

	// while the Unit and VecT symbols vary during rewriting,
	// specific units and vector types do not.
	for unit := range setUnitTypes {
		env.simd[unit] = unit
	}
	for vect := range setVecTypes {
		env.simd[vect] = vect
	}

	return env
}

func (env *genEnv) setTarget(target simd.Dispatch) {
	env.simd["ISA"] = target.ISA()
	env.simd["Unit"] = target.Unit()
	env.simd["Lanes"] = target.Lanes()
	env.simd["Tag"] = target.Tag()
	env.simd["VecT"] = target.VecT()

	// ops are dispatch-tagged
	for op := range setOps {
		env.simd[op] = string(env.tag(op))
	}
}

func (env *genEnv) tag(op string) symbol {
	s, ok := parseSymbol(op + `_` + env.simd["Tag"])
	if !ok {
		throw("tagging %s: result doesn't parse: %s", op, op+`_`+env.simd["Tag"])
	}
	return s
}

func (env *genEnv) initSource(name string, data string) {
	env.src[name] = &source{text: data}
	for isa := range env.isa {
		env.initDst(name, isa)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, parser.SkipObjectResolution)
	if err != nil {
		throw("parsing source %s: %w", name, err)
	}
	astutil.DeleteImport(fset, f, env.apiDir)

	env.scanAST(f, name)
}

func (env *genEnv) initDst(name, isa string) *strings.Builder {
	name_isa := name + `_` + isa
	checkDuplicateFile(env.dst, name_isa)
	dst := new(strings.Builder)
	dst.WriteString("package simd\n")
	env.dst[name_isa] = dst
	return dst
}

func (env *genEnv) initISA(isa string) *strings.Builder {
	checkDuplicateFile(env.isa, isa)
	dst := new(strings.Builder)
	dst.WriteString("package simd\n")
	env.isa[isa] = dst
	return dst
}

func (env *genEnv) initAPI(name string) *strings.Builder {
	checkDuplicateFile(env.api, name)
	dst := new(strings.Builder)
	env.api[name] = dst
	return dst
}

func checkDuplicateFile(m map[string]*strings.Builder, name string) {
	if _, found := m[name]; found {
		throw("duplicate file init: %s", name)
	}
}

func (env *genEnv) publishNode(dst *strings.Builder, n ast.Node) {
	if dst.Len() == 0 {
		throw("publishNode called on blank builder")
	}
	dst.WriteString("\n")
	format.Node(dst, token.NewFileSet(), n)
	dst.WriteString("\n")
}

func (env *genEnv) publishFile(dst *strings.Builder, f *ast.File) {
	if dst.Len() == 0 {
		throw("publishFile called on blank builder")
	}
	for _, decl := range f.Decls {
		env.publishNode(dst, decl)
	}
}

// validated to be unexported and contain a valid tag
type symbol string

func parseSymbol(s string) (symbol, bool) {
	i, ii := symbol(s).index()
	if i < 0 {
		return "", false
	}

	isa, vect := s[i+1:ii], s[ii+1:]
	if !supportedTags[isa][vect] {
		throw("parsed unsupported tag %s_%s", isa, vect)
	}
	return symbol(unexport(s)), true
}

func (s symbol) index() (int, int) {
	ii := strings.LastIndex(string(s), "_")
	if ii < 0 {
		return -1, -1
	}
	i := strings.LastIndex(string(s)[:ii], "_")
	if i < 0 {
		return -1, -1
	}
	return i, ii
}

func (s symbol) op() string {
	i, _ := s.index()
	return string(s[:i])
}

func (s symbol) isa() string {
	i, ii := s.index()
	return string(s[i+1 : ii])
}

func (s symbol) tag() string {
	i, _ := s.index()
	return string(s)[i+1:]
}

func (s symbol) exportOp() string {
	return export(s.op())
}

func export(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func unexport(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
