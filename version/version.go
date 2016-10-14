package version

import (
	"fmt"
)

// Version is the application's current version
var Version = newVersionNumber(0, 3, 1)

type semverNumber struct {
	major int
	minor int
	patch int
}

// String returns the version number as a string
func (v *semverNumber) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func newVersionNumber(major, minor, patch int) semverNumber {
	return semverNumber{major: major, minor: minor, patch: patch}
}
