package gensimd

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"
)

// expectThrow catches expected panics with appropriate test logic
// `caught` indication is useful in table-driven loops (`continue` if caught)
func expectThrow(t *testing.T, cond bool, do func()) (caught bool) {
	t.Helper()

	// if cond isn't set, don't catch, and don't indicate a panic was caught
	if !cond {
		do()
		return false
	}

	// anticipating a throw ...
	defer catch(func(err error) {
		if _, ok := err.(genError); !ok {
			// an unexpected panic causes test failure
			t.Fatal(err)
		}
		// an expected panic sets `caught` to true
		caught = true
	})

	do()
	// if a panic was expected but not encountered, test fails
	if !caught {
		t.Fatal("expected throw")
	}
	return
}

func TestRuleScanVirtualOp(t *testing.T) {
	var env genEnv
	env.virtualOps = make(map[string]bool)

	exprs := []struct {
		text   string
		lookup string
		expect bool
	}{
		{`simd.Add`, "Add", true},
		{`simd.Nop`, "Nop", false},
		{`dmis.Add`, "Add", false},
	}

	for _, expr := range exprs {
		nn, err := parser.ParseExpr(expr.text)
		if err != nil {
			t.Fatal(err)
		}
		n, ok := nn.(*ast.SelectorExpr)
		if !ok {
			t.Fatal("expected *ast.SelectorExpr")
		}
		env.ruleScanVirtualOp(n)

		found := env.virtualOps[expr.lookup]
		if found != expr.expect {
			t.Fatalf("looking up %s: %v, expected %v", expr.lookup, found, expr.expect)
		}
		delete(env.virtualOps, expr.lookup)
	}
}

func TestRuleScanFileDecl(t *testing.T) {
	var env genEnv
	env.src = make(map[string]*source)
	env.src["foo"] = new(source)

	decls := []struct {
		decl   ast.Decl
		throws bool
	}{
		{&ast.GenDecl{Tok: token.IMPORT}, false},
		{&ast.GenDecl{Tok: token.CONST}, true},
		{new(ast.FuncDecl), false},
	}

	for _, item := range decls {
		expectThrow(t, item.throws, func() {
			env.ruleScanFileDecl(item.decl, "foo")
		})
	}

	if len(env.src["foo"].importDecls) != 1 {
		t.Fatal("bad decl count")
	}
	if _, ok := env.src["foo"].importDecls[0].(*ast.GenDecl); !ok {
		t.Fatal("expected *ast.GenDecl")
	}
}

func TestRuleScanFuncDecl(t *testing.T) {
	var env genEnv
	env.isa = make(map[string]*strings.Builder)
	env.initISA("emuA")
	env.export = make(map[string]string)
	env.local = make(map[string]bool)
	env.overload = make(map[symbol]bool)

	fns := []struct {
		name string
		kind string
	}{
		{"Foo", "export"},
		{"Baz", "export"},
		{"Add", ""},
		{"add", ""},
		{"add_emuA_F64x4", ""},
		{"bar", "local"},
		{"baz_emuA_F64x4", "overload"},
		{"foo", ""},
		{"Bar", ""},
		{"foo_bad_F64x4", ""},
		{"foo_emuA_F64s", ""},
	}

	n, err := parser.ParseExpr("func(){}")
	if err != nil {
		t.Fatal(err)
	}
	lit, ok := n.(*ast.FuncLit)
	if !ok {
		t.Fatal("expected *ast.FuncLit")
	}

	for _, fn := range fns {
		n := &ast.FuncDecl{
			Name: &ast.Ident{Name: fn.name},
			Type: lit.Type,
			Body: lit.Body,
		}

		if caught := expectThrow(t, fn.kind == "", func() {
			env.ruleScanFuncDecl(n)
		}); caught {
			continue
		}

		switch fn.kind {
		case "export":
			if !(n.Name.Name == "Foo" || n.Name.Name == "Baz") {
				t.Fatalf("expected export: %s", fn.name)
			}
		case "overload":
			if n.Name.Name != "baz_emuA_F64x4" {
				t.Fatalf("expected overload: %s", fn.name)
			}
		case "local":
			if n.Name.Name != "bar" {
				t.Fatalf("expected local: %s", fn.name)
			}
		}
	}
}

func TestRuleScanFuncTypeText(t *testing.T) {
	var env genEnv
	env.export = make(map[string]string)

	parse := func(s string) *ast.FuncType {
		t.Helper()
		expr, err := parser.ParseExpr(s)
		if err != nil {
			t.Fatal(s, err)
		}
		ft, ok := expr.(*ast.FuncType)
		if !ok {
			t.Fatal("expected *ast.FuncType")
		}
		return ft
	}

	pre := parse("func(a, b any) (c, d any)")
	env.ruleScanFuncTypeText(pre, "Foo")
	post := parse(env.export["Foo"])

	if !reflect.DeepEqual(pre, post) {
		t.Fatal("bad clone")
	}

	bad := parse("func(vect simd.VecT)")
	expectThrow(t, true, func() {
		env.ruleScanFuncTypeText(bad, "bad")
	})
}

func TestRuleRewriteFileDecl(t *testing.T) {
	var env genEnv

	f := new(ast.File)
	f.Decls = []ast.Decl{
		new(ast.GenDecl),
		new(ast.FuncDecl),
	}

	env.ruleRewriteFileDecls(f)

	if len(f.Decls) != 1 {
		t.Fatal("bad decl count")
	}

	n := f.Decls[0]
	if _, ok := n.(*ast.FuncDecl); !ok {
		t.Fatal("expected *ast.GenDecl")
	}
}

func TestRuleRewriteSimdSelector(t *testing.T) {
	var env genEnv
	env.simd = make(map[string]string)
	env.simd["ISA"] = "emuA"
	env.simd["Unit"] = "F64x4"

	exprs := []struct {
		text string
		want string
	}{
		{"simd.ISA", "emuA"},
		{"simd.Unit", "F64x4"},
		{"simd.Nop", ""},
		{"dmis.ISA", "dmis.ISA"},
	}

	for _, expr := range exprs {
		n, err := parser.ParseExpr(expr.text)
		if err != nil {
			t.Fatal(err)
		}
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			t.Fatal("expected *ast.SelectorExpr")
		}

		got := new(strings.Builder)
		if caught := expectThrow(t, expr.want == "", func() {
			format.Node(got, token.NewFileSet(), env.ruleRewriteSimdSelector(sel))
		}); caught {
			continue
		}

		if expr.want != got.String() {
			t.Fatalf("%s: want %s, got %s", expr.text, expr.want, got.String())
		}
	}
}

func TestRuleRewriteCallExpr(t *testing.T) {
	var env genEnv
	env.local = make(map[string]bool)
	env.export = make(map[string]string)
	env.simd = make(map[string]string)
	env.local["foo"] = true
	env.export["Baz"] = "..."
	env.simd["Tag"] = "emuA_F64x4"

	calls := []struct {
		name   string
		expect string
	}{
		{"foo", "foo_emuA_F64x4"},
		{"bar_emuA_F64x4", "bar_emuA_F64x4"},
		{"Baz", "baz_emuA_F64x4"},
		{"panic", "panic"},
	}

	for _, call := range calls {
		expr := &ast.CallExpr{
			Fun: &ast.Ident{Name: call.name},
		}

		env.ruleRewriteCallExpr(expr)
		ident, ok := expr.Fun.(*ast.Ident)
		if !ok {
			t.Fatal("expected *ast.Ident")
		}

		if call.expect != ident.Name {
			t.Fatalf("want %s, got %s", call.expect, ident.Name)
		}
	}
}

func TestRuleDeleteOverload(t *testing.T) {
	var env genEnv
	env.dst = make(map[string]*strings.Builder)
	env.isa = make(map[string]*strings.Builder)
	env.overload = make(map[symbol]bool)
	env.simd = make(map[string]string)
	env.initISA("emuA")
	env.initDst("filename", "emuA")
	env.simd["Tag"] = "emuA_F64x4"
	env.overload["foo_emuA_F64x4"] = true

	names := []struct {
		text string
		want bool
	}{
		{"foo_emuA_F64x4", true},
		{"foo", true},
		{"bar_emuA_F64x4", true},
		{"baz", false},
	}

	fn := &ast.FuncDecl{
		Name: new(ast.Ident),
		Type: new(ast.FuncType),
		Body: new(ast.BlockStmt),
	}

	for _, name := range names {
		fn.Name.Name = name.text

		remove := env.ruleDeleteOverload(fn, "filename")
		if name.want != remove {
			t.Fatalf("%s: want %v, got %v", name.text, name.want, remove)
		}
	}
}

func TestRuleRewriteFuncDecl(t *testing.T) {
	var env genEnv
	env.isa = make(map[string]*strings.Builder)
	env.overload = make(map[symbol]bool)
	env.simd = make(map[string]string)
	env.simd["Tag"] = "emuA_F64x4"
	env.initISA("emuA")

	names := []struct {
		text string
		want string
	}{
		{"Foo", "foo_emuA_F64x4"},
		{"foo", "foo_emuA_F64x4"},
	}

	fn := &ast.FuncDecl{
		Name: new(ast.Ident),
		Type: new(ast.FuncType),
		Body: new(ast.BlockStmt),
	}

	for _, name := range names {
		fn.Name.Name = name.text

		env.ruleRewriteFuncDecl(fn)
		if name.want != fn.Name.Name {
			t.Fatalf("%s: want %s, got %s", name.text, name.want, fn.Name.Name)
		}
	}
}

func TestRuleReifyVirtualOp(t *testing.T) {
	var env genEnv
	env.virtualOps = make(map[string]bool)
	env.virtualOps["Add"] = true

	ops := []struct {
		name    string
		reifies bool
		throws  bool
	}{
		{"add_emuA_F64x4", true, false},
		{"mul_emuA_F64x4", false, false},
		{"Add_emuA_F64x4", false, true},
		{"add", false, true},
	}

	fn := &ast.FuncDecl{
		Name: new(ast.Ident),
		Type: new(ast.FuncType),
		Body: new(ast.BlockStmt),
	}

	for _, op := range ops {
		fn.Name.Name = op.name

		var reified bool
		if caught := expectThrow(t, op.throws, func() {
			reified = env.ruleReifyVirtualOp(fn)
		}); caught {
			continue
		}

		if op.reifies != reified {
			t.Fatalf("%s: want %v, got %v", op.name, op.reifies, reified)
		}
	}
}
