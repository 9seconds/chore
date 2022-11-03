package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateMachineIdTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateMachineIdTestSuite) TestNoEnv() {
	env.GenerateMachineId(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Len(data, 1)
	suite.NotEmpty(data[env.EnvMachineId])
}

func (suite *GenerateMachineIdTestSuite) TestWithEnv() {
	suite.Setenv(env.EnvMachineId, "xxx")
	env.GenerateMachineId(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Empty(data)
}

func TestGenerateMachineId(t *testing.T) {
	t.Parallel()
	suite.Run(t, &GenerateMachineIdTestSuite{})
}
