package osversion

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"gopkg.in/ini.v1"
)

const (
	LinuxOSReleaseID        = "ID"
	LinuxOSReleaseVersionID = "VERSION_ID"
	LinuxOSReleaseCodename  = "VERSION_CODENAME"
)

func ParseLinuxOSRelease(path string) (OSVersion, error) {
	version := OSVersion{}

	cfg, err := ini.Load(path)
	if err != nil {
		return version, fmt.Errorf("cannot load ini file: %w", err)
	}

	section := cfg.Section("")

	for _, key := range []string{LinuxOSReleaseID, LinuxOSReleaseVersionID, LinuxOSReleaseCodename} {
		if !section.HasKey(key) {
			return version, fmt.Errorf("cannot find out %s value", key)
		}
	}

	version.ID = strings.ToLower(section.Key(LinuxOSReleaseID).MustString(""))
	version.Version = section.Key(LinuxOSReleaseVersionID).MustString("")
	version.Codename = strings.ToLower(section.Key(LinuxOSReleaseCodename).MustString(""))

	parsed, err := semver.NewVersion(version.Version)
	if err != nil {
		return version, fmt.Errorf("cannot parse version: %w", err)
	}

	version.Major = uint64(parsed.Major())
	version.Minor = uint64(parsed.Minor())

	return version, nil
}
