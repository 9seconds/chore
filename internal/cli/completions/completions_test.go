package completions_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type CompletionsTestSuite struct {
	suite.Suite
	testlib.CustomRootTestSuite

	fn func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)
}

func (suite *CompletionsTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())
}

func (suite *CompletionsTestSuite) Run(args ...string) ([]string, cobra.ShellCompDirective) {
	return suite.fn(nil, args, "")
}

type CompleteNamespacesTestSuite struct {
	CompletionsTestSuite
}

func (suite *CompleteNamespacesTestSuite) SetupTest() {
	suite.CompletionsTestSuite.SetupTest()

	suite.fn = completions.CompleteNamespaces
}

func (suite *CompleteNamespacesTestSuite) TestError() {
	namespaces, directive := suite.Run()

	suite.Empty(namespaces)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteNamespacesTestSuite) TestEmpty() {
	suite.EnsureDir(paths.ConfigRoot())

	namespaces, directive := suite.Run()

	suite.Empty(namespaces)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespacesTestSuite) TestSomething() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("yy"))
	suite.EnsureDir(paths.ConfigNamespace("__"))

	namespaces, directive := suite.Run()

	suite.Equal([]string{"__", "xx", "yy"}, namespaces)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

type CompleteNamespaceScriptTestSuite struct {
	CompletionsTestSuite
}

func (suite *CompleteNamespaceScriptTestSuite) SetupTest() {
	suite.CompletionsTestSuite.SetupTest()

	suite.fn = completions.CompleteNamespaceScript
}

func (suite *CompleteNamespaceScriptTestSuite) TestNamespaceError() {
	values, directive := suite.Run()

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestNamespaces() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := suite.Run()

	suite.Equal([]string{"_", "xx"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestEmptyNamespace() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := suite.Run("xx")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestUnknownNamespace() {
	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureDir(paths.ConfigNamespace("_"))

	values, directive := suite.Run("a")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveError, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestListNamespace() {
	suite.EnsureScript("xx", "yy", "")
	suite.EnsureScript("xx", "y2", "")

	values, directive := suite.Run("xx")

	suite.Equal([]string{"y2", "yy"}, values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func (suite *CompleteNamespaceScriptTestSuite) TestManyArgs() {
	suite.EnsureScript("xx", "yy", "")
	suite.EnsureScript("xx", "y2", "")

	values, directive := suite.Run("xx", "xx")

	suite.Empty(values)
	suite.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestCompleteNamespaces(t *testing.T) {
	suite.Run(t, &CompleteNamespacesTestSuite{})
}

func TestCompleteNamespaceScript(t *testing.T) {
	suite.Run(t, &CompleteNamespaceScriptTestSuite{})
}
