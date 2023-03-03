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
		Parameters: map[string]string{
			"k": "1",
		},
		Positional: []string{"1", "2"},
	}
}

func (suite *GenerateIdsTestSuite) TestNoEnvs() {
	env.GenerateIds(suite.Context(), suite.values, suite.wg, "xx", "xz", suite.args)

	data := suite.Collect()

	suite.Len(data, 4)
	suite.NotEmpty(data[env.EnvIDRun])
	suite.NotEmpty(data[env.EnvIDChainRun])
	suite.Equal(
		"0NuX7xYp98TzxeCRL3WWmEDd7_m8m-8ZX9mOjvIUeXI",
		data[env.EnvIDIsolated])
	suite.Equal(
		"Em5rkM3HRPOqufbGU5YOjJ3ZFKNfPtzYBlMWVV6vwkY",
		data[env.EnvIDChainIsolated])
}

func TestGenerateIds(t *testing.T) {
	suite.Run(t, &GenerateIdsTestSuite{})
}
