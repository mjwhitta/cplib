package cplib

import (
	"debug/elf"
	"slices"
	"strings"

	"github.com/mjwhitta/errors"
)

// ELFExports will read the provided filename and, if it's an ELF, it
// will return a list of functions that are exported.
func ELFExports(fn string) ([]string, error) {
	var bin *elf.File
	var e error
	var exports []string
	var index int
	var syms []elf.Symbol

	if bin, e = elf.Open(fn); e != nil {
		return nil, errors.Newf("failed to open binary: %w", e)
	}

	for i, s := range bin.Sections {
		if s.Name == ".text" {
			index = i
		}
	}

	if syms, e = bin.DynamicSymbols(); e != nil {
		return nil, errors.Newf("failed to parse exports: %w", e)
	}

	for _, sym := range syms {
		switch {
		case sym.Library != "": // Must be local
			continue
		case sym.Section != elf.SectionIndex(index): // .text section
			continue
		case sym.VersionIndex.Index() == 0: // Not private
			continue
		case strings.HasPrefix(sym.Name, "_cgo"): // Not cgo
			continue
		case strings.HasPrefix(sym.Name, "x_cgo"): // Not cgo
			continue
		}

		exports = append(exports, sym.Name)
	}

	slices.SortFunc(exports, sortStringsCaseInsensitive)

	return exports, nil
}

// ELFImports will read the provided filename and, if it's an ELF, it
// will return a list of functions that are imported.
func ELFImports(fn string) ([]Import, error) {
	var bin *elf.File
	var e error
	var imports []Import
	var entries []elf.ImportedSymbol

	if bin, e = elf.Open(fn); e != nil {
		return nil, errors.Newf("failed to open binary: %w", e)
	}

	if entries, e = bin.ImportedSymbols(); e != nil {
		return nil, errors.Newf("failed to read imports: %w", e)
	}

	slices.SortFunc(entries, sortELFImportsCaseInsensitive)

	for _, entry := range entries {
		imports = append(
			imports,
			Import{Lib: entry.Library, Name: entry.Name},
		)
	}

	return imports, nil
}
