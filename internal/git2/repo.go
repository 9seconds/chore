package git2

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Repo struct {
	repo *git.Repository

	head    *plumbing.Reference
	isDirty bool

	branches map[string]bool
	remotes  map[string]bool
	notes    map[string]bool
	tags     map[string]bool

	collectReferences func() (bool, error)
	collectIsDirty    func() (bool, error)
	collectHead       func() (*plumbing.Reference, error)
}

func (r *Repo) HasBranch(name string) (bool, error) {
	if _, err := r.collectReferences(); err != nil {
		return false, err
	}

	return r.branches[name], nil
}

func (r *Repo) HasRemote(name string) (bool, error) {
	if _, err := r.collectReferences(); err != nil {
		return false, err
	}

	return r.remotes[name], nil
}

func (r *Repo) HasNote(name string) (bool, error) {
	if _, err := r.collectReferences(); err != nil {
		return false, err
	}

	return r.notes[name], nil
}

func (r *Repo) HasTag(name string) (bool, error) {
	if _, err := r.collectReferences(); err != nil {
		return false, err
	}

	return r.tags[name], nil
}

func (r *Repo) HasRevision(rev string) (ok bool, err error) {
	// go-git is horrible here: https://github.com/go-git/go-git/issues/674
	// seems very little tested
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic: %v", rec)
			ok = false
		}
	}()

	hash, err := r.repo.ResolveRevision(plumbing.Revision(rev))

	switch {
	case errors.Is(err, plumbing.ErrReferenceNotFound):
		return false, nil
	// https://github.com/go-git/go-git/issues/673
	case err != nil && strings.Contains(err.Error(), "Revision invalid"):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("cannot resolve revision: %w", err)
	case hash != &plumbing.ZeroHash, hash != nil && *hash != plumbing.ZeroHash:
		return true, nil
	}

	return r.HasBranch(rev)
}

func (r *Repo) Head() (*plumbing.Reference, error) {
	return r.collectHead()
}

func (r *Repo) IsDirty() (bool, error) {
	return r.collectIsDirty()
}

func New() (*Repo, error) {
	gitDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot find out current working dir: %w", err)
	}

	if value, ok := os.LookupEnv("GIT_DIR"); ok {
		gitDir = value
	}

	repo, err := git.PlainOpenWithOptions(gitDir, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot open git repository: %w", err)
	}

	value := &Repo{
		repo:     repo,
		branches: make(map[string]bool),
		remotes:  make(map[string]bool),
		notes:    make(map[string]bool),
		tags:     make(map[string]bool),
	}

	value.collectReferences = memoize(func() (bool, error) {
		iter, err := repo.Storer.IterReferences()
		if err != nil {
			return false, err
		}

		defer iter.Close()

		return false, iter.ForEach(func(ref *plumbing.Reference) error {
			name := ref.Name()
			short := name.Short()

			switch {
			case name.IsBranch():
				value.branches[short] = true
			case name.IsTag():
				value.tags[short] = true
			case name.IsNote():
				value.notes[short] = true
			case name.IsRemote():
				value.remotes[short] = true
			}

			return nil
		})
	})
	value.collectHead = memoize(repo.Head)
	value.collectIsDirty = memoize(func() (bool, error) {
		tree, err := repo.Worktree()
		if err != nil {
			return false, fmt.Errorf("cannot get worktree: %w", err)
		}

		status, err := tree.Status()
		if err != nil {
			return false, fmt.Errorf("cannot get status: %w", err)
		}

		return !status.IsClean(), nil
	})

	return value, nil
}

var (
	globalRepo     *Repo
	globalError    error
	globalMakeSync sync.Once
)

func Get() (*Repo, error) {
	globalMakeSync.Do(func() {
		globalRepo, globalError = New()
	})

	return globalRepo, globalError
}
