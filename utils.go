package cplib

import (
	"debug/elf"
	"debug/pe"
	"encoding/binary"
	"strings"

	"github.com/mjwhitta/errors"
)

func findSection(
	addr uint32,
	sections []*pe.Section,
) (*pe.Section, error) {
	var e error
	var start uint32
	var stop uint32

	for _, section := range sections {
		start = section.VirtualAddress
		stop = start + section.VirtualSize

		if (addr >= start) && (addr <= stop) {
			return section, nil
		}
	}

	e = errors.Newf("failed to find section for address 0x%08x", addr)

	return nil, e
}

func leUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func leUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func sortELFImportsCaseInsensitive(a, b elf.ImportedSymbol) int {
	var l string = strings.ToLower(a.Name)
	var r string = strings.ToLower(b.Name)

	switch {
	case l < r:
		return -1
	case l > r:
		return 1
	default:
		return 0
	}
}

func sortStringsCaseInsensitive(a string, b string) int {
	var l string = strings.ToLower(a)
	var r string = strings.ToLower(b)

	switch {
	case l < r:
		return -1
	case l > r:
		return 1
	default:
		return 0
	}
}
