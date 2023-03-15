package script_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
}

func (suite *ListTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.EnsureDir(paths.ConfigNamespace("ns"))
	suite.EnsureScript("nb", "aa", "echo 1")
	suite.EnsureScript("nb", "ab", "echo 1")
	suite.EnsureFile(paths.ConfigNamespaceScript("nb", "a"), "1", 0o400)
	suite.EnsureDir(paths.ConfigNamespaceScript("nb", "b"))
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

func TestExtractRealNamespace(t *testing.T) {
	if _, ok := os.LookupEnv(env.EnvNamespace); ok {
		t.Skip("environment variable is defined")
	}

	namespace, exists := script.ExtractRealNamespace("ns")
	assert.True(t, exists)
	assert.Equal(t, namespace, "ns")

	_, exists = script.ExtractRealNamespace(".")
	assert.False(t, exists)

	t.Setenv(env.EnvNamespace, "xx")

	namespace, exists = script.ExtractRealNamespace(".")
	assert.True(t, exists)
	assert.Equal(t, namespace, "xx")

	namespace, exists = script.ExtractRealNamespace("ns")
	assert.True(t, exists)
	assert.Equal(t, namespace, "ns")
}
