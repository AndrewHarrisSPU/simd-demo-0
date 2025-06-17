package gensimd

import (
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gosimd/stdlib/simd"
)

func GenSimdCode() {
	dA, err := simd.NewDispatch("emuA", "float64", 4)
	if err != nil {
		log.Fatal(err)
	}

	dB, err := simd.NewDispatch("emuB", "float64", 4)
	if err != nil {
		log.Fatal(err)
	}

	ds := []simd.Dispatch{dA, dB}

	clearSimd()
	generateFiles(ds)
	cloneImpls()
}

func clearSimd() {
	wd := os.DirFS(filepath.Join(".", "simd"))
	paths, err := fs.Glob(wd, "*.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range paths {
		os.Remove(path)
	}
}

func generateFiles(ds []simd.Dispatch) {
	wd := os.DirFS(".")
	paths, err := fs.Glob(wd, "*_simd.go")
	if err != nil {
		log.Fatal(err)
	}

	ex := newExportSet(ds)

	for _, srcPath := range paths {
		for _, d := range ds {
			fset := token.NewFileSet()
			src, err := parser.ParseFile(fset, srcPath, nil, 0)
			if err != nil {
				log.Fatal(err)
			}

			rewriteFile(d, ex, fset, src)

			dstPath := strings.Replace(srcPath, "simd", d.Tag(), 1)
			dst, err := os.Create(filepath.Join("simd", dstPath))
			if err != nil {
				log.Fatal(err)
			}
			defer dst.Close()

			format.Node(dst, fset, src)
		}
	}

	// exports
	dst, err := os.Create(filepath.Join(".", "simd", "export.go"))
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	dst.WriteString(`package simd` + "\n\n")

	for _, fn := range ex.publish {
		format.Node(dst, token.NewFileSet(), fn)
		dst.WriteString("\n")
	}
}

func cloneImpls() {
	cp := func(dir, srcPath string) {
		src, err := os.Open(filepath.Join(dir, srcPath))
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(".", "simd", srcPath))
		if err != nil {
			log.Fatal(err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			log.Fatal(err)
		}
	}

	// impls from stdlib/simd
	dir := os.DirFS(filepath.Join(".", "stdlib", "simd"))
	paths, err := fs.Glob(dir, "*.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, srcPath := range paths {
		if srcPath == "api.go" {
			continue
		}

		cp(filepath.Join(".", "stdlib", "simd"), srcPath)
	}
}
