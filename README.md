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

- `*_simd.go` files in the `simd` directory are slurped up
- code generation (from the `gensimd` directory) generates various concrete impls
- concrete impls are based on `Dispatch` values - mostly substitutions and tagging, respecting overloads
- generated and overloaded impls are written to the `simd` directory
- necessary parts of the `stdlib/simd` directory are also cloned to the `simd` directory
- finally, an `export.go` file is written to the `simd` directory to export routines

Step 2:

Testing reveals:

- a floating point precision issue between 32-bit and 64-bit sigmoid
- a (deliberately) broken overload that yields an incorrect result.

---

### Project Status

It's runnable. Only a very small amount of proposed API and no real SIMD (just emulation) is implemented. However, I think it's looking at organizational concepts in a substantive way. It's running enough to have forced some clarity in the details.

update: Did a lot of code golf, or code jenga in the sense that it's trying to stay ahead of blocks falling over. Proper testing, various edge cases identified. Not a lot of headway on simd problems, but a much better foundation. Would like to implement base-e exponentials the way Highway does, it would be a good test of whether the scheme holds up.

---

### What does this repo demonstrate?

- A two-step workflow for generating simd code - a generative step, then a compile-as-usual step.
- A consistent naming scheme

    Some examples of renaming and tagging:
    1. `func Wip` -> `func Wip_F64` (`Wip` is exported, and rewritten per-`simd.Unit` flavor)
    2. `simd.VecT` -> `F64s` (substitution of a scalable float64 vector type)
    3. `func foo` -> `foo_emuB_F64s` (`foo` is rewritten per simd ISA / vector type combination)
    4. `func foo_emuA_F64x4` overloads `foo` for that simd isa and vector type
    5. `sigmoid_simd.go` -> `sigmoid_emuA_native.go` (rewritten files indicate simd and cpu ISA)

- How some concepts like pure/virtual API surface, templates - employed by Google Highway, in C++ - can be realized via code generation in Go.

---

### Portable? Extensible?

Big picture, I think portable vs. extensible are in some tension.

In a sense this code generation approach suggests two APIs:

- A virutal API in some `stdlib/simd/api.go` that is _universally portable_
- A concrete API generated in the local `simd` package that is _extensible_

The virutal API is defined in a way that enables some useful type checking, but it's virtual. It doesn't get called and doesn't impelement anything.

The concrete API is self-contained, and can be called upon through symbols in `simd/export.go`.

Also in `stdlib/simd/` are per-ISA impelementation files `emuA.go` and `emuB.go`. These files contain the implementations of the virtual API. The idea here is that with low-level intrinsics, files like `avx2.go` could similarly implement each symbol of the virtual, portable API with very simple redirects to the intrinsics. Or, at the cost of maintaining mapping from virtual symbols to low-level intrinsics some other way, some generative ideas are also maybe possible here.

---

### Does this scale?

Not sure. I think this could scale to handle different ISAs. I'm not sure what is possible or convenient or pleasent w/r/t a much larger virtual API surface, including masking, gather/scatter, etc.