package cplib

import (
	"debug/elf"
	"slices"

	"github.com/mjwhitta/errors"
)

// ELFExports will read the provided filename and, if it's an ELF, it
// will return a list of functions that are exported.
func ELFExports(_ string) ([]string, error) {
	return nil, errors.New("not implemented")
}

// ELFImports will read the provided filename and, if it's an ELF, it
// will return a list of functions that are imported.
func ELFImports(fn string) ([]Import, error) {
	var exe *elf.File
	var e error
	var imports []Import
	var entries []elf.ImportedSymbol

	if exe, e = elf.Open(fn); e != nil {
		return nil, errors.Newf("failed to open DLL: %w", e)
	}

	if entries, e = exe.ImportedSymbols(); e != nil {
		return nil, errors.Newf("failed to read imports: %w", e)
	}

	slices.SortFunc(entries, sortELFImportsCaseInsensitive)

	for _, entry := range entries {
		imports = append(
			imports,
			Import{Lib: entry.Library, Name: entry.Name},
		)
	}

	_ = imports
	// return imports, nil

	return nil, errors.New("not implemented")
}
