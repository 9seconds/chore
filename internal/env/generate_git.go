package env

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/9seconds/chore/internal/git"
)

const GitCommitHashShortLength = 12

func GenerateGit(ctx context.Context, results chan<- string, waiters *sync.WaitGroup, mode git.AccessMode) { //nolint: cyclop
	switch mode {
	case git.AccessModeNo:
		return
	case git.AccessModeIfUndefined:
		if _, ok := os.LookupEnv(EnvGitReference); ok {
			return
		}
	}

	waiters.Add(1)

	go func() {
		defer waiters.Done()

		repo, err := git.Get()
		if err != nil {
			log.Printf("cannot find out correct git repo: %v", err)

			return
		}

		head, err := repo.Head()
		if err != nil {
			log.Printf("cannot lookup HEAD: %v", err)

			return
		}

		refName := head.Name()
		commitHash := head.Hash().String()
		commitHashShort := string([]rune(commitHash)[:GitCommitHashShortLength])

		sendValue(ctx, results, EnvGitReference, refName.String())
		sendValue(ctx, results, EnvGitReferenceShort, refName.Short())
		sendValue(ctx, results, EnvGitCommitHash, commitHash)
		sendValue(ctx, results, EnvGitCommitHashShort, commitHashShort)

		switch {
		case refName.IsBranch():
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeBranch.String())
		case refName.IsTag():
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeTag.String())
		case refName.IsRemote():
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeRemote.String())
		case refName.IsNote():
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeNote.String())
		default:
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeCommit.String())
		}

		if isDirty, err := repo.IsDirty(); err != nil {
			log.Printf("cannot detect if repository is dirty: %v", err)
		} else {
			sendValue(ctx, results, EnvGitIsDirty, strconv.FormatBool(isDirty))
		}
	}()
}
