package main

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/cplib"
	"github.com/mjwhitta/log"
)

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

	var bin string
	var e error
	var exports []string
	var fn string
	var imports []cplib.Import
	var keep []cplib.Import
	var tags string

	validate()

	bin = cli.Arg(0)

	switch strings.ToLower(filepath.Ext(bin)) {
	case "", ".bin":
		fn = "generated_unix.go"
		tags = "!windows"
		imports, e = cplib.ELFImports(bin)
	case ".dll":
		fn = "generated_windows.go"
		tags = "dll && windows"
		exports, e = cplib.PEExports(bin)
	case ".exe":
		fn = "generated_windows.go"
		tags = "dll && windows"
		imports, e = cplib.PEImports(bin)
	case ".so":
		fn = "generated_unix.go"
		tags = "!windows"
		exports, e = cplib.ELFExports(bin)
	}

	if e != nil {
		panic(e)
	}

	// Filter imports
	for _, im := range imports {
		if len(flags.libs) > 0 {
			if !slices.Contains(flags.libs, strings.ToLower(im.Lib)) {
				continue
			}
		}

		keep = append(keep, im)
	}

	fn = filepath.Join(flags.output, fn)
	if e = cplib.GenerateGo(bin, fn, tags, exports, keep); e != nil {
		panic(e)
	}

	log.Goodf("Source written to %s", fn)
}
