package osversion

import "errors"

var ErrOSVersionNotSupported = errors.New("os version introspection is not supported")

type OSVersion struct {
	ID       string
	Version  string
	Codename string
	Major    uint64
	Minor    uint64
}
