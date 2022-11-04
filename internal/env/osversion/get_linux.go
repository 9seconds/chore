//go:build linux

package osversion

func Get() (OSVersion, error) {
	return ParseLinuxOSRelease("/etc/os-release")
}
