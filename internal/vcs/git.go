package vcs

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

var ErrGitAccessModeUnknown = errors.New("access mode is unknown")

const (
	GitRefTypeTag    = "tag"
	GitRefTypeBranch = "branch"
	GitRefTypeRemote = "remote"
	GitRefTypeNote   = "note"
	GitRefTypeCommit = "commit"
)

type GitAccessMode byte

const (
	GitAccessNo GitAccessMode = iota
	GitAccessIfUndefined
	GitAccessAlways
)

func (m GitAccessMode) String() string {
	switch m {
	case GitAccessIfUndefined:
		return "if_undefined"
	case GitAccessNo:
		return "no"
	case GitAccessAlways:
		return "always"
	}

	return "<unknown>"
}

func GetGitAccessMode(value string) (GitAccessMode, error) {
	switch value {
	case "", "if_undefined":
		return GitAccessIfUndefined, nil
	case "no":
		return GitAccessNo, nil
	case "always":
		return GitAccessAlways, nil
	}

	return 0, ErrGitAccessModeUnknown
}

func GetGitRepo() (*git.Repository, error) {
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
