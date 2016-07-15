package version

import (
	"fmt"
)

var Version VersionNumber = NewVersionNumber(0, 0, 7)

type VersionNumber struct {
	major int
	minor int
	patch int
}

func (v *VersionNumber) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func NewVersionNumber(major int, minor int, patch int) VersionNumber {
	return VersionNumber{major: major, minor: minor, patch: patch}
}
