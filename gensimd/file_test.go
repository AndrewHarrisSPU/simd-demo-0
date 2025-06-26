package gensimd

import (
	"go/ast"
	"go/token"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"gosimd/stdlib/simd"
)

func TestFileFoo(t *testing.T) {
	defer catch(func(err error) {
		t.Fatal(err)
	})

	d, err := simd.NewDispatch("emuA", "float64", 4)
	if err != nil {
		t.Fatal(err)
	}

	targets := []simd.Dispatch{d}
	apiDir := filepath.Join("testdata", "stdlib", "simd")
	srcDir := filepath.Join("testdata", "simd")
	dstDir := filepath.Join("testdata", "simd")
	env := newGenEnv(targets, apiDir, srcDir, dstDir)

	env.scanSources()
	testFileFooKeys(t, env)
	testFileFooScanAST(t, env)
	testFileFooExportText(t, env)

	env.rewriteSources()
	testFileFooPublishFile(t, env)

	env.generateExportFile()
	testFileFooExport(t, env)
	testFileScanList(t, env)
	testFileDeleteList(t, env)
	testFileWriteList(t, env)
}

func testFileFooKeys(t *testing.T, env *genEnv) {
	srcKeys := slices.Collect(maps.Keys(env.src))
	wantSrcKeys := []string{"foo"}
	slices.Sort(srcKeys)
	if !slices.Equal(wantSrcKeys, srcKeys) {
		t.Fatalf("want %+v, got %+v", wantSrcKeys, srcKeys)
	}

	dstKeys := slices.Collect(maps.Keys(env.dst))
	wantDstKeys := []string{"foo_emuA"}
	slices.Sort(dstKeys)
	if !slices.Equal(wantDstKeys, dstKeys) {
		t.Fatalf("want %+v, got %+v", wantDstKeys, dstKeys)
	}
}

func testFileFooScanAST(t *testing.T, env *genEnv) {
	if len(env.src["foo"].importDecls) < 1 {
		t.Fatal()
	}

	decl := env.src["foo"].importDecls[0]
	gdecl, ok := decl.(*ast.GenDecl)
	if !ok {
		t.Fatal("expected *ast.GenDecl")
	}
	if gdecl.Tok != token.IMPORT {
		t.Fatal("expeted import declaration")
	}
	if len(gdecl.Specs) != 1 {
		t.Fatal("expected 1 spec")
	}
	n, ok := gdecl.Specs[0].(*ast.ImportSpec)
	if !ok {
		t.Fatal("expected *ast.ImportSpec")
	}
	if n.Path.Value != `"math"` {
		t.Fatalf(`want "math", got %s`, n.Path.Value)
	}
}

func testFileFooExportText(t *testing.T, env *genEnv) {
	want := `func(in simd.Unit) (out simd.Unit)`
	got := env.export["Foo"]
	if want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func testFileFooPublishFile(t *testing.T, env *genEnv) {
	want, err := os.ReadFile(filepath.Join(env.dstDir, "foo_emuA_native.go"))
	if err != nil {
		t.Fatal(err)
	}

	got := env.dst["foo_emuA"].String()
	if string(want) != got {
		t.Fatalf("\nwant:\n%s\n\ngot:\n%s\n", want, got)
	}
}

func testFileFooExport(t *testing.T, env *genEnv) {
	want, err := os.ReadFile(filepath.Join(env.dstDir, "export.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := env.api["export"].String()

	if string(want) != got {
		t.Fatalf("\nwant:\n%s\n\ngot:\n%s\n", want, got)
	}
}

func testFileScanList(t *testing.T, env *genEnv) {
	want := []string{"foo_simd.go"}
	got := slices.Collect(env.listFilesToScan())
	slices.Sort(got)

	if !slices.Equal(want, got) {
		t.Fatalf("want %+v, got %+v", want, got)
	}
}

func testFileDeleteList(t *testing.T, env *genEnv) {
	want := []string{"dispatch.go", "export.go", "foo_emuA_native.go", "types.go"}
	got := slices.Collect(env.listFilesToDelete())
	slices.Sort(got)

	if !slices.Equal(want, got) {
		t.Fatalf("want %+v, got %+v", want, got)
	}
}

func testFileWriteList(t *testing.T, env *genEnv) {
	want := []string{"dispatch.go", "export.go", "emuA_native.go", "foo_emuA_native.go", "types.go"}

	for got := range env.listFilesToWrite() {
		if !slices.Contains(want, got) {
			t.Fatalf("%s not found in %+v", got, want)
		}
	}
}
