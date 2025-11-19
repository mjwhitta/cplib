package cplib

import "regexp"

// Version is the package version.
const Version = "1.2.3"

var sharedObject *regexp.Regexp = regexp.MustCompile(`.+\.so\.\d+`)
