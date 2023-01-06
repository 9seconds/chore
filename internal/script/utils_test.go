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

func (suite *ListTestSuite) TestListNamespaces() {
	namespaces, err := script.ListNamespaces()
	suite.NoError(err)
	suite.Equal([]string{"nb", "ns"}, namespaces)
}

func (suite *ListTestSuite) TestListScripts() {
	scripts, err := script.ListScripts("nb")
	suite.NoError(err)
	suite.Equal([]string{"aa", "ab"}, scripts)
}

func (suite *ListTestSuite) TestListScriptsNothing() {
	scripts, err := script.ListScripts("ns")
	suite.NoError(err)
	suite.Empty(scripts)
}

func TestList(t *testing.T) {
	suite.Run(t, &ListTestSuite{})
}
