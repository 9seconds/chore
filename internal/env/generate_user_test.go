package env_test

import (
	"os/user"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateUserTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateUserTestSuite) TestNoEnv() {
	user, err := user.Current()
	if err != nil {
		suite.T().Skipf("Test skipped because of %v", err)
	}

	env.GenerateUser(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Len(data, 3)
	suite.Equal(user.Uid, data[env.EnvUserUID])
	suite.Equal(user.Gid, data[env.EnvUserGID])
	suite.Equal(user.Username, data[env.EnvUserName])
}

func (suite *GenerateUserTestSuite) TestEnv() {
	suite.Setenv(env.EnvUserName, "xx")

	env.GenerateUser(suite.Context(), suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func TestGenerateUser(t *testing.T) {
	suite.Run(t, &GenerateUserTestSuite{})
}
