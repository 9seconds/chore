package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/9seconds/chore/internal/git"
)

const ParameterGit = "git"

var ErrGitIncorrectRefType = errors.New("value does not conform an expected reference type")

type parameterGit struct {
	baseParameter

	repo     *git.Repo
	refTypes []git.RefType
}

func (p parameterGit) Type() string {
	return ParameterGit
}

func (p parameterGit) Validate(_ context.Context, value string) error { //nolint: cyclop
	refTypes := p.refTypes

	if len(p.refTypes) == 0 {
		if ok, err := p.repo.HasRevision(value); err != nil || ok {
			return err
		}

		return ErrGitIncorrectRefType
	}

	for _, refType := range refTypes {
		var validate func(string) (bool, error)

		switch refType {
		case git.RefTypeTag:
			validate = p.repo.HasTag
		case git.RefTypeBranch:
			validate = p.repo.HasBranch
		case git.RefTypeRemote:
			validate = p.repo.HasRemote
		case git.RefTypeNote:
			validate = p.repo.HasNote
		case git.RefTypeCommit:
			validate = p.repo.HasCommit
		default:
			panic(fmt.Sprintf("unexpected reference type %v", refType))
		}

		if ok, err := validate(value); err != nil || ok {
			return err
		}
	}

	return ErrGitIncorrectRefType
}

func NewGit(
	description string,
	required bool,
	spec map[string]string,
	createRepo func() (*git.Repo, error),
) (Parameter, error) {
	param := parameterGit{
		baseParameter: baseParameter{
			required:      required,
			description:   description,
			specification: spec,
		},
	}

	refTypes := make(map[git.RefType]bool)

	for _, value := range parseCSV(spec["type"]) {
		refType, err := git.GetRefType(value)
		if err != nil {
			return nil, fmt.Errorf("incorrect value %s for 'type': %w", refType, err)
		}

		refTypes[refType] = true
	}

	for k := range refTypes {
		param.refTypes = append(param.refTypes, k)
	}

	repo, err := createRepo()
	if err != nil {
		return nil, fmt.Errorf("cannot initialize repo: %w", err)
	}

	param.repo = repo

	return param, nil
}
