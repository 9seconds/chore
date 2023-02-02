package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/git"
	"github.com/stretchr/testify/suite"
)

type GenerateGitTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateGitTestSuite) TestGitAccessNo() {
	env.GenerateGit(suite.Context(), suite.values, suite.wg, git.AccessModeNo)
	suite.Empty(suite.Collect())
}

func (suite *GenerateGitTestSuite) TestGitAccessIfPresent() {
	suite.Setenv(env.EnvGitReference, "xx")
	env.GenerateGit(suite.Context(), suite.values, suite.wg, git.AccessModeIfUndefined)
	suite.Empty(suite.Collect())
}

func (suite *GenerateGitTestSuite) TestGitAccess() {
	env.GenerateGit(suite.Context(), suite.values, suite.wg, git.AccessModeAlways)
	data := suite.Collect()

	suite.Len(data, 6)
	suite.Contains(data, env.EnvGitReference)
	suite.Contains(data, env.EnvGitReferenceShort)
	suite.Contains(data, env.EnvGitReferenceType)
	suite.Contains(data, env.EnvGitCommitHash)
	suite.Contains(data, env.EnvGitCommitHashShort)
	suite.Contains(data, env.EnvGitIsDirty)
}

func TestGenerateGit(t *testing.T) {
	suite.Run(t, &GenerateGitTestSuite{})
}
