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
	var bin *pe.File
	var e error
	var table *ExportTable

	if bin, e = pe.Open(fn); e != nil {
		return nil, errors.Newf("failed to open binary: %w", e)
	}

	if table, e = GetExportTable(bin); e != nil {
		return nil, e
	}

	return table.Names(), nil
}

// PEImports will read the provided filename and, if it's a PE, it
// will return a list of functions that are imported.
func PEImports(fn string) ([]Import, error) {
	var bin *pe.File
	var e error
	var imports []Import
	var entries []string

	if bin, e = pe.Open(fn); e != nil {
		return nil, errors.Newf("failed to open binary: %w", e)
	}

	if entries, e = bin.ImportedSymbols(); e != nil {
		return nil, errors.Newf("failed to read imports: %w", e)
	}

	slices.SortFunc(entries, sortStringsCaseInsensitive)

	for _, entry := range entries {
		if name, lib, ok := strings.Cut(entry, ":"); ok {
			imports = append(
				imports,
				Import{Library: lib, Name: name},
			)
		}
	}

	return imports, nil
}
