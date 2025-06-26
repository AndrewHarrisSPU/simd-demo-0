package gensimd

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

func (env *genEnv) rewriteSources() {
	for name, src := range env.src {
		// per-source, per-isa, clone imports
		for isa := range env.isa {
			for _, decl := range src.importDecls {
				name_isa := name + `_` + isa
				env.publishNode(env.dst[name_isa], decl)
			}
		}

		// per-source, per-target, rewrite source
		for _, target := range env.targets {
			env.setTarget(target)
			env.rewriteSourceToTarget(name)
		}
	}
}

func (env *genEnv) rewriteSourceToTarget(name string) {
	text := env.src[name].text
	f, err := parser.ParseFile(token.NewFileSet(), "", text, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	env.ruleRewriteFileDecls(f)

	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.SelectorExpr:
			c.Replace(env.ruleRewriteSimdSelector(n))
		case *ast.CallExpr:
			env.ruleRewriteCallExpr(n)
		case *ast.FuncDecl:
			if env.ruleDeleteOverload(n, name) {
				c.Delete()
			} else {
				env.ruleRewriteFuncDecl(n)
			}
		}
		return true
	})

	isa := env.simd["ISA"]
	dst := env.dst[name+`_`+isa]
	env.publishFile(dst, f)
}

func (env *genEnv) ruleRewriteFileDecls(f *ast.File) {
	var retain []ast.Decl
	for _, n := range f.Decls {
		if _, ok := n.(*ast.GenDecl); ok {
			continue
		}
		retain = append(retain, n)
	}

	f.Decls = retain
}

func (env *genEnv) ruleRewriteSimdSelector(parent *ast.SelectorExpr) ast.Node {
	ident, ok := parent.X.(*ast.Ident)
	if !ok || ident.Name != "simd" {
		return parent
	}

	name, ok := env.simd[parent.Sel.Name]
	if !ok {
		throw("env.simd lookup: simd.%s not found", name)
	}

	return &ast.Ident{Name: name}
}

func (env *genEnv) ruleRewriteCallExpr(parent *ast.CallExpr) {
	ident, ok := parent.Fun.(*ast.Ident)
	if !ok {
		return
	}

	s := ident.Name
	if len(env.export[s]) != 0 {
		ident.Name = string(env.tag(s))
	}
	if env.local[s] {
		ident.Name = string(env.tag(s))
	}
}

func (env *genEnv) ruleDeleteOverload(n *ast.FuncDecl, srcName string) (remove bool) {
	fnName := n.Name.Name
	// case 1: remove when n's untagged name parses as a tagged symbol
	if s, ok := parseSymbol(fnName); ok {
		remove = true
		if s.tag() == env.simd["Tag"] {
			srcName_isa := srcName + `_` + s.isa()
			// if n's name matches current target, publish before removal
			env.publishNode(env.dst[srcName_isa], n)
		}
		return
	}

	// case 2: remove when n's tagged name matches an existing manual overload
	s := env.tag(fnName)
	remove = env.overload[s]
	return
}

func (env *genEnv) ruleRewriteFuncDecl(n *ast.FuncDecl) {
	s := env.tag(n.Name.Name)
	n.Name.Name = string(s)
}
