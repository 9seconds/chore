//go:build !linux && !darwin

package osversion

func Get() (OSVersion, error) {
	return OSVersion{}, ErrOSVersionNotSupported
}
