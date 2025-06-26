package gensimd

import (
	"fmt"
	"log"
	"path/filepath"

	"gosimd/stdlib/simd"
)

type genError struct {
	error
}

func throw(f string, args ...any) {
	err := fmt.Errorf(f, args...)
	panic(genError{err})
}

func catch(handle func(error)) {
	if e := recover(); e != nil {
		if err, ok := e.(genError); ok {
			handle(err)
		} else {
			panic(e)
		}
	}
}

func GenSimdCode() {
	defer catch(func(err error) {
		log.Fatal(err)
	})

	// TODO: in a not-pretending reality parts of this would be less hard-coded
	dA32, err := simd.NewDispatch("emuA", "float32", 8)
	if err != nil {
		log.Fatal(err)
	}

	dA64, err := simd.NewDispatch("emuA", "float64", 4)
	if err != nil {
		log.Fatal(err)
	}

	dB32, err := simd.NewDispatch("emuB", "float32", 64)
	if err != nil {
		log.Fatal(err)
	}

	dB64, err := simd.NewDispatch("emuB", "float64", 4)
	if err != nil {
		log.Fatal(err)
	}

	apiDir := filepath.Join("stdlib", "simd")
	srcDir := filepath.Join("simd")
	dstDir := filepath.Join("simd")

	ds := []simd.Dispatch{dA32, dA64, dB32, dB64}

	env := newGenEnv(ds, apiDir, srcDir, dstDir)
	env.scanSources()
	env.rewriteSources()
	env.reifyVirtualAPI()
	env.generateExportFile()
	env.commitFiles()
}
