package gensimd

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"gosimd/stdlib/simd"

	"golang.org/x/tools/go/ast/astutil"
)

func rewriteFile(d simd.Dispatch, ex *exportSet, fset *token.FileSet, f *ast.File) {
	dispatchTags := makeSymbolMap(d)

	// find FuncDecls, generate replacements, store association
	overloads := make(map[string]string)
	for n := range ast.Preorder(f) {
		switch x := n.(type) {
		case *ast.FuncDecl:
			name := x.Name.Name
			if detectOverload(name) {
				overloads[name] = name
			} else {
				overloads[name] = name + `_` + d.Tag()
			}
		}
	}

	// rule: delete "simd" import (assumed to be stdlib simd package)
	astutil.DeleteImport(fset, f, `gosimd/stdlib/simd`)

	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {

		// rule: set package name to "simd"
		case *ast.File:
			x.Name.Name = "simd"

		// rule: simd.Symbol -> symbol_tag
		case *ast.SelectorExpr:
			if id, ok := x.X.(*ast.Ident); ok {
				if id.Name == "simd" {
					replace := dispatchTags[x.Sel.Name]
					c.Replace(&ast.Ident{Name: replace})
				}
			}

		case *ast.FuncDecl:
			// rule: maybe export Foo
			ex.genExportFunc(x)

			if replace, found := overloads[x.Name.Name]; found {
				overloadTag := parseOverloadTag(replace)

				condA := overloadTag != d.Tag()

				condB := x.Name.Name != replace
				_, condC := overloads[replace]

				if condA && condC {
					c.Delete()
				} else if condB && condC {
					c.Delete()
				}

				// rule: if overloaded, delete func node
				// if x.Name.Name != replace {
				// 	if _, overloaded := overloads[replace]; overloaded {
				// 		c.Delete()
				// 	}
				// }
				// rule: func symbol -> func symbol_tag
				x.Name.Name = unexport(replace)
			}

		// rule: foo() called, and foo() defined in this file -> foo_tag()
		case *ast.CallExpr:
			if id, ok := x.Fun.(*ast.Ident); ok {
				if replace, found := overloads[id.Name]; found {
					id.Name = replace
				}
			}
		}

		return true
	})
}

// detect if a function looks like a deliberate overload
func detectOverload(name string) bool {
	ss := strings.Split(name, "_")
	if len(ss) < 3 {
		return false
	}
	tail := len(ss) - 1

	if !slices.Contains(simdISAs, ss[tail-1]) {
		return false
	}
	if !slices.Contains(simdVecTypes, ss[tail]) {
		return false
	}
	return true
}

func parseOverloadTag(name string) string {
	ss := strings.Split(name, "_")
	tail := len(ss) - 1
	return ss[tail-1] + "_" + ss[tail]
}

// generation of runtime dispatch of exported SIMD funcs
// rule: Foo([]simd.T) -> Foo_T([]T)
type exportSet struct {
	support map[string][]string
	seen    map[string]bool
	publish []*ast.FuncDecl
}

func newExportSet(ds []simd.Dispatch) *exportSet {
	ex := new(exportSet)
	ex.seen = make(map[string]bool)
	ex.support = make(map[string][]string)

	for _, d := range ds {
		t, tag := d.T(), d.Tag()
		ex.support[t] = append(ex.support[t], tag)
	}

	return ex
}

func (ex *exportSet) genExportFunc(fn *ast.FuncDecl) {
	// check that fn needs to export
	// TODO: think about different way to do this, without ex.seen maybe?
	fnName := fn.Name.Name
	if !ast.IsExported(fnName) {
		return
	}
	if ex.seen[fnName] {
		return
	}
	ex.seen[fnName] = true

	paramNames := genExportParams(fn.Type.Params.List)
	hasRet := fn.Type.Results != nil

	// per unit type type
	for t, tags := range ex.support {
		var cs []*ast.CaseClause
		// per isa, generate a switch clause
		for _, tag := range tags {
			isa, vect := simd.ParseTag(tag)
			c := genExportClause(isa, vect, fnName, paramNames, hasRet)
			cs = append(cs, c)
		}

		// generate a switch statement
		sw := genExportSwitch(cs)

		// generate a new function
		gfn := new(ast.FuncDecl)
		gfn.Name = genExportName(fnName, simd.NativeToSimdT[t])

		// function params / results shouldn't change ...
		gfn.Type = fn.Type

		// the function body is just the switch statement
		gfn.Body = genExportBody(sw, hasRet)

		// hold onto the generated function
		ex.publish = append(ex.publish, gfn)
	}
}

func genExportName(fnName, vect string) *ast.Ident {
	return &ast.Ident{
		Name: fnName + `_` + vect,
	}
}

func genExportBody(sw *ast.SwitchStmt, hasRet bool) *ast.BlockStmt {
	block := new(ast.BlockStmt)
	block.List = []ast.Stmt{sw}
	if hasRet {
		oops := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{
					Name: "panic",
				},
				Args: []ast.Expr{
					&ast.BinaryExpr{
						Op: token.ADD,
						X: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"invalid dispatch"`,
						},
						Y: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "RuntimeD",
								},
								Sel: &ast.Ident{
									Name: "ISA",
								},
							},
						},
					},
				},
			},
		}
		block.List = append(block.List, oops)
	}

	return block
}

func genExportParams(list []*ast.Field) []string {
	var params []string
	for _, field := range list {
		for _, name := range field.Names {
			params = append(params, name.Name)
		}
	}
	return params
}

func genExportSwitch(clauses []*ast.CaseClause) *ast.SwitchStmt {
	sw := new(ast.SwitchStmt)
	sw.Tag = &ast.Ident{
		Name: "RuntimeD.ISA()",
	}
	sw.Body = new(ast.BlockStmt)
	for _, clause := range clauses {
		sw.Body.List = append(sw.Body.List, clause)
	}

	return sw
}

func genExportClause(isa, vect, fnName string, params []string, hasRet bool) *ast.CaseClause {
	clause := new(ast.CaseClause)

	// predicate
	lit := &ast.BasicLit{
		Kind:  token.STRING,
		Value: `"` + isa + `"`,
	}
	clause.List = append(clause.List, lit)

	// body
	call := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: unexport(fnName) + `_` + isa + `_` + vect,
		},
	}

	for _, param := range params {
		id := &ast.Ident{
			Name: param,
		}
		call.Args = append(call.Args, id)
	}

	if hasRet {
		clause.Body = append(clause.Body, &ast.ReturnStmt{
			Results: []ast.Expr{call},
		})
	} else {
		clause.Body = append(clause.Body, &ast.ExprStmt{
			X: call,
		})
	}

	return clause
}
