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
		Keywords: map[string]string{
			"k": "1",
		},
		Positional: []string{"1", "2"},
	}
}

func (suite *GenerateIdsTestSuite) TestNoEnvs() {
	env.GenerateIds(suite.ctx, suite.values, suite.wg, "xx", suite.args)

	data := suite.Collect()

	suite.Len(data, 4)
	suite.NotEmpty(data[env.EnvIdUnique])
	suite.NotEmpty(data[env.EnvIdChainUnique])
	suite.Equal(
		"1lcK3kcAuwj9dF91XbXnmdZybO8Rj51HYMywNbWfWNI",
		data[env.EnvIdIsolated])
	suite.Equal(
		"wZkkTeL4GOiyL0PkOKEqGenGbuw0xdIvVvXRnN4qEi0",
		data[env.EnvIdChainIsolated])
}

func (suite *GenerateIdsTestSuite) TestUniquesAreDifferent() {
	env.GenerateIds(suite.ctx, suite.values, suite.wg, "xx", suite.args)

	data := suite.Collect()
	unique := data[env.EnvIdUnique]
	suite.values = make(chan string, 1)

	env.GenerateIds(suite.ctx, suite.values, suite.wg, "xx", suite.args)

	data = suite.Collect()

	suite.NotEqual(unique, data[env.EnvIdUnique])
}

func (suite *GenerateIdsTestSuite) TestChainUnique() {
	suite.Setenv(env.EnvIdChainUnique, "xx2")
	env.GenerateIds(suite.ctx, suite.values, suite.wg, "xx", suite.args)

	data := suite.Collect()

	suite.NotContains(env.EnvIdChainUnique, data)
}

func TestGenerateIds(t *testing.T) {
	t.Parallel()
	suite.Run(t, &GenerateIdsTestSuite{})
}
