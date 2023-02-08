package testlib

import (
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"
)

type GitTestSuite struct {
	t        *testing.T
	workTree *git.Worktree
	repo     *git.Repository
}

func (suite *GitTestSuite) Setup(t *testing.T, path string) {
	t.Helper()

	suite.t = t

	t.Setenv("GIT_DIR", path)

	repo, err := git.PlainInit(path, false)
	require.NoError(t, err)

	suite.repo = repo

	tree, err := repo.Worktree()
	require.NoError(t, err)

	suite.workTree = tree
}

func (suite *GitTestSuite) GitCommit(message string) {
	suite.t.Helper()

	_, err := suite.workTree.Commit(message, &git.CommitOptions{})
	require.NoError(suite.t, err)
}

func (suite *GitTestSuite) GitAdd(path string) {
	suite.t.Helper()

	require.NoError(suite.t, suite.workTree.AddGlob(path))
}

func (suite *GitTestSuite) GitHead() *plumbing.Reference {
	suite.t.Helper()

	head, err := suite.repo.Head()
	require.NoError(suite.t, err)

	return head
}

func (suite *GitTestSuite) GitCreateRemote(name, url string) {
	suite.t.Helper()

	if url == "" {
		url = fmt.Sprintf("https://github.com/9seconds/%s.git", name)
	}

	_, err := suite.repo.CreateRemote(&config.RemoteConfig{
		Name: "remote1",
		URLs: []string{url},
	})
	require.NoError(suite.t, err)
}

func (suite *GitTestSuite) GitCreateBranch(name string, hash plumbing.Hash) {
	suite.t.Helper()

	branchName := plumbing.NewBranchReferenceName(name)
	branchRef := plumbing.NewHashReference(branchName, hash)

	require.NoError(suite.t, suite.repo.Storer.SetReference(branchRef))
}

func (suite *GitTestSuite) GitCreateRemoteBranch(name, remote string, hash plumbing.Hash) {
	suite.t.Helper()

	branchName := plumbing.NewBranchReferenceName(name)
	remoteRef := plumbing.NewRemoteReferenceName(remote, name)

	require.NoError(
		suite.t,
		suite.repo.CreateBranch(&config.Branch{
			Name:   name,
			Remote: remote,
			Merge:  branchName,
		}))
	require.NoError(
		suite.t,
		suite.repo.Storer.SetReference(
			plumbing.NewSymbolicReference(branchName, remoteRef)))
	require.NoError(
		suite.t,
		suite.repo.Storer.SetReference(
			plumbing.NewHashReference(remoteRef, hash)))
}

func (suite *GitTestSuite) GitCreateTag(name string, hash plumbing.Hash) {
	_, err := suite.repo.CreateTag(name, hash, nil)
	require.NoError(suite.t, err)
}

func (suite *GitTestSuite) GitCreateAnnotatedTag(name, message string, hash plumbing.Hash) {
	_, err := suite.repo.CreateTag(name, hash, &git.CreateTagOptions{
		Message: message,
	})
	require.NoError(suite.t, err)
}
