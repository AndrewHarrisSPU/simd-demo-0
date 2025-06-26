package gensimd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/ast/astutil"
)

func (env *genEnv) reifyVirtualAPI() {
	env.cloneAPI("dispatch")
	env.cloneAPI("types")

	for isa, dst := range env.isa {
		srcPath := filepath.Join(env.apiDir, isa+".go")
		f, err := parser.ParseFile(token.NewFileSet(), srcPath, nil, parser.SkipObjectResolution)
		if err != nil {
			throw("parsing virtual API: %w", err)
		}

		astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.FuncDecl:
				if env.ruleReifyVirtualOp(n) {
					env.publishNode(dst, n)
				}
			}
			return true
		})
	}
}

func (env *genEnv) cloneAPI(name string) {
	srcPath := filepath.Join(env.apiDir, name+".go")
	src, err := os.Open(srcPath)
	if err != nil {
		throw("cloning api: %w", err)
	}
	defer src.Close()

	dst := env.initAPI(name)
	_, err = io.Copy(dst, src)
	if err != nil {
		throw("cloning api: %w", err)
	}
}

func (env *genEnv) ruleReifyVirtualOp(n *ast.FuncDecl) bool {
	name := n.Name.Name
	r, _ := utf8.DecodeRuneInString(name)
	if unicode.IsUpper(r) {
		throw("unexpected exported function name %s", name)
	}

	s, ok := parseSymbol(name)
	if !ok {
		throw("bad symbol %s", name)
	}

	if !env.virtualOps[s.exportOp()] {
		return false
	}

	return true
}
