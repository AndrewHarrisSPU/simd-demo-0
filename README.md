This repo explores some concepts for organization, generation, and dispatch of SIMD code in Go.

Nomenclature and ideas borrowed heavily from:
- [Go Simd Proposal: Issue #73787](https://github.com/golang/go/issues/73787)
- [Google Highway](https://github.com/google/highway)
- [Another sketch of generated Go SIMD code](https://docs.google.com/document/d/1Wfidbmy4KYiNYfrFJfqrrFYwuT2ITP5vceku8_LFsXU/)

Run it:

1. `go run .` to generate
2. `go test` to test

---

### What happened?

Step 1:
- `*_simd.go` files are slurped up
- code gen (from the `gensimd` directory) generates concrete impls, based on `Dispatch` data, respecting overloads.
- generated and overloaded impls are written to the `simd` directory
- parts of the `stdlib/simd` directory are also written to the `simd` directory
- finally, an `export.go` file is written to the `simd` directory to export routines

Step 2:

Testing reveals a (deliberately) broken overload that yields an incorrect result.

To fix, look in `sigmoid_simd.go` - a comment shows what's wrong.

---

### Project Status:

It's runnable. Only a very small amount of proposed API and no real SIMD (just emulation) is implemented. However, I think it's looking at organizational concepts in a substantive way. It's running enough to have forced some clarity in the details.

---

### What does this repo demonstrate?

- A two-step workflow for generating simd code - a generative step, then a compile-as-usual step.
- A consistent naming scheme. Generated functions and files are tagged with ISA and numeric types. Overloads result from following the scheme.

    Some examples of renaming:
    1. `func Wip` -> `func Wip_F64`
    2. `simd.VecT` -> `F64s`
    3. `sigmoid.go` -> `sigmoid_emuA_F64x4.go`

- How some concepts like pure/virtual API surface, templates - employed by Google Highway, in C++ - can be realized via code generation in Go.

---

### Portable? Extensible?

Big picture, I think portable vs. extensible are in some tension.

In a sense this code generation approach suggests two APIs:
- A virutal API in `stdlib/simd/api.go` that is _universally portable_
- A concrete API generated in the local `simd` package that is _extensible_

The virutal API is defined in a way that enables some useful type checking, but it's virtual. It doesn't get called and doesn't impelement anything.

The concrete API is self-contained, and can be called upon through symbols in `simd/export.go`.

Also in `stdlib/simd/` are per-ISA impelementation files `emuA.go` and `emuB.go`. These files contain the implementations of the virtual API. The idea here is that with low-level intrinsics, files like `avx2.go` could similarly implement each symbol of the virtual, portable API with very simple redirects to the intrinsics. Or, at the cost of maintaining mapping from virtual symbols to low-level intrinsics some other way, some generative ideas are also maybe possible here.

--- 

### Does this scale?

Not sure. I think this could scale to handle different ISAs. I'm not sure what is possible or convenient or pleasent w/r/t a much larger virtual API surface, including masking, gather/scatter, etc.