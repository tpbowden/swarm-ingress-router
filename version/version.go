package version

import (
	"fmt"
)

var Version = newVersionNumber(0, 1, 0)

type semverNumber struct {
	major int
	minor int
	patch int
}

func (v *semverNumber) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func newVersionNumber(major int, minor int, patch int) semverNumber {
	return semverNumber{major: major, minor: minor, patch: patch}
}
