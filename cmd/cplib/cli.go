package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/cplib"
	hl "github.com/mjwhitta/hilighter"
)

// Exit status
const (
	Good = iota
	InvalidOption
	MissingOption
	InvalidArgument
	MissingArgument
	ExtraArgument
	Exception
)

// Flags
var flags struct {
	append   bool
	exports  bool
	generate bool
	imports  bool
	libs     cli.StringList
	nocolor  bool
	output   string
	verbose  bool
	version  bool
}

func init() {
	// Configure cli package
	cli.Align = true // Defaults to false
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = filepath.Base(os.Args[0]) + " [OPTIONS] <binary>"
	cli.BugEmail = "cplib.bugs@whitta.dev"

	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		fmt.Sprintf("  %d: Invalid option\n", InvalidOption),
		fmt.Sprintf("  %d: Missing option\n", MissingOption),
		fmt.Sprintf("  %d: Invalid argument\n", InvalidArgument),
		fmt.Sprintf("  %d: Missing argument\n", MissingArgument),
		fmt.Sprintf("  %d: Extra argument\n", ExtraArgument),
		fmt.Sprintf("  %d: Exception", Exception),
	)
	cli.Info(
		"This tool can create a Go source file with a template for",
		"the exports/imports of the specified binary.",
	)

	cli.SeeAlso = []string{"koppeling"}
	cli.Title = "Copy Library"

	// Parse cli flags
	cli.Flag(
		&flags.append,
		"a",
		"append",
		false,
		"Append to source file, if it already exists.",
	)
	cli.Flag(
		&flags.exports,
		"e",
		"exports",
		false,
		"Only parse exports.",
	)
	cli.Flag(
		&flags.generate,
		"g",
		"generate",
		false,
		"Generate source for exports/imports.",
	)
	cli.Flag(
		&flags.imports,
		"i",
		"imports",
		false,
		"Only parse imports.",
	)
	cli.Flag(&flags.libs, "f", "filter", "Filter imports by library.")
	cli.Flag(
		&flags.nocolor,
		"no-color",
		false,
		"Disable colorized output.",
	)
	cli.Flag(
		&flags.output,
		"o",
		"output",
		".",
		"Write Go source to the specified directory.",
	)
	cli.Flag(
		&flags.verbose,
		"v",
		"verbose",
		false,
		"Show stacktrace, if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit, if version was requested
	if flags.version {
		fmt.Println(
			filepath.Base(os.Args[0]) + " version " + cplib.Version,
		)
		os.Exit(Good)
	}

	// Validate cli flags
	switch {
	case cli.NArg() < 1:
		cli.Usage(MissingArgument)
	case cli.NArg() > 1:
		cli.Usage(ExtraArgument)
	}

	// If neither, default to both
	if !flags.exports && !flags.imports {
		flags.exports = true
		flags.imports = true
	}

	// Normalize to lowercase
	for i := range flags.libs {
		flags.libs[i] = strings.ToLower(flags.libs[i])
	}
}
