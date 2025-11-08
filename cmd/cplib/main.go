package main

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/cplib"
	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/pathname"
)

type state struct {
	bin     string
	exports []string
	fn      string
	imports []cplib.Import
	tags    string
}

func list(s *state) {
	if len(s.exports) > 0 {
		fmt.Printf("[*] %s exports\n", s.bin)
	}

	for i := range s.exports {
		fmt.Println(s.exports[i])
	}

	if len(s.imports) > 0 {
		fmt.Printf("[*] %s imports\n", s.bin)
	}

	for i := range s.imports {
		fmt.Printf("%s (%s)\n", s.imports[i].Name, s.imports[i].Lib)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r)
			}

			switch r := r.(type) {
			case error:
				log.ErrX(Exception, r.Error())
			case string:
				log.ErrX(Exception, r)
			}
		}
	}()

	var e error
	var s *state

	validate()

	if s, e = setup(cli.Arg(0)); e != nil {
		panic(e)
	}

	switch {
	case flags.generate:
		// If the file does not exist, no need to append
		if ok, e := pathname.DoesExist(s.fn); (e != nil) || !ok {
			flags.append = false
		}

		e = cplib.GenerateGo(
			s.bin,
			s.fn,
			s.tags,
			s.exports,
			s.imports,
			flags.append,
		)
		if e != nil {
			panic(e)
		}

		log.Goodf("Source written to %s", s.fn)
	default:
		list(s)
	}
}

//nolint:wrapcheck // I'm not wrapping my own errors
func setup(bin string) (*state, error) {
	var e error
	var keep []cplib.Import
	var s *state = &state{bin: bin}

	switch cplib.NaiveFileType(bin) {
	case cplib.TypeELF:
		s.fn = filepath.Join(flags.output, "generated_linux_so.go")
		s.tags = "linux && so"

		if flags.exports {
			if s.exports, e = cplib.ELFExports(bin); e != nil {
				return nil, e
			}
		}

		if flags.imports {
			if s.imports, e = cplib.ELFImports(bin); e != nil {
				return nil, e
			}
		}
	case cplib.TypePE:
		s.fn = filepath.Join(flags.output, "generated_windows_dll.go")
		s.tags = "dll && windows"

		if flags.exports {
			if s.exports, e = cplib.PEExports(bin); e != nil {
				return nil, e
			}
		}

		if flags.imports {
			if s.imports, e = cplib.PEImports(bin); e != nil {
				return nil, e
			}
		}
	default:
		return nil, errors.Newf("unsupported filetype: %s", bin)
	}

	// Filter imports
	for _, im := range s.imports {
		if len(flags.libs) > 0 {
			if !slices.Contains(flags.libs, strings.ToLower(im.Lib)) {
				continue
			}
		}

		keep = append(keep, im)
	}

	s.imports = keep

	return s, nil
}
