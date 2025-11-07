package main

import (
	"os"
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
	libs    cli.StringList
	nocolor bool
	output  string
	verbose bool
	version bool
}

func init() {
	// Configure cli package
	cli.Align = true // Defaults to false
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = os.Args[0] + " [OPTIONS] <dll/exe>"
	cli.BugEmail = "cplib.bugs@whitta.dev"

	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		hl.Sprintf("  %d: Invalid option\n", InvalidOption),
		hl.Sprintf("  %d: Missing option\n", MissingOption),
		hl.Sprintf("  %d: Invalid argument\n", InvalidArgument),
		hl.Sprintf("  %d: Missing argument\n", MissingArgument),
		hl.Sprintf("  %d: Extra argument\n", ExtraArgument),
		hl.Sprintf("  %d: Exception", Exception),
	)
	cli.Info(
		"This tool will create Go source files with a template for",
		"either the exports of the input library or the imports of",
		"the input executable.",
	)

	cli.SeeAlso = []string{"koppling"}
	cli.Title = "Copy Library"

	// Parse cli flags
	cli.Flag(
		&flags.libs,
		"l",
		"lib",
		"Filter imports to only include specified libraries.",
	)
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
		hl.Printf("cplib version %s\n", cplib.Version)
		os.Exit(Good)
	}

	// Validate cli flags
	switch {
	case cli.NArg() < 1:
		cli.Usage(MissingArgument)
	case cli.NArg() > 1:
		cli.Usage(ExtraArgument)
	}

	// Normalize to lowercase
	for i := range flags.libs {
		flags.libs[i] = strings.ToLower(flags.libs[i])
	}
}
