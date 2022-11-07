package sem

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	Invalid = iota
	Major
	Minor
	Patch
)

// ErrInvalidKind is returned when a Kind is not within the accepted constants
var ErrInvalidKind = errors.New("Kind is invalid, needs to be either patch, minor, major")

var ErrParseVersionFault = errors.New("given string could not be parsed into valid semver version")

type Ver struct {
	Prefix string
	Major  int
	Minor  int
	Patch  int
}

func (v *Ver) String() string {
	return fmt.Sprintf("%s%d.%d.%d", v.Prefix, v.Major, v.Minor, v.Patch)
}

type ByLatest []*Ver

func (b ByLatest) Len() int {
	return len(b)
}

func (b ByLatest) Less(i, j int) bool {
	if b[i].Major < b[j].Major {
		return true
	}
	if b[i].Major > b[j].Major {
		return false
	}

	if b[i].Minor < b[j].Minor {
		return true
	}
	if b[i].Minor > b[j].Minor {
		return false
	}

	if b[i].Patch < b[j].Patch {
		return true
	}
	if b[i].Patch > b[j].Patch {
		return false
	}

	// Equal has to return false
	return false
}

func (b ByLatest) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Kind denotes specific SemVer versions that the API accepts
type Kind int

func (vk *Kind) IsValid() error {
	if *vk < Major || *vk > Patch {
		return ErrInvalidKind
	}
	return nil
}

// ParseKind attempts to create a kind constant from the input string
func ParseKind(input string) (Kind, error) {
	var kind Kind
	switch strings.ToLower(input) {
	case "major":
		kind = Major
	case "minor":
		kind = Minor
	case "patch":
		kind = Patch
	}
	if err := kind.IsValid(); err != nil {
		return Invalid, err
	}
	return kind, nil
}

// Next accepts a Kind to return the next available SemVer version tag
func (v *Ver) Next(kind Kind) error {
	if err := kind.IsValid(); err != nil {
		return err
	}
	switch kind {
	case Major:
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case Minor:
		v.Minor++
		v.Patch = 0
	case Patch:
		v.Patch++
	}
	return nil
}

var semverRegexp = regexp.MustCompile(`^([a-zA-Z]*)(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)

// ParseVersion attempts to parse the given string into a Version object.
// This will only work with valid semver versions like 1.0.0 or v1.0.0
// The prefix matches /[a-zA-Z]/ so no emojis! :(
func ParseVersion(version string) (*Ver, error) {
	versionRegexMatch := semverRegexp.FindStringSubmatch(version)
	if versionRegexMatch == nil || len(versionRegexMatch) != 5 {
		return nil, ErrParseVersionFault
	}

	majorVersionInt, err := strconv.Atoi(versionRegexMatch[2])
	if err != nil {
		return nil, ErrParseVersionFault
	}

	minorVersionInt, err := strconv.Atoi(versionRegexMatch[3])
	if err != nil {
		return nil, ErrParseVersionFault
	}

	patchVersionInt, err := strconv.Atoi(versionRegexMatch[4])
	if err != nil {
		return nil, ErrParseVersionFault
	}

	ver := &Ver{
		Prefix: versionRegexMatch[1],
		Major:  majorVersionInt,
		Minor:  minorVersionInt,
		Patch:  patchVersionInt,
	}

	return ver, nil
}

// GetAllVersions returns a sorted list of all available git tags, if they all are SemVer compliant
func GetAllVersions(dir string, skipParseErrors bool) ([]*Ver, error) {
	if err := os.Chdir(dir); err != nil {
		return nil, err
	}
	cmd := exec.Command("git", "tag")
	cmdOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	tagList := strings.Split(string(cmdOut), "\n")
	versionList := make([]*Ver, 0)
	for _, tag := range tagList {
		if tag == "" {
			continue
		}
		version, err := ParseVersion(tag)
		if err != nil {
			if skipParseErrors {
				continue
			}
			return nil, err
		}
		versionList = append(versionList, version)
	}
	sort.Sort(ByLatest(versionList))
	return versionList, nil
}

var ErrNoVersionsAvailable = errors.New("no versions available")

// GetLatestVersion gets all versions and returns the latest one
func GetLatestVersion(dir string, ignoreNonSemVerTags bool) (*Ver, error) {
	versionList, err := GetAllVersions(dir, ignoreNonSemVerTags)
	if err != nil {
		return nil, err
	}
	if len(versionList) == 0 {
		return nil, ErrNoVersionsAvailable
	}
	if len(versionList) == 1 {
		return versionList[0], nil
	}
	return versionList[len(versionList)-1], nil
}
