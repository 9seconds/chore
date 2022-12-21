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
	env.GenerateIds(suite.Context(), suite.values, suite.wg, "xx", suite.args)

	data := suite.Collect()

	suite.Len(data, 4)
	suite.NotEmpty(data[env.EnvIDUnique])
	suite.NotEmpty(data[env.EnvIDChainUnique])
	suite.Equal(
		"v3CalS7yGSohMYSBoKlD5xFQCFQ8Bi02nVWOD2DAFZE",
		data[env.EnvIDIsolated])
	suite.Equal(
		"1Ek0XW68Lh0bis9SCJqkEjGSWe0STHuldW7l-DABoe4",
		data[env.EnvIDChainIsolated])
}

func (suite *GenerateIdsTestSuite) TestUniquesAreDifferent() {
	env.GenerateIds(suite.Context(), suite.values, suite.wg, "xx", suite.args)

	data := suite.Collect()
	unique := data[env.EnvIDUnique]
	suite.values = make(chan string, 1)

	env.GenerateIds(suite.Context(), suite.values, suite.wg, "xx", suite.args)

	data = suite.Collect()

	suite.NotEqual(unique, data[env.EnvIDUnique])
}

func (suite *GenerateIdsTestSuite) TestChainUnique() {
	suite.Setenv(env.EnvIDChainUnique, "xx2")
	env.GenerateIds(suite.Context(), suite.values, suite.wg, "xx", suite.args)

	data := suite.Collect()

	suite.NotContains(env.EnvIDChainUnique, data)
}

func TestGenerateIds(t *testing.T) {
	suite.Run(t, &GenerateIdsTestSuite{})
}
