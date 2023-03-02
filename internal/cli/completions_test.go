package cli

import (
	"testing"

	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type CompletionsTestSuite struct {
	suite.Suite
	testlib.CustomRootTestSuite

	cmd *cobra.Command
}

func (suite *CompletionsTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	root := NewRoot("version")
	root.AddCommand(&cobra.Command{})

	suite.cmd = root
}

type CompleteNamespacesTestSuite struct {
	CompletionsTestSuite
}

func (suite *CompleteNamespacesTestSuite) TestError() {
	namespaces, directive := completeNamespaces()

	suite.Empty(namespaces)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteNamespacesTestSuite) TestEmpty() {
	suite.EnsureDir(paths.ConfigRoot())

	namespaces, directive := completeNamespaces()

	suite.Empty(namespaces)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespacesTestSuite) TestSomething() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("yy"))
	suite.EnsureDir(paths.ConfigNamespace("__"))

	namespaces, directive := completeNamespaces()

	suite.Equal([]string{"__", "xx", "yy"}, namespaces)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

type CompleteScriptsTestSuite struct {
	CompletionsTestSuite
}

func (suite *CompleteScriptsTestSuite) TestError() {
	scripts, directive := completeScripts("xx")

	suite.Empty(scripts)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteScriptsTestSuite) TestEmpty() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))

	scripts, directive := completeScripts("xx")

	suite.Empty(scripts)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteScriptsTestSuite) TestList() {
	suite.EnsureScript("xx", "yy", "")
	suite.EnsureScript("xx", "y2", "")
	suite.EnsureScript("xx", "_y2", "")
	suite.EnsureScript("xz", "y2", "")
	suite.EnsureDir(paths.ConfigNamespace("xx"))

	scripts, directive := completeScripts("xx")

	suite.Equal([]string{"_y2", "y2", "yy"}, scripts)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

type CompleteNamespaceScriptTestSuite struct {
	CompletionsTestSuite
}

func (suite *CompleteNamespaceScriptTestSuite) TestNamespaceError() {
	values, directive := completeNamespaceScript(suite.cmd, nil, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestNamespaces() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := completeNamespaceScript(suite.cmd, nil, "")

	suite.Equal([]string{"_", "xx"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestEmptyNamespace() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := completeNamespaceScript(suite.cmd, []string{"xx"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestUnknownNamespace() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := completeNamespaceScript(suite.cmd, []string{"a"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestListNamespace() {
	suite.EnsureScript("xx", "yy", "")
	suite.EnsureScript("xx", "y2", "")

	values, directive := completeNamespaceScript(suite.cmd, []string{"xx"}, "")

	suite.Equal([]string{"y2", "yy"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestManyArgs() {
	suite.EnsureScript("xx", "yy", "")
	suite.EnsureScript("xx", "y2", "")

	values, directive := completeNamespaceScript(suite.cmd, []string{"xx", "xx"}, "")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

type CompleteRunTestSuite struct {
	CompletionsTestSuite
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
		"+flag1\tflag1 description (yes)",
		"+param\tworks too (yes)",
		"-flag1\tflag1 description (no)",
		"-param\tworks too (no)",
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

	suite.Equal([]string{"+flag1\tflag1 description (yes)"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestCompleteNamespaces(t *testing.T) {
	suite.Run(t, &CompleteNamespacesTestSuite{})
}

func TestCompleteScripts(t *testing.T) {
	suite.Run(t, &CompleteScriptsTestSuite{})
}

func TestCompleteNamespaceScript(t *testing.T) {
	suite.Run(t, &CompleteNamespaceScriptTestSuite{})
}

func TestCompleteRun(t *testing.T) {
	suite.Run(t, &CompleteRunTestSuite{})
}
