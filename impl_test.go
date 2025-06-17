package main

import (
	"gosimd/simd"
	"slices"
	"testing"
)

func TestSigmoid(t *testing.T) {
	dA, err := simd.NewDispatch("emuA", "float64", 4)
	if err != nil {
		t.Fatal(err)
	}

	dB, err := simd.NewDispatch("emuB", "float64", 4)
	if err != nil {
		t.Fatal(err)
	}

	xs := []float64{1.0, 1.0, 1.0, 1.0}
	ys := []float64{0.0, -1.0, 1.0, 100.0}
	a, b := slices.Clone(xs), slices.Clone(xs)

	simd.RuntimeD = dA
	simd.Sigmoid_F64(a, ys)
	simd.RuntimeD = dB
	simd.Sigmoid_F64(b, ys)

	if !slices.Equal(a, b) {
		t.Errorf("Sigmoid: a != b: %+v, %+v", a, b)
	}
}

func TestWip(t *testing.T) {
	dA, err := simd.NewDispatch("emuA", "float64", 4)
	if err != nil {
		t.Fatal(err)
	}

	dB, err := simd.NewDispatch("emuB", "float64", 4)
	if err != nil {
		t.Fatal(err)
	}

	vs := [][]float64{
		{1.0, 2.0, 1.0, 1.0, 1.0},
		{1.0, 2.0, 1.0, 1.0, 1.0},
		{1.0, 2.0, 1.0, 1.0, 1.0},
	}

	simd.RuntimeD = dA
	a := simd.Wip_F64(vs)
	simd.RuntimeD = dB
	b := simd.Wip_F64(vs)

	if a != b {
		t.Errorf("Wip: a != b: %f, %v", a, b)
	}
}
