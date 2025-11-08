package cplib

import "regexp"

// Version is the package version.
const Version = "1.1.1"

var sharedObject *regexp.Regexp = regexp.MustCompile(`.+\.so\.\d+`)
