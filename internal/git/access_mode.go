package git

import "errors"

type AccessMode string

const (
	AccessModeNo          AccessMode = "no"
	AccessModeIfUndefined AccessMode = "if_undefined"
	AccessModeAlways      AccessMode = "always"
)

var (
	ErrInvalidAccessMode = errors.New("invalid mode")
)

func (a AccessMode) String() string {
	return string(a)
}

func (a AccessMode) Valid() bool {
	switch string(a) {
	case "if_undefined", "no", "always":
		return true
	}

	return false
}

func GetAccessMode(value string) (AccessMode, error) {
	if value == "" {
		value = AccessModeNo.String()
	}

	mode := AccessMode(value)

	if !mode.Valid() {
		return "", ErrInvalidAccessMode
	}

	return mode, nil
}
