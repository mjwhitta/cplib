package cplib

import (
	"debug/pe"
	"slices"

	"github.com/mjwhitta/errors"
)

// ExportTable is a struct containing relevant export data for a PE
// file.
type ExportTable struct {
	AddressOfFuncs        uint32
	AddressOfNameOrdinals uint32
	AddressOfNames        uint32
	Base                  uint32
	Characteristics       uint32
	Checksum              uint32
	MajorVersion          uint16
	MinorVersion          uint16
	Name                  uint32
	NumberOfFuncs         uint32
	NumberOfNames         uint32

	exportOffset  uint32
	functions     map[string]uint32
	pe            *pe.File
	sectionData   []byte
	sectionOffset uint32
}

// GetExportTable will return the export table for the provided PE
// file.
func GetExportTable(pf *pe.File) (*ExportTable, error) {
	var b []byte
	var dataDir pe.DataDirectory
	var e error
	var nameOff uint32
	var s *pe.Section
	var table *ExportTable

	// Read header for export data directory
	switch hdr := pf.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		dataDir = hdr.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_EXPORT]
	case *pe.OptionalHeader64:
		dataDir = hdr.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_EXPORT]
	default:
		return nil, errors.New("invalid optional header format")
	}

	// Get nameOff for export table
	s, nameOff, e = sectionOffset(dataDir.VirtualAddress, pf.Sections)
	if e != nil {
		return nil, e
	}

	if b, e = s.Data(); e != nil {
		return nil, errors.Newf("failed to read section data: %w", e)
	}

	if len(b) < int(nameOff)+40 {
		return nil, errors.New("truncated export table")
	}

	table = &ExportTable{
		AddressOfFuncs:        leUint32(b[nameOff+28 : nameOff+32]),
		AddressOfNameOrdinals: leUint32(b[nameOff+36 : nameOff+40]),
		AddressOfNames:        leUint32(b[nameOff+32 : nameOff+36]),
		Base:                  leUint32(b[nameOff+16 : nameOff+20]),
		Characteristics:       leUint32(b[nameOff+0 : nameOff+4]),
		Checksum:              leUint32(b[nameOff+4 : nameOff+8]),
		exportOffset:          dataDir.VirtualAddress,
		functions:             map[string]uint32{},
		MajorVersion:          leUint16(b[nameOff+8 : nameOff+10]),
		MinorVersion:          leUint16(b[nameOff+10 : nameOff+12]),
		Name:                  leUint32(b[nameOff+12 : nameOff+16]),
		NumberOfFuncs:         leUint32(b[nameOff+20 : nameOff+24]),
		NumberOfNames:         leUint32(b[nameOff+24 : nameOff+28]),
		pe:                    pf,
		sectionData:           b,
		sectionOffset:         s.Offset,
	}

	if e = table.parse(); e != nil {
		return nil, e
	}

	return table, nil
}

// Names will return a list of exported function names.
func (et *ExportTable) Names() []string {
	var names []string

	for k := range et.functions {
		names = append(names, k)
	}

	slices.SortFunc(names, sortStringsCaseInsensitive)

	return names
}

func (et *ExportTable) parse() error {
	var addr uint32
	var e error
	var name []byte
	var nameOff uint32
	var ordinal uint16
	var ordOff uint32
	var start uint32
	var stop uint32

	if et.NumberOfNames == 0 {
		return nil
	}

	_, nameOff, e = sectionOffset(et.AddressOfNames, et.pe.Sections)
	if e != nil {
		return e
	}

	_, ordOff, e = sectionOffset(
		et.AddressOfNameOrdinals,
		et.pe.Sections,
	)
	if e != nil {
		return e
	}

	for i := range et.NumberOfNames {
		start = nameOff + 4*i //nolint:mnd // 4 is size of uint32

		if int(start)+4 > len(et.sectionData) {
			return errors.New("export offset is out of range")
		}

		addr = leUint32(et.sectionData[start : start+4])

		_, start, e = sectionOffset(addr, et.pe.Sections)
		if e != nil {
			return e
		}

		if int(start) > len(et.sectionData) {
			return errors.New("export name offset is out of range")
		}

		for stop = start; et.sectionData[stop] != 0; stop++ {
		}

		name = et.sectionData[start:stop]

		start = ordOff + 2*i //nolint:mnd // 2 is size of uint16
		stop = start + 2     //nolint:mnd // 2 is size of uint16

		if int(stop) > len(et.sectionData) {
			e = errors.New("export oridinal offset is out of range")
			return e
		}

		ordinal = leUint16(et.sectionData[start:stop])

		et.functions[string(name)] = uint32(ordinal) + et.Base
	}

	return nil
}
