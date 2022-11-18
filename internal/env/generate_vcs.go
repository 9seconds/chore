package env

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/9seconds/chore/internal/vcs"
)

const GitCommitHashShortLength = 12

func GenerateGit(ctx context.Context, results chan<- string, waiters *sync.WaitGroup, mode vcs.GitAccessMode) { //nolint: cyclop
	switch mode {
	case vcs.GitAccessNo:
		return
	case vcs.GitAccessIfUndefined:
		if _, ok := os.LookupEnv(EnvGitReference); ok {
			return
		}
	}

	waiters.Add(1)

	go func() {
		defer waiters.Done()

		repo, err := vcs.GetGitRepo()
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
			sendValue(ctx, results, EnvGitReferenceType, vcs.GitRefTypeBranch)
		case refName.IsTag():
			sendValue(ctx, results, EnvGitReferenceType, vcs.GitRefTypeTag)
		case refName.IsRemote():
			sendValue(ctx, results, EnvGitReferenceType, vcs.GitRefTypeRemote)
		case refName.IsNote():
			sendValue(ctx, results, EnvGitReferenceType, vcs.GitRefTypeNote)
		default:
			sendValue(ctx, results, EnvGitReferenceType, vcs.GitRefTypeCommit)
		}
	}()
}
