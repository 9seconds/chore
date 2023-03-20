package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateIdsTestSuite struct {
	EnvBaseTestSuite

	args argparse.ParsedArgs
}

func (suite *GenerateIdsTestSuite) SetupTest() {
	suite.EnvBaseTestSuite.SetupTest()

	suite.args = argparse.ParsedArgs{
		Parameters: map[string][]string{
			"k": {"1"},
		},
		Positional: []string{"1", "2"},
	}
}

func (suite *GenerateIdsTestSuite) TestNoEnvs() {
	env.GenerateIds(suite.Context(), suite.values, suite.wg, "xx", suite.args)

	data := suite.Collect()

	suite.Len(data, 3)
	suite.NotEmpty(data[env.EnvIDChainRun])
	suite.Equal(
		"bPw4mf0i7ORf4zXimc4AJl0AjO5uiSFqWgmdhPTrJ-A",
		data[env.EnvIDIsolated])
	suite.Equal(
		"OANtkcb4mtiB_O-4ovEDuNE21yga8uQOvXHpH60aldM",
		data[env.EnvIDChainIsolated])
}

func TestGenerateIds(t *testing.T) {
	suite.Run(t, &GenerateIdsTestSuite{})
}
