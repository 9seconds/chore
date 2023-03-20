package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateSelfTestSuite struct {
	BaseTestSuite

	args argparse.ParsedArgs
}

func (suite *GenerateSelfTestSuite) SetupTest() {
	suite.BaseTestSuite.SetupTest()

	suite.args = argparse.ParsedArgs{
		Parameters: map[string][]string{
			"param1": {"33"},
			"param2": {"34 35"},
		},
		Flags: map[string]bool{
			"flag1": true,
			"flag2": false,
		},
		Positional: []string{"pos1", "pos2", "pos3"},
	}
}

func (suite *GenerateSelfTestSuite) Test() {
	env.GenerateSelf(
		suite.Context(),
		suite.values,
		suite.wg,
		"namespace2",
		"script1",
		suite.args)

	data := suite.Collect()

	suite.Len(data, 2)
	suite.Contains(data, env.Bin)
	suite.Contains(data[env.Self], "run namespace2 script1")
	suite.Contains(data[env.Self], "param1=33")
	suite.Contains(data[env.Self], "'param2=34 35'")
	suite.Contains(data[env.Self], "+flag1")
	suite.NotContains(data[env.Self], "flag2")
	suite.NotContains(data[env.Self], "pos1")
	suite.NotContains(data[env.Self], "pos2")
	suite.NotContains(data[env.Self], "pos3")
	suite.NotContains(data[env.Self], "--")
}

func TestGenerateSelf(t *testing.T) {
	suite.Run(t, &GenerateSelfTestSuite{})
}
