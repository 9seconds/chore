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
		if _, ok := os.LookupEnv(GitReference); ok {
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

		sendValue(ctx, results, GitReference, refName.String())
		sendValue(ctx, results, GitReferenceShort, refName.Short())
		sendValue(ctx, results, GitCommitHash, commitHash)
		sendValue(ctx, results, GitCommitHashShort, commitHashShort)

		switch {
		case refName.IsBranch():
			sendValue(ctx, results, GitReferenceType, git.RefTypeBranch.String())
		case refName.IsTag():
			sendValue(ctx, results, GitReferenceType, git.RefTypeTag.String())
		case refName.IsRemote():
			sendValue(ctx, results, GitReferenceType, git.RefTypeRemote.String())
		case refName.IsNote():
			sendValue(ctx, results, GitReferenceType, git.RefTypeNote.String())
		default:
			sendValue(ctx, results, GitReferenceType, git.RefTypeCommit.String())
		}

		if isDirty, err := repo.IsDirty(); err != nil {
			log.Printf("cannot detect if repository is dirty: %v", err)
		} else {
			sendValue(ctx, results, GitIsDirty, strconv.FormatBool(isDirty))
		}
	}()
}
