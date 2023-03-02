package gc_test

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/9seconds/chore/internal/gc"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GCTestSuite struct {
	suite.Suite
	testlib.CustomRootTestSuite

	validScripts []*script.Script
}

func (suite *GCTestSuite) EqualStrings(expected, actual []string) {
	sort.Strings(expected)
	sort.Strings(actual)
	suite.Equal(expected, actual)
}

func (suite *GCTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.EnsureScript("x", "valid_script_with_config", "echo 1")
	suite.EnsureScriptConfig("x", "valid_script_with_config", "description = '1'")

	suite.EnsureScript("x", "valid_script_without_config", "echo 2")

	suite.EnsureScript("x", "valid_script_with_incorrect_config", "echo 2")
	suite.EnsureScriptConfig("x", "valid_script_with_incorrect_config", "{")

	suite.EnsureDir(paths.ConfigNamespaceScriptConfig("y", "script_config_dir"))
	suite.EnsureDir(paths.ConfigNamespaceScript("y", "script_dir"))
	suite.EnsureDir(paths.DataNamespace("y1"))
	suite.EnsureFile(paths.DataNamespace("y2"), "", 0o600)
	suite.EnsureFile(paths.CacheNamespace("y2"), "", 0o600)
	suite.EnsureDir(paths.CacheNamespaceScript("y", "script_dir"))
	suite.EnsureDir(
		filepath.Join(paths.StateNamespaceScript("x", "valid_script_without_config"), "a"),
	)
	suite.EnsureDir(paths.CacheNamespaceScript("x", "valid_script_with_config"))

	namespaces, err := script.ListNamespaces()
	require.NoError(suite.T(), err)

	suite.validScripts = nil

	for _, namespace := range namespaces {
		scripts, err := script.ListScripts(namespace)
		require.NoError(suite.T(), err)

		for _, name := range scripts {
			suite.validScripts = append(suite.validScripts, &script.Script{
				Namespace:  namespace,
				Executable: name,
			})
		}
	}

	require.NotEmpty(suite.T(), suite.validScripts)
}

func (suite *GCTestSuite) TestCollect() {
	filenames, err := gc.Collect(suite.validScripts)
	suite.NoError(err)

	suite.EqualStrings([]string{
		paths.CacheNamespace("y"),
		paths.CacheNamespace("y2"),
		paths.ConfigNamespaceScriptConfig("x", "valid_script_with_incorrect_config"),
		paths.ConfigNamespace("y"),
		paths.DataNamespace("y1"),
		paths.DataNamespace("y2"),
	}, filenames)
}

func (suite *GCTestSuite) TestRemove() {
	filenames, err := gc.Collect(suite.validScripts)
	suite.NoError(err)

	suite.NoError(gc.Remove(filenames))
	suite.DirExists(paths.ConfigRoot())
	suite.DirExists(paths.StateRoot())
	suite.DirExists(paths.CacheRoot())
	suite.DirExists(paths.DataRoot())
	suite.NoDirExists(paths.DataNamespace("y1"))
	suite.NoDirExists(paths.DataNamespace("y2"))
	suite.NoDirExists(paths.CacheNamespace("y1"))
	suite.NoDirExists(paths.CacheNamespace("y2"))
	suite.NoDirExists(paths.ConfigNamespace("y"))
	suite.NoFileExists(paths.ConfigNamespaceScriptConfig("x", "valid_script_without_config"))
}

func TestGC(t *testing.T) {
	suite.Run(t, &GCTestSuite{})
}
