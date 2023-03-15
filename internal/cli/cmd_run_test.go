package cli_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/stretchr/testify/suite"
)

type CmdRunTestSuite struct {
	CmdTestSuite
}

func (suite *CmdRunTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("run", cli.NewRun)

	suite.EnsureScriptConfig("ns", "s", `
description = "ZZY"
git = "always"
network = true

[flags.flag1]
description = "This is a description for flag1"
required = false  # default value

[parameters.param]
description = "Never knows best"
type = "string"
required = true

[paramters.param.spec]
ascii = true
regexp = '^\d\w+$'`)
	suite.EnsureScript("ns", "s", "echo $CHORE_P_PARAM")
}

func (suite *CmdRunTestSuite) TestOk() {
	_, err := suite.ExecuteCommand("ns", "s", "param=ppp")
	suite.NoError(err)
}

func TestCmdRun(t *testing.T) {
	suite.Run(t, &CmdRunTestSuite{})
}
