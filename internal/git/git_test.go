package git_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetRepoTestSuite struct {
	suite.Suite
}

func (suite *GetRepoTestSuite) SetupTest() {
	wd, _ := os.Getwd()
	t := suite.T()

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(wd))
	})
}

func (suite *GetRepoTestSuite) TestInNonGit() {
	suite.NoError(os.Chdir(suite.T().TempDir()))

	_, err := git.GetRepo()
	suite.ErrorContains(err, "cannot open git repository")
}

func (suite *GetRepoTestSuite) TestInvalidWorkingDir() {
	dir := suite.T().TempDir()

	suite.NoError(os.Chdir(dir))
	suite.NoError(os.RemoveAll(dir))

	_, err := git.GetRepo()
	suite.ErrorContains(err, "cannot find out current working dir")
}

func (suite *GetRepoTestSuite) TestOk() {
	repo, err := git.GetRepo()
	suite.NoError(err)

	_, err = repo.Head()
	suite.NoError(err)
}

func TestGetGitRepo(t *testing.T) {
	suite.Run(t, &GetRepoTestSuite{})
}

func TestGetAccessMode(t *testing.T) {
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

			mode, err := git.GetAccessMode(testName)

			if isValid {
				assert.NoError(t, err)

				if testName == "" {
					assert.Equal(t, git.AccessIfUndefined, mode)
				} else {
					assert.Equal(t, testName, mode.String())
				}
			} else {
				assert.ErrorIs(t, err, git.ErrAccessModeUnknown)
			}
		})
	}
}
