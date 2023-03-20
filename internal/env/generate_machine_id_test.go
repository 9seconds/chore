package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateMachineIDTestSuite struct {
	BaseTestSuite
}

func (suite *GenerateMachineIDTestSuite) TestNo() {
	env.GenerateMachineID(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Len(data, 1)
	suite.NotEmpty(data[env.MachineID])
}

func (suite *GenerateMachineIDTestSuite) TestWith() {
	suite.Setenv(env.MachineID, "xxx")
	env.GenerateMachineID(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Empty(data)
}

func TestGenerateMachineID(t *testing.T) {
	suite.Run(t, &GenerateMachineIDTestSuite{})
}
