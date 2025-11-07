package cplib

import (
	"debug/pe"
	"slices"
	"strings"

	"github.com/mjwhitta/errors"
)

// PEExports will read the provided filename and, if it's a PE, it
// will return a list of functions that are exported.
func PEExports(fn string) ([]string, error) {
	var dll *pe.File
	var e error
	var table *ExportTable

	if dll, e = pe.Open(fn); e != nil {
		return nil, errors.Newf("failed to open DLL: %w", e)
	}

	if table, e = GetExportTable(dll); e != nil {
		return nil, e
	}

	return table.Names(), nil
}

// PEImports will read the provided filename and, if it's a PE, it
// will return a list of functions that are imported.
func PEImports(fn string) ([]Import, error) {
	var exe *pe.File
	var e error
	var imports []Import
	var entries []string

	if exe, e = pe.Open(fn); e != nil {
		return nil, errors.Newf("failed to open DLL: %w", e)
	}

	if entries, e = exe.ImportedSymbols(); e != nil {
		return nil, errors.Newf("failed to read imports: %w", e)
	}

	slices.SortFunc(entries, sortStringsCaseInsensitive)

	for _, entry := range entries {
		if name, lib, ok := strings.Cut(entry, ":"); ok {
			imports = append(imports, Import{Lib: lib, Name: name})
		}
	}

	return imports, nil
}
