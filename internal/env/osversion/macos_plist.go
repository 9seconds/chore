package osversion

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"howett.net/plist"
)

const (
	MacOSProductName = "ProductName"
	MacOSVersion     = "ProductUserVisibleVersion"
)

var MacOSCodeNames = map[uint64]string{
	11: "big sur",
	12: "monterey",
	13: "ventura",
}

func ParseMacOSPlist(path string) (OSVersion, error) {
	data := make(map[string]interface{})
	version := OSVersion{}

	file, err := os.Open(path)
	if err != nil {
		return version, fmt.Errorf("cannot open plist: %w", err)
	}

	defer file.Close()

	if err := plist.NewDecoder(file).Decode(data); err != nil {
		return version, fmt.Errorf("cannot parse plist: %w", err)
	}

	for _, param := range []string{MacOSProductName, MacOSVersion} {
		if _, ok := data[param]; !ok {
			return version, fmt.Errorf("cannot find out %s value", param)
		}

		if _, ok := data[param].(string); !ok {
			return version, fmt.Errorf("incorrect string value for %s", param)
		}
	}

	version.ID = strings.ToLower(data[MacOSProductName].(string))
	version.Version = data[MacOSVersion].(string)

	parsed, err := semver.NewVersion(version.Version)
	if err != nil {
		return version, fmt.Errorf("cannot parse version: %w", err)
	}

	version.Major = uint64(parsed.Major())
	version.Minor = uint64(parsed.Minor())

	if codename, ok := MacOSCodeNames[version.Major]; ok {
		version.Codename = codename
	} else {
		return version, fmt.Errorf("unknown major mac os version: %d", version.Major)
	}

	return version, nil
}
