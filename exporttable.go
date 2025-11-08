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

	functions    map[string]uint32
	sectionData  []byte
	sectionStart uint32
}

// GetExportTable will return the export table for the provided PE
// file.
func GetExportTable(pf *pe.File) (*ExportTable, error) {
	var b []byte
	var dataDir pe.DataDirectory
	var e error
	var offset uint32
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

	if dataDir.VirtualAddress == 0 {
		return &ExportTable{}, nil
	}

	// Find section for export table
	s, e = findSection(dataDir.VirtualAddress, pf.Sections)
	if e != nil {
		return nil, e
	}

	// Get offset for export table
	offset = dataDir.VirtualAddress - s.VirtualAddress

	if b, e = s.Data(); e != nil {
		return nil, errors.Newf("failed to read section data: %w", e)
	}

	if len(b) < int(offset)+40 {
		return nil, errors.New("truncated export table")
	}

	table = &ExportTable{
		AddressOfFuncs:        leUint32(b[offset+28 : offset+32]),
		AddressOfNameOrdinals: leUint32(b[offset+36 : offset+40]),
		AddressOfNames:        leUint32(b[offset+32 : offset+36]),
		Base:                  leUint32(b[offset+16 : offset+20]),
		Characteristics:       leUint32(b[offset+0 : offset+4]),
		Checksum:              leUint32(b[offset+4 : offset+8]),
		functions:             map[string]uint32{},
		MajorVersion:          leUint16(b[offset+8 : offset+10]),
		MinorVersion:          leUint16(b[offset+10 : offset+12]),
		Name:                  leUint32(b[offset+12 : offset+16]),
		NumberOfFuncs:         leUint32(b[offset+20 : offset+24]),
		NumberOfNames:         leUint32(b[offset+24 : offset+28]),
		sectionData:           b,
		sectionStart:          s.VirtualAddress,
	}

	if e = table.parse(); e != nil {
		return nil, e
	}

	return table, nil
}

// Names will return a list of exported function names.
func (et *ExportTable) Names() []string {
	var names []string

	if et.functions == nil {
		return nil
	}

	for k := range et.functions {
		names = append(names, k)
	}

	slices.SortFunc(names, sortStringsCaseInsensitive)

	return names
}

func (et *ExportTable) parse() error {
	var e error
	var name []byte
	var namePtr uint32
	var namesOff uint32
	var ordinal uint16
	var ordsOff uint32
	var start uint32
	var stop uint32

	if et.NumberOfNames == 0 {
		return nil
	}

	namesOff = et.AddressOfNames - et.sectionStart
	ordsOff = et.AddressOfNameOrdinals - et.sectionStart

	for i := range et.NumberOfNames {
		// Get export name pointer offset
		start = namesOff + 4*i //nolint:mnd // 4 is size of uint32
		stop = start + 4       //nolint:mnd // 4 is size of uint32

		if int(stop) > len(et.sectionData) {
			return errors.New("export is out of range")
		}

		// Get pointer to name
		namePtr = leUint32(et.sectionData[start:stop])

		// Get export name from pointer
		start = namePtr - et.sectionStart

		// Find null character
		for stop = start; ; stop++ {
			if int(stop) > len(et.sectionData) {
				return errors.New("export name is out of range")
			} else if et.sectionData[stop] == 0 {
				break
			}
		}

		// Extract export name
		name = et.sectionData[start:stop]

		// Get ordinal offset
		start = ordsOff + 2*i //nolint:mnd // 2 is size of uint16
		stop = start + 2      //nolint:mnd // 2 is size of uint16

		if int(stop) > len(et.sectionData) {
			e = errors.New("export oridinal is out of range")
			return e
		}

		// Extract export ordinal
		ordinal = leUint16(et.sectionData[start:stop])

		// Store function
		et.functions[string(name)] = uint32(ordinal) + et.Base
	}

	return nil
}
