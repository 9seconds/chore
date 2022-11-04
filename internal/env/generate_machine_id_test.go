package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateMachineIDTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateMachineIDTestSuite) TestNoEnv() {
	env.GenerateMachineID(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Len(data, 1)
	suite.NotEmpty(data[env.EnvMachineID])
}

func (suite *GenerateMachineIDTestSuite) TestWithEnv() {
	suite.Setenv(env.EnvMachineID, "xxx")
	env.GenerateMachineID(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Empty(data)
}

func TestGenerateMachineID(t *testing.T) {
	suite.Run(t, &GenerateMachineIDTestSuite{})
}
