package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type BumpType string

const (
	Major BumpType = "major"
	Minor BumpType = "minor"
	Patch BumpType = "patch"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`)

func Parse(s string) (*Version, error) {
	s = strings.TrimSpace(s)
	matches := semverRegex.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("invalid semver format: %s", s)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func (v *Version) Bump(t BumpType) *Version {
	switch t {
	case Major:
		return &Version{Major: v.Major + 1, Minor: 0, Patch: 0}
	case Minor:
		return &Version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}
	case Patch:
		return &Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
	default:
		return &Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
	}
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func IsBumpType(s string) bool {
	s = strings.ToLower(s)
	return s == "major" || s == "minor" || s == "patch"
}

func ParseBumpType(s string) (BumpType, error) {
	s = strings.ToLower(s)
	switch s {
	case "major":
		return Major, nil
	case "minor":
		return Minor, nil
	case "patch":
		return Patch, nil
	default:
		return "", fmt.Errorf("invalid bump type: %s (must be major, minor, or patch)", s)
	}
}
