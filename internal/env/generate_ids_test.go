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
		"MIXSka8V-ASGnPVXEQIRCTeyUWVD8fv9C35EiLMG1vg",
		data[env.EnvIDIsolated])
	suite.Equal(
		"ygwLMN-2WP1nlT8qS5-IIBI6SR0l0hF0jZA5xthqKfY",
		data[env.EnvIDChainIsolated])
}

func TestGenerateIds(t *testing.T) {
	suite.Run(t, &GenerateIdsTestSuite{})
}
