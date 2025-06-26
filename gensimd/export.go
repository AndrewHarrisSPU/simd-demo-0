package gensimd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"iter"
	"maps"
	"strings"

	"gosimd/stdlib/simd"

	"golang.org/x/tools/go/ast/astutil"
)

func (env *genEnv) generateExportFile() {
	text := new(strings.Builder)
	text.WriteString("package simd\n")

	tags := newTagSet(env.targets)
	for unit := range tags.intersectUnits() {
		for name, signature := range env.export {
			fn := genExportFunc(name, unit, signature, tags)
			env.publishNode(text, fn)
		}
	}

	env.initAPI("export")
	env.api["export"] = text
}

type exportFunc struct {
	*ast.FuncType
	name       string
	unit       string
	tags       tagSet
	paramNames []ast.Expr
	hasResults bool
}

func genExportFunc(name, unit, signature string, tags tagSet) *ast.FuncDecl {
	f := exportFunc{
		name: name,
		unit: unit,
		tags: tags,
	}
	f.initFuncType(signature)
	f.paramNames = f.scanParamNames()
	f.hasResults = f.FuncType.Results != nil
	f.rewriteUnits()
	return f.genFuncDecl()
}

func (f *exportFunc) initFuncType(signature string) {
	expr, err := parser.ParseExpr(signature)
	if err != nil {
		throw("parsing %s: %w", signature, err)
	}

	n, ok := expr.(*ast.FuncType)
	if !ok {
		throw("expected *ast.FuncType")
	}
	f.FuncType = n
}

func (f exportFunc) scanParamNames() []ast.Expr {
	var names []ast.Expr
	for _, param := range f.FuncType.Params.List {
		for _, ident := range param.Names {
			names = append(names, ident)
		}
	}
	return names
}

func (f exportFunc) rewriteUnits() {
	astutil.Apply(f.FuncType, nil, func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.SelectorExpr:
			xIdent, ok := n.X.(*ast.Ident)
			if ok && xIdent.Name == "simd" && n.Sel.Name == "Unit" {
				unit := dictUnitsToNative[f.unit]
				ident := &ast.Ident{Name: unit}
				c.Replace(ident)
			}
		}
		return true
	})
}

func (f exportFunc) genFuncDecl() *ast.FuncDecl {
	ident := &ast.Ident{
		Name: f.name + `_` + f.unit,
	}

	return &ast.FuncDecl{
		Name: ident,
		Type: f.FuncType,
		Body: f.genFuncBody(),
	}
}

func (f exportFunc) genFuncBody() *ast.BlockStmt {
	block := new(ast.BlockStmt)
	block.List = []ast.Stmt{f.genSwitch()}
	if f.hasResults {
		block.List = append(block.List, &exprStmtPanicInvalidDispatch)
	}
	return block
}

func (f exportFunc) genSwitch() *ast.SwitchStmt {
	body := new(ast.BlockStmt)
	for isa, vect := range f.tags.matchUnit(f.unit) {
		body.List = append(body.List, f.genClause(isa, vect))
	}

	return &ast.SwitchStmt{
		Tag:  &exprRuntimeDispatch,
		Body: body,
	}
}

func (f exportFunc) genClause(isa string, vect string) *ast.CaseClause {
	return &ast.CaseClause{
		List: f.genClauseList(isa),
		Body: f.genClauseBody(isa, vect),
	}
}

func (f exportFunc) genClauseList(isa string) []ast.Expr {
	return []ast.Expr{&ast.BasicLit{
		Kind:  token.STRING,
		Value: `"` + isa + `"`,
	}}
}

func (f exportFunc) genClauseBody(isa string, vect string) []ast.Stmt {
	name := unexport(f.name)
	name = name + `_` + isa + `_` + vect

	ident := &ast.Ident{
		Name: name,
	}

	call := &ast.CallExpr{
		Fun:  ident,
		Args: f.paramNames,
	}

	body := make([]ast.Stmt, 1)
	if f.hasResults {
		body[0] = &ast.ReturnStmt{
			Results: []ast.Expr{call},
		}
	} else {
		body[0] = &ast.ExprStmt{
			X: call,
		}
	}
	return body
}

// constant nodes
var (
	identSimd    = ast.Ident{Name: "simd"}
	identRuntime = ast.Ident{Name: "RuntimeD"}
	identISA     = ast.Ident{Name: "ISA"}
	identUnit    = ast.Ident{Name: "Unit"}

	exprRuntimeDispatch = ast.CallExpr{Fun: &ast.SelectorExpr{
		X:   &identRuntime,
		Sel: &identISA,
	}}

	exprSimdUnit = ast.SelectorExpr{
		X:   &identSimd,
		Sel: &identUnit,
	}

	litInvalidDispatch = ast.BasicLit{Kind: token.STRING, Value: `"invalid dispatch: "`}

	exprInvalidDispatch = ast.BinaryExpr{
		Op: token.ADD,
		X:  &litInvalidDispatch,
		Y:  &exprRuntimeDispatch,
	}

	exprStmtPanicInvalidDispatch = ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "panic"},
			Args: []ast.Expr{&exprInvalidDispatch},
		},
	}
)

type tagSet map[string]map[string]bool

func newTagSet(targets []simd.Dispatch) tagSet {
	tags := make(tagSet)
	for _, target := range targets {
		isa, vect := target.ISA(), target.VecT()
		if tags[isa] == nil {
			tags[isa] = make(map[string]bool)
		}
		tags[isa][vect] = true
	}

	return tags
}

// returns an iterator over units supported by all isas in the tag set
func (tags tagSet) intersectUnits() iter.Seq[string] {
	units := make(map[string]int)

	// count appearances
	for _, vset := range tags {
		for vect := range vset {
			unit := dictVecTsToUnit[vect]
			units[unit] += 1
		}
	}

	// remove units that aren't seen in every isa
	expect := len(tags)
	for unit, count := range units {
		if count != expect {
			delete(units, unit)
		}
	}

	return maps.Keys(units)
}

func (tags tagSet) matchUnit(unit string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for isa, vset := range tags {
			for vect := range vset {
				if unit != dictVecTsToUnit[vect] {
					continue
				}
				if !yield(isa, vect) {
					return
				}
			}
		}
	}
}
