package git_test

import (
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/git"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	testlib.GitTestSuite

	repo *git.Repo
}

func (suite *RepoTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	rootPath := suite.RootPath()

	suite.GitTestSuite.Setup(suite.T(), rootPath)

	suite.EnsureFile(filepath.Join(rootPath, "aa"), "a", 0o600)
	suite.GitAdd("aa")
	suite.GitCommit("finally")

	suite.GitCreateRemote("remote1", "")

	headHash := suite.GitHead().Hash()

	suite.GitCreateBranch("br1", headHash)
	suite.GitCreateBranch("br2", headHash)
	suite.GitCreateRemoteBranch("remote-branch", "remote1", headHash)
	suite.GitCreateTag("light", headHash)
	suite.GitCreateAnnotatedTag("annotated", "Hello", headHash)

	gitRepo, err := git.New()
	require.NoError(suite.T(), err)

	suite.repo = gitRepo
}

func (suite *RepoTestSuite) TestBranches() {
	testTable := map[string]bool{
		"master":        true,
		"br1":           true,
		"br2":           true,
		"remote-branch": true,
		"br3":           false,
		"":              false,
	}

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			present, err := suite.repo.HasBranch(testValue)

			assert.NoError(t, err)

			if isValid {
				assert.True(t, present)
			} else {
				assert.False(t, present)
			}
		})
	}
}

func (suite *RepoTestSuite) TestRemotes() {
	testTable := map[string]bool{
		"remote1/remote-branch": true,
		"remote1/master":        false,
		"origin/master":         false,
		"":                      false,
	}

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			present, err := suite.repo.HasRemote(testValue)

			assert.NoError(t, err)

			if isValid {
				assert.True(t, present)
			} else {
				assert.False(t, present)
			}
		})
	}
}

func (suite *RepoTestSuite) TestTags() {
	testTable := map[string]bool{
		"":          false,
		"v1":        false,
		"light":     true,
		"annotated": true,
	}

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			present, err := suite.repo.HasTag(testValue)

			assert.NoError(t, err)

			if isValid {
				assert.True(t, present)
			} else {
				assert.False(t, present)
			}
		})
	}
}

func (suite *RepoTestSuite) TestRevision() {
	testTable := map[string]bool{
		"    ":                  false,
		"light":                 true,
		"annotated":             true,
		"remote-branch":         true,
		"remote1/remote-branch": true,
		"origin/master":         false,
		"br1":                   true,
		"master":                true,
	}

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			present, err := suite.repo.HasRevision(testValue)

			assert.NoError(t, err)

			if isValid {
				assert.True(t, present)
			} else {
				assert.False(t, present)
			}
		})
	}

	head, err := suite.repo.Head()
	suite.NoError(err)

	hash := head.Hash().String()

	suite.T().Run("full-hash", func(t *testing.T) {
		present, err := suite.repo.HasRevision(hash)

		assert.NoError(t, err)
		assert.True(t, present)
	})

	suite.T().Run("short-hash", func(t *testing.T) {
		present, err := suite.repo.HasRevision(hash[:8])

		assert.NoError(t, err)
		assert.True(t, present)
	})
}

func (suite *RepoTestSuite) TestCleanDir() {
	ok, err := suite.repo.IsDirty()

	suite.NoError(err)
	suite.False(ok)
}

func (suite *RepoTestSuite) TestModifiedFile() {
	suite.EnsureFile(filepath.Join(suite.RootPath(), "aa"), "b", 0o600)

	ok, err := suite.repo.IsDirty()

	suite.NoError(err)
	suite.True(ok)
}

func (suite *RepoTestSuite) TestUntracked() {
	suite.EnsureFile(filepath.Join(suite.RootPath(), "bb"), "b", 0o600)

	ok, err := suite.repo.IsDirty()

	suite.NoError(err)
	suite.True(ok)
}

func (suite *RepoTestSuite) TestCommit() {
	head := suite.GitHead().Hash().String()

	ok, err := suite.repo.HasCommit(head)
	suite.True(ok)
	suite.NoError(err)

	ok, err = suite.repo.HasCommit(head[:8])
	suite.True(ok)
	suite.NoError(err)

	ok, err = suite.repo.HasCommit("xxx")
	suite.False(ok)
	suite.NoError(err)

	ok, err = suite.repo.HasCommit("br1")
	suite.False(ok)
	suite.NoError(err)
}

func TestRepo(t *testing.T) {
	suite.Run(t, &RepoTestSuite{})
}
