//go:build darwin

package osversion

func Get() (OSVersion, error) {
	return ParseMacOSPlist("/System/Library/CoreServices/SystemVersion.plist")
}
