package env_test

import (
	"os/user"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateUserTestSuite struct {
	BaseTestSuite
}

func (suite *GenerateUserTestSuite) TestNo() {
	user, err := user.Current()
	if err != nil {
		suite.T().Skipf("Test skipped because of %v", err)
	}

	env.GenerateUser(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Len(data, 3)
	suite.Equal(user.Uid, data[env.UserUID])
	suite.Equal(user.Gid, data[env.UserGID])
	suite.Equal(user.Username, data[env.UserName])
}

func (suite *GenerateUserTestSuite) Test() {
	suite.Setenv(env.UserName, "xx")

	env.GenerateUser(suite.Context(), suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func TestGenerateUser(t *testing.T) {
	suite.Run(t, &GenerateUserTestSuite{})
}
