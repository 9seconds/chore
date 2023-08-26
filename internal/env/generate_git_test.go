package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/git"
	"github.com/stretchr/testify/suite"
)

type GenerateGitTestSuite struct {
	BaseTestSuite
}

func (suite *GenerateGitTestSuite) TestGitAccessNo() {
	env.GenerateGit(suite.Context(), suite.values, suite.wg, git.AccessModeNo)
	suite.Empty(suite.Collect())
}

func (suite *GenerateGitTestSuite) TestGitAccessIfPresent() {
	suite.Setenv(env.GitReference, "xx")
	env.GenerateGit(suite.Context(), suite.values, suite.wg, git.AccessModeIfUndefined)
	suite.Empty(suite.Collect())
}

func (suite *GenerateGitTestSuite) TestGitAccess() {
	env.GenerateGit(suite.Context(), suite.values, suite.wg, git.AccessModeAlways)
	data := suite.Collect()

	suite.Len(data, 6)
	suite.Contains(data, env.GitReference)
	suite.Contains(data, env.GitReferenceShort)
	suite.Contains(data, env.GitReferenceType)
	suite.Contains(data, env.GitCommitHash)
	suite.Contains(data, env.GitCommitHashShort)
	suite.Contains(data, env.GitIsDirty)
}

func TestGenerateGit(t *testing.T) {
	suite.Run(t, &GenerateGitTestSuite{})
}
