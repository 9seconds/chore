package script_test

import (
	"testing"

	"github.com/9seconds/chore/internal/script"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
}

func (suite *ListTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.EnsureDir(suite.ConfigNamespacePath("ns"))
	suite.EnsureScript("nb", "aa", "echo 1")
	suite.EnsureScript("nb", "ab", "echo 1")
	suite.EnsureFile(suite.ConfigScriptPath("nb", "a"), "1", 0o400)
	suite.EnsureDir(suite.ConfigScriptPath("nb", "b"))
}

type ListNamespacesTestSuite struct {
	ListTestSuite
}

func (suite *ListNamespacesTestSuite) TestListAll() {
	namespaces, err := script.ListNamespaces("")
	suite.NoError(err)
	suite.Equal([]string{"nb", "ns"}, namespaces)
}

func (suite *ListNamespacesTestSuite) TestPrefixed() {
	namespaces, err := script.ListNamespaces("ns")
	suite.NoError(err)
	suite.Equal([]string{"ns"}, namespaces)
}

func (suite *ListNamespacesTestSuite) TestEmpty() {
	namespaces, err := script.ListNamespaces("f")
	suite.NoError(err)
	suite.Empty(namespaces)
}

type ListScriptsTestSuite struct {
	ListTestSuite
}

func (suite *ListScriptsTestSuite) TestListAll() {
	scripts, err := script.ListScripts("nb", "")
	suite.NoError(err)
	suite.Equal([]string{"aa", "ab"}, scripts)
}

func (suite *ListScriptsTestSuite) TestPrefixed() {
	scripts, err := script.ListScripts("nb", "aa")
	suite.NoError(err)
	suite.Equal([]string{"aa"}, scripts)
}

func (suite *ListScriptsTestSuite) TestNothing() {
	scripts, err := script.ListScripts("nb", "b")
	suite.NoError(err)
	suite.Empty(scripts)
}

func (suite *ListScriptsTestSuite) TestEmpty() {
	scripts, err := script.ListScripts("ns", "")
	suite.NoError(err)
	suite.Empty(scripts)
}

type FindScriptTestSuite struct {
	ListTestSuite
}

func (suite *FindScriptTestSuite) TestExactFound() {
	scr, err := script.FindScript("nb", "aa")
	suite.NoError(err)
	suite.Equal("nb", scr.Namespace)
	suite.Equal("aa", scr.Executable)
}

func (suite *FindScriptTestSuite) TestNotExactFound() {
	scr, err := script.FindScript("n", "aa")
	suite.NoError(err)
	suite.Equal("nb", scr.Namespace)
	suite.Equal("aa", scr.Executable)
}

func (suite *FindScriptTestSuite) TestAmbigous() {
	_, err := script.FindScript("n", "a")
	suite.ErrorContains(err, "aa")
	suite.ErrorContains(err, "ab")
}

func (suite *FindScriptTestSuite) TestNotFound() {
	_, err := script.FindScript("n", "b")
	suite.ErrorContains(err, "cannot find such script")
}

func TestListNamespaces(t *testing.T) {
	suite.Run(t, &ListNamespacesTestSuite{})
}

func TestListScripts(t *testing.T) {
	suite.Run(t, &ListScriptsTestSuite{})
}

func TestFindScript(t *testing.T) {
	suite.Run(t, &FindScriptTestSuite{})
}
