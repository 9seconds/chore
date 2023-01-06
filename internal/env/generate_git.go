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
	case git.AccessNo:
		return
	case git.AccessIfUndefined:
		if _, ok := os.LookupEnv(EnvGitReference); ok {
			return
		}
	}

	waiters.Add(1)

	go func() {
		defer waiters.Done()

		repo, err := git.GetRepo()
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
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeBranch)
		case refName.IsTag():
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeTag)
		case refName.IsRemote():
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeRemote)
		case refName.IsNote():
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeNote)
		default:
			sendValue(ctx, results, EnvGitReferenceType, git.RefTypeCommit)
		}

		workTree, err := repo.Worktree()
		if err != nil {
			log.Printf("cannot get work tree: %v", err)

			return
		}

		status, err := workTree.Status()
		if err != nil {
			log.Printf("cannot get work tree status: %v", err)

			return
		}

		sendValue(ctx, results, EnvGitIsDirty, strconv.FormatBool(!status.IsClean()))
	}()
}
