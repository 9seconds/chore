package config_test

import (
	"errors"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/git"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParameterGitTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
	testlib.CustomRootTestSuite
	testlib.GitTestSuite

	repo *git.Repo
}

func (suite *ParameterGitTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())
	suite.CtxTestSuite.Setup(suite.T())

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

func (suite *ParameterGitTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewGit("", testValue, nil, git.New)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterGitTestSuite) TestType() {
	param, err := config.NewGit("", false, nil, git.New)
	suite.NoError(err)
	suite.Equal(config.ParameterGit, param.Type())
}

func (suite *ParameterGitTestSuite) TestFailIfCantInit() {
	referenceError := errors.New("")

	_, err := config.NewGit("", false, nil, func() (*git.Repo, error) {
		return nil, referenceError
	})
	suite.ErrorIs(err, referenceError)
}

func (suite *ParameterGitTestSuite) TestTags() {
	testTable := map[string]bool{
		"light":                 true,
		"annotated":             true,
		"":                      false,
		"xx":                    false,
		"br1":                   false,
		"remote1":               false,
		"remote1/remote-branch": false,
	}

	param, err := config.NewGit("", false, map[string]string{
		"type": "tag",
	}, git.New)
	suite.NoError(err)

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			err := param.Validate(suite.Context(), testValue)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, config.ErrGitIncorrectRefType)
			}
		})
	}

	suite.T().Run("hash", func(t *testing.T) {
		head := suite.GitHead().Hash()
		assert.ErrorIs(
			t,
			param.Validate(suite.Context(), head.String()),
			config.ErrGitIncorrectRefType)
	})
}

func (suite *ParameterGitTestSuite) TestBranch() {
	testTable := map[string]bool{
		"light":                 false,
		"annotated":             false,
		"":                      false,
		"xx":                    false,
		"br1":                   true,
		"remote1":               false,
		"remote1/remote-branch": false,
	}

	param, err := config.NewGit("", false, map[string]string{
		"type": "branch",
	}, git.New)
	suite.NoError(err)

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			err := param.Validate(suite.Context(), testValue)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, config.ErrGitIncorrectRefType)
			}
		})
	}

	suite.T().Run("hash", func(t *testing.T) {
		head := suite.GitHead().Hash()
		assert.ErrorIs(
			t,
			param.Validate(suite.Context(), head.String()),
			config.ErrGitIncorrectRefType)
	})
}

func (suite *ParameterGitTestSuite) TestRemote() {
	testTable := map[string]bool{
		"light":                 false,
		"annotated":             false,
		"":                      false,
		"xx":                    false,
		"br1":                   false,
		"remote1":               false,
		"remote1/remote-branch": true,
	}

	param, err := config.NewGit("", false, map[string]string{
		"type": "remote",
	}, git.New)
	suite.NoError(err)

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			err := param.Validate(suite.Context(), testValue)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, config.ErrGitIncorrectRefType)
			}
		})
	}

	suite.T().Run("hash", func(t *testing.T) {
		head := suite.GitHead().Hash()
		assert.ErrorIs(
			t,
			param.Validate(suite.Context(), head.String()),
			config.ErrGitIncorrectRefType)
	})
}

func (suite *ParameterGitTestSuite) TestCommit() {
	testTable := []string{
		"light",
		"annotated",
		"",
		"xx",
		"br1",
		"remote1",
		"remote1/remote-branch",
	}

	param, err := config.NewGit("", false, map[string]string{
		"type": "commit",
	}, git.New)
	suite.NoError(err)

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.Error(t, param.Validate(suite.Context(), testValue))
		})
	}

	suite.T().Run("hash", func(t *testing.T) {
		head := suite.GitHead().Hash()
		assert.NoError(t, param.Validate(suite.Context(), head.String()))
	})
}

func (suite *ParameterGitTestSuite) TestAll() {
	testTable := map[string]bool{
		"light":                 true,
		"annotated":             true,
		"":                      false,
		"xx":                    false,
		"br1":                   true,
		"remote1":               false,
		"remote1/remote-branch": true,
	}

	param, err := config.NewGit("", false, nil, git.New)
	suite.NoError(err)

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		suite.T().Run(testValue, func(t *testing.T) {
			err := param.Validate(suite.Context(), testValue)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}

	suite.T().Run("hash", func(t *testing.T) {
		head := suite.GitHead().Hash()
		assert.NoError(t, param.Validate(suite.Context(), head.String()))
	})
}

func TestParameterGit(t *testing.T) {
	suite.Run(t, &ParameterGitTestSuite{})
}
