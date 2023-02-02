package git2

import "errors"

type RefType string

const (
	RefTypeRevision RefType = "revision"
	RefTypeTag      RefType = "tag"
	RefTypeBranch   RefType = "branch"
	RefTypeRemote   RefType = "remote"
	RefTypeNote     RefType = "note"
	RefTypeCommit   RefType = "commit"
)

var (
	ErrInvalidRefType = errors.New("invalid reftype")
)

func (r RefType) String() string {
	return string(r)
}

func (r RefType) Valid() bool {
	switch string(r) {
	case "tag", "branch", "remote", "note", "commit", "revision":
		return true
	}

	return false
}

func GetRefType(value string) (RefType, error) {
	if value == "" {
		value = RefTypeRevision.String()
	}

	ref := RefType(value)

	if !ref.Valid() {
		return "", ErrInvalidRefType
	}

	return ref, nil
}
