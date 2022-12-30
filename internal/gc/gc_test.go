package gc_test

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/gc"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/suite"
)

type GCTestSuite struct {
	suite.Suite
	testlib.CustomRootTestSuite
}

func (suite *GCTestSuite) EqualStrings(expected, actual []string) {
	sort.Strings(expected)
	sort.Strings(actual)
	suite.Equal(expected, actual)
}

func (suite *GCTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.EnsureScript("x", "valid_script_with_config", "echo 1")
	suite.EnsureScriptConfig("x", "valid_script_with_config", "{}")

	suite.EnsureScript("x", "valid_script_without_config", "echo 2")

	suite.EnsureScript("x", "valid_script_with_incorrect_config", "echo 2")
	suite.EnsureScriptConfig("x", "valid_script_with_incorrect_config", "{")

	suite.EnsureDir(suite.ConfigScriptConfigPath("y", "script_config_dir"))
	suite.EnsureDir(suite.ConfigScriptPath("y", "script_dir"))
	suite.EnsureDir(suite.DataNamespacePath("y1"))
	suite.EnsureFile(suite.DataNamespacePath("y2"), "", 0o600)
	suite.EnsureFile(suite.CacheNamespacePath("y2"), "", 0o600)
	suite.EnsureDir(suite.CacheScriptPath("y", "script_dir"))
	suite.EnsureDir(
		filepath.Join(suite.StateScriptPath("x", "valid_script_without_config"), "a"),
	)
	suite.EnsureDir(suite.CacheScriptPath("x", "valid_script_with_config"))
}

func (suite *GCTestSuite) TestCollect() {
	filenames, err := gc.Collect()
	suite.NoError(err)

	suite.EqualStrings([]string{
		suite.CacheNamespacePath("y"),
		suite.CacheNamespacePath("y2"),
		suite.ConfigScriptConfigPath("x", "valid_script_with_incorrect_config"),
		suite.ConfigNamespacePath("y"),
		suite.DataNamespacePath("y1"),
		suite.DataNamespacePath("y2"),
	}, filenames)
}

func (suite *GCTestSuite) TestRemove() {
	filenames, err := gc.Collect()
	suite.NoError(err)

	suite.NoError(gc.Remove(filenames))
	suite.DirExists(env.RootPathConfig())
	suite.DirExists(env.RootPathState())
	suite.DirExists(env.RootPathCache())
	suite.DirExists(env.RootPathData())
	suite.NoDirExists(suite.DataNamespacePath("y1"))
	suite.NoDirExists(suite.DataNamespacePath("y2"))
	suite.NoDirExists(suite.CacheNamespacePath("y1"))
	suite.NoDirExists(suite.CacheNamespacePath("y2"))
	suite.NoDirExists(suite.ConfigNamespacePath("y"))
	suite.NoFileExists(suite.ConfigScriptConfigPath("x", "valid_script_without_config"))
}

func TestGC(t *testing.T) {
	suite.Run(t, &GCTestSuite{})
}
