package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateRecursionTestSuite struct {
	EnvBaseTestSuite

	args argparse.ParsedArgs
}

func (suite *GenerateRecursionTestSuite) SetupTest() {
	suite.EnvBaseTestSuite.SetupTest()

	suite.args = argparse.ParsedArgs{
		Parameters: map[string][]string{
			"param1": {"33"},
			"param2": {"34 35"},
		},
		Flags: map[string]string{
			"flag1": argparse.FlagTrue,
			"flag2": argparse.FlagFalse,
		},
		Positional: []string{"pos1", "pos2", "pos3"},
	}
}

func (suite *GenerateRecursionTestSuite) TestEnv() {
	env.GenerateSelf(
		suite.Context(),
		suite.values,
		suite.wg,
		"namespace2",
		"script1",
		suite.args)

	data := suite.Collect()

	suite.Len(data, 1)
	suite.Contains(data[env.EnvSelf], "run namespace2 script1")
	suite.Contains(data[env.EnvSelf], "param1=33")
	suite.Contains(data[env.EnvSelf], "'param2=34 35'")
	suite.Contains(data[env.EnvSelf], "+flag1")
	suite.Contains(data[env.EnvSelf], "-flag2")
	suite.NotContains(data[env.EnvSelf], "pos1")
	suite.NotContains(data[env.EnvSelf], "pos2")
	suite.NotContains(data[env.EnvSelf], "pos3")
	suite.NotContains(data[env.EnvSelf], "--")
}

func TestGenerateRecursion(t *testing.T) {
	suite.Run(t, &GenerateRecursionTestSuite{})
}
