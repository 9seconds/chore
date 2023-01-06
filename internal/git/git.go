package git

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

var ErrAccessModeUnknown = errors.New("access mode is unknown")

const (
	RefTypeTag    = "tag"
	RefTypeBranch = "branch"
	RefTypeRemote = "remote"
	RefTypeNote   = "note"
	RefTypeCommit = "commit"
)

type AccessMode byte

const (
	AccessNo AccessMode = iota
	AccessIfUndefined
	AccessAlways
)

func (m AccessMode) String() string {
	switch m {
	case AccessIfUndefined:
		return "if_undefined"
	case AccessNo:
		return "no"
	case AccessAlways:
		return "always"
	}

	return "<unknown>"
}

func GetAccessMode(value string) (AccessMode, error) {
	switch value {
	case "", "if_undefined":
		return AccessIfUndefined, nil
	case "no":
		return AccessNo, nil
	case "always":
		return AccessAlways, nil
	}

	return 0, ErrAccessModeUnknown
}

func GetRepo() (*git.Repository, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot find out current working dir: %w", err)
	}

	repo, err := git.PlainOpenWithOptions(currentDir, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot open git repository: %w", err)
	}

	return repo, nil
}
