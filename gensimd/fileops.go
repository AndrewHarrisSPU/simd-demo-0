package gensimd

import (
	"errors"
	"io"
	"io/fs"
	"iter"
	"os"
	"strings"
)

func (env *genEnv) listFilesToScan() iter.Seq[string] {
	wd, err := os.OpenRoot(env.srcDir)
	if err != nil {
		throw("scanning %s: %w", env.srcDir, err)
	}
	defer wd.Close()

	names, err := fs.Glob(wd.FS(), "*_simd.go")
	if err != nil {
		throw("scanning %s: %w", env.srcDir, err)
	}

	return func(yield func(string) bool) {
		for _, name := range names {
			if !yield(name) {
				return
			}
		}
	}
}

func (env *genEnv) listFilesToWrite() iter.Seq2[string, *strings.Builder] {
	return func(yield func(string, *strings.Builder) bool) {
		// rewrites
		for name, data := range env.dst {
			i := strings.LastIndex(name, "_")
			isa := name[i+1:]
			name_isa_arch := name + `_` + dictISAtoArch[isa]

			if !yield(name_isa_arch+".go", data) {
				return
			}
		}

		// isa
		for name, data := range env.isa {
			isa_arch := name + `_` + dictISAtoArch[name]

			if !yield(isa_arch+".go", data) {
				return
			}
		}

		// api
		for name, data := range env.api {
			if !yield(name+".go", data) {
				return
			}
		}
	}
}

func (env *genEnv) listFilesToDelete() iter.Seq[string] {
	wd, err := os.OpenRoot(env.dstDir)
	if err != nil {
		throw("listing %s: %w", env.dstDir, err)
	}

	given := []string{"dispatch.go", "types.go", "export.go"}

	return func(yield func(string) bool) {
		defer wd.Close()

		for _, name := range given {
			_, err := wd.Stat(name)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				throw("listing %s: %w", name, err)
			}
			if !yield(name) {
				return
			}
		}

		for isa := range supportedTags {
			isa_arch := isa + `_` + dictISAtoArch[isa]
			names, err := fs.Glob(wd.FS(), "*_"+isa_arch+".go")
			if err != nil {
				throw("listing %s: %w", env.dstDir, err)
			}
			for _, name := range names {
				if !yield(name) {
					return
				}
			}
		}
	}
}

func (env *genEnv) scanSources() {
	wd, err := os.OpenRoot(env.srcDir)
	if err != nil {
		throw("scanning %s: %w", env.srcDir, err)
	}
	defer wd.Close()

	for path := range env.listFilesToScan() {
		f, err := wd.Open(path)
		if err != nil {
			throw("scanning %s: %w", path, err)
		}

		data, err := io.ReadAll(f)
		if err != nil {
			throw("scanning %s: %w", path, err)
		}
		f.Close()

		name := strings.TrimSuffix(path, "_simd.go")
		env.initSource(name, string(data))
	}
}

func (env *genEnv) commitFiles() {
	wd, err := os.OpenRoot(env.dstDir)
	if err != nil {
		throw("deleting in directory %s: %w", env.dstDir, err)
	}
	defer wd.Close()

	for path := range env.listFilesToDelete() {
		if err := wd.Remove(path); err != nil {
			throw("deleting %s: %w", path, err)
		}
	}

	for path, data := range env.listFilesToWrite() {
		f, err := wd.Create(path)
		if err != nil {
			throw("writing %s: %w", path, err)
		}

		if _, err := f.WriteString(data.String()); err != nil {
			throw("writing %s: %w", path, err)
		}
		f.Close()
	}
}
