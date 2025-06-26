package gensimd

import (
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func (env *genEnv) scanAST(f *ast.File, name string) {
	for _, decl := range f.Decls {
		env.ruleScanFileDecl(decl, name)
	}

	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.SelectorExpr:
			env.ruleScanVirtualOp(n)
		case *ast.FuncDecl:
			env.ruleScanFuncDecl(n)
		}
		return true
	})
}

func (env *genEnv) ruleScanFileDecl(decl ast.Decl, name string) {
	n, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}
	if n.Tok != token.IMPORT {
		throw("%s: file-level declaration: want import, got %s)", name, n.Tok)
	}

	src := env.src[name]
	src.importDecls = append(src.importDecls, n)
}

func (env *genEnv) ruleScanVirtualOp(parent *ast.SelectorExpr) {
	n, ok := parent.X.(*ast.Ident)
	if !ok || n.Name != "simd" {
		return
	}

	name := parent.Sel.Name
	if setOps[name] {
		env.virtualOps[name] = true
	}
}

func (env *genEnv) ruleScanFuncDecl(n *ast.FuncDecl) {
	name := n.Name.Name

	// export
	if ast.IsExported(name) {
		if setOps[name] {
			throw("exported func %s collides with virtual API operation", name)
		}
		if env.local[unexport(name)] {
			throw("exported func %s collides with local func %s", name, unexport(name))
		}
		env.ruleScanFuncTypeText(n.Type, name)
		return
	}

	// overload
	if s, ok := parseSymbol(name); ok {
		if setOps[s.exportOp()] {
			throw("overload %s collides with virtual API operation", s)
		}
		if env.isa[s.isa()] == nil {
			return
		}
		env.overload[s] = true
		return
	}

	// local
	if setOps[export(name)] {
		throw("local func %s collides with virtual API operation", name)
	}
	if env.export[export(name)] != "" {
		throw("local func %s collides with exported func %s", name, export(name))
	}
	env.local[name] = true
}

func (env *genEnv) ruleScanFuncTypeText(n *ast.FuncType, name string) {
	env.ruleScanFuncTypeParams(n.Params)

	var clone strings.Builder
	format.Node(&clone, token.NewFileSet(), n)
	env.export[name] = clone.String()
}

func (env *genEnv) ruleScanFuncTypeParams(params *ast.FieldList) {
	for _, param := range params.List {
		n, ok := param.Type.(*ast.SelectorExpr)
		if !ok {
			return
		}
		if ident, ok := n.X.(*ast.Ident); !ok || ident.Name != "simd" {
			return
		}
		if n.Sel.Name == "VecT" {
			throw("invalid FuncType: includes simd.VecT-typed paramater")
		}
	}
}
