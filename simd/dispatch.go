package simd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var RuntimeD Dispatch

var nativeToSimdUnit = map[string]string{
	"float32": "F32",
	"float64": "F64",
}

type Dispatch struct {
	isa  string
	unit string
	n    int
	vect string
	tag  string
}

func NewDispatch(isa, unit string, lanes int) (Dispatch, error) {
	var d Dispatch
	var err error

	d.isa, err = checkISA(isa)
	if err != nil {
		return Dispatch{}, fmt.Errorf("NewDispatch: %w", err)
	}

	d.unit, err = checkUnit(unit)
	if err != nil {
		return Dispatch{}, fmt.Errorf("NewDispatch: %w", err)
	}

	d.n, err = checkLanes(lanes)
	if err != nil {
		return Dispatch{}, fmt.Errorf("NewDispatch: %w", err)
	}

	d.initVecT()
	d.initTag()

	return d, nil
}

func checkISA(isa string) (string, error) {
	switch isa {
	case `emuA`, `emuB`:
		return isa, nil
	}
	return "", errors.New("unknown isa: " + isa)
}

func checkUnit(t string) (string, error) {
	switch t {
	case `float64`, `float32`:
		return t, nil
	}
	return "", errors.New("unknown type: " + t)
}

func checkLanes(lanes int) (int, error) {
	switch lanes {
	case 1, 2, 4, 8, 16, 32, 64:
		return lanes, nil
	}
	return 0, errors.New("invalid lanes: " + strconv.Itoa(lanes))
}

func (d *Dispatch) initVecT() {
	d.vect = nativeToSimdUnit[d.unit]

	if d.Scalable() {
		d.vect += `s`
	} else {
		d.vect += `x` + d.Lanes()
	}
}

func (d *Dispatch) initTag() {
	d.tag = d.isa + `_` + d.vect
}

func (d *Dispatch) Tag() string {
	return d.tag
}

func ParseTag(tag string) (isa, vect string) {
	isa, vect, found := strings.Cut(tag, `_`)
	if !found {
		panic("invalid tag: " + tag)
	}
	return
}

func (d Dispatch) ISA() string {
	return d.isa
}

func (d Dispatch) Unit() string {
	return d.unit
}

func (d Dispatch) Lanes() string {
	return strconv.Itoa(d.n)
}

func (d Dispatch) VecT() string {
	return d.vect
}

func (d Dispatch) Scalable() bool {
	switch d.isa {
	case `emuA`:
		return false
	case `emuB`:
		return true
	}
	panic("unknown isa")
}

func (d Dispatch) Arch() string {
	switch d.isa {
	//	case `avx2`, `avx512`:
	//		return `amd64`
	case `emuA`, `emuB`:
		return `emuArch`
	}
	panic("unknown isa")
}
