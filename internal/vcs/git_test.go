package vcs_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/vcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetGitRepoTestSuite struct {
	suite.Suite
}

func (suite *GetGitRepoTestSuite) SetupTest() {
	wd, _ := os.Getwd()
	t := suite.T()

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(wd))
	})
}

func (suite *GetGitRepoTestSuite) TestInNonGit() {
	suite.NoError(os.Chdir(suite.T().TempDir()))

	_, err := vcs.GetGitRepo()
	suite.ErrorContains(err, "cannot open git repository")
}

func (suite *GetGitRepoTestSuite) TestInvalidWorkingDir() {
	dir := suite.T().TempDir()

	suite.NoError(os.Chdir(dir))
	suite.NoError(os.RemoveAll(dir))

	_, err := vcs.GetGitRepo()
	suite.ErrorContains(err, "cannot find out current working dir")
}

func (suite *GetGitRepoTestSuite) TestOk() {
	repo, err := vcs.GetGitRepo()
	suite.NoError(err)

	_, err = repo.Head()
	suite.NoError(err)
}

func TestGetGitRepo(t *testing.T) {
	suite.Run(t, &GetGitRepoTestSuite{})
}

func TestGetGitAccessMode(t *testing.T) {
	testTable := map[string]bool{
		"if_undef":     false,
		"":             true,
		"xj":           false,
		"if_undefined": true,
		"yes":          false,
		"1":            false,
		"no":           true,
		"always":       true,
		"never":        false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			mode, err := vcs.GetGitAccessMode(testName)

			if isValid {
				assert.NoError(t, err)

				if testName == "" {
					assert.Equal(t, vcs.GitAccessIfUndefined, mode)
				} else {
					assert.Equal(t, testName, mode.String())
				}
			} else {
				assert.ErrorIs(t, err, vcs.ErrGitAccessModeUnknown)
			}
		})
	}
}
