package cplib

import "regexp"

// Version is the package version.
const Version = "1.2.4"

var sharedObject *regexp.Regexp = regexp.MustCompile(`.+\.so\.\d+`)
