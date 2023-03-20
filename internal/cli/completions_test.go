package cli

import (
	"testing"

	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type CompleteRunTestSuite struct {
	suite.Suite
	testlib.CustomRootTestSuite

	cmd *cobra.Command
}

func (suite *CompleteRunTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.cmd = NewRoot("version")
	suite.cmd.AddCommand(&cobra.Command{})
}

func (suite *CompleteRunTestSuite) TestNamespaceError() {
	values, directive := completeRun(suite.cmd, nil, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteRunTestSuite) TestNamespaces() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := completeRun(suite.cmd, nil, "")

	suite.Equal([]string{"_", "xx"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteRunTestSuite) TestEmptyNamespace() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := completeRun(suite.cmd, []string{"xx"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteRunTestSuite) TestUnknownNamespace() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := completeRun(suite.cmd, []string{"a"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteRunTestSuite) TestListNamespace() {
	suite.EnsureScript("xx", "yy", "")
	suite.EnsureScript("xx", "y2", "")

	values, directive := completeRun(suite.cmd, []string{"xx"}, "")

	suite.Equal([]string{"y2", "yy"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteRunTestSuite) TestCompleteUnknownScript() {
	suite.EnsureScript("xx", "yy", "")
	suite.EnsureScript("xx", "y2", "")

	values, directive := completeRun(suite.cmd, []string{"xx", "xx"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteRunTestSuite) TestCompleteScriptWithoutConfig() {
	suite.EnsureScript("xx", "y", "")

	values, directive := completeRun(suite.cmd, []string{"xx", "y"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveDefault, directive)
}

func (suite *CompleteRunTestSuite) TestCompleteScriptWithConfig() {
	suite.EnsureScript("xx", "y", "")
	suite.EnsureScriptConfig("xx", "y", `
[parameters.param]
type = "string"

[flags.flag1]
description = "flag1 description"

[flags.param]
description = "works too"
	`)

	values, directive := completeRun(suite.cmd, []string{"xx", "y"}, "")

	suite.Equal([]string{
		"+flag1\tflag1 description (set)",
		"+param\tworks too (set)",
		"_flag1\tflag1 description (clear)",
		"_param\tworks too (clear)",
		"param=",
	}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteRunTestSuite) TestCompleteScriptWithConfigPositionalTime() {
	suite.EnsureScript("xx", "y", "")
	suite.EnsureScriptConfig("xx", "y", `
[parameters.param]
type = "string"

[flags.flag1]
description = "flag1 description"

[flags.param]
description = "works too"
	`)

	values, directive := completeRun(suite.cmd, []string{"xx", "y", "zz"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveDefault, directive)
}

func (suite *CompleteRunTestSuite) TestCompleteScriptWithConfigParam() {
	suite.EnsureScript("xx", "y", "")
	suite.EnsureScriptConfig("xx", "y", `
[parameters.param]
type = "string"

[flags.flag1]
description = "flag1 description"

[flags.param]
description = "works too"
	`)

	values, directive := completeRun(suite.cmd, []string{"xx", "y"}, "para")

	suite.Equal([]string{"param="}, values)
	suite.Equal(cobra.ShellCompDirectiveNoSpace, directive)
}

func (suite *CompleteRunTestSuite) TestCompleteScriptWithConfigFlag() {
	suite.EnsureScript("xx", "y", "")
	suite.EnsureScriptConfig("xx", "y", `
[parameters.param]
type = "string"

[flags.flag1]
description = "flag1 description"

[flags.param]
description = "works too"
	`)

	values, directive := completeRun(suite.cmd, []string{"xx", "y"}, "+f")

	suite.Equal([]string{"+flag1\tflag1 description (set)"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestCompleteRun(t *testing.T) {
	suite.Run(t, &CompleteRunTestSuite{})
}
