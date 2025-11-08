package cplib

import (
	"path/filepath"
	"strings"
)

// FileType is an enum
type FileType uint32

// FiltType consts
const (
	TypeUnknown = iota
	TypeELF
	TypePE
)

// NaiveFileType does some very naive checks to determine file type.
func NaiveFileType(bin string) FileType {
	var ext string = strings.ToLower(filepath.Ext(bin))

	switch ext {
	case "", ".bin", ".so": // Linux
		return TypeELF
	case ".dll", ".exe": // Windows
		return TypePE
	}

	if sharedObject.MatchString(bin) {
		return TypeELF // Linux shared objects
	}

	return TypeUnknown
}
