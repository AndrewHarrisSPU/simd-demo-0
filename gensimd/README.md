### gensimd rules

#### scanning phase

1. virtual operations are tracked (e.g. `simd.Add`)
2. only imports and functions are read from `_simd.go` files
    - other file-level declarations throw
3. Scanned functions are categorized like:
    - exported (e.g. `Foo`) (exports may not use simd.Vect as parameters or results)
    - local (e.g. `foo`)
    - overloads (e.g. `foo_emuA_F64x4`)
    - namespace collisions (e.g. resulting from delcaration of both `Foo` and `foo`) may throw
4. function types of exports are retained for the export generation phase

#### source rewriting phase

5. per-source-file, per-isa, imports are written
6. `simd` selector expressions are rewritten (e.g. `simd.Unit` -> `F64`)
7. per-dispatch, function calls are tagged (e.g. `foo(vs)` -> `foo_emuA_F64x4(vs)`)
8. per-dispatch, exported and local function declarations are tagged (e.g. `func Foo(){}` -> `func foo_emuA_F64x4(){}`)

#### reification phase

9. implementations of the virtual `simd` dispatch and types are cloned
10. per-isa, virtual operations (see rule 1) are cloned

#### export generation phase

11. per-exported function, per-unit-type, export functions are derived that perform runtime dispatch based on CPU features. The functions just switch over a `dispatch.ISA()` value, into matching tagged implementations.

#### publishing phase

12. existing files matching any known `_isa_arch` tag, and `export.go`, `dispatch.go`, `types.go` are deleted
13. the following files are regenerated:
    - per-isa rewrites of `_simd.go` files (e.g. `foo_simd.go` -> `foo_emuA_native.go`)
    - per-isa clones of virtual API operations (e.g. `emuA_native.go`)
    - `export.go`, synthesized exports
    - `dispatch.go`, `types.go` cloned from standard library

---

todo: some filenames might be worth discouraging, like `dispatch_simd.go`