package paths_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/paths"
	"github.com/adrg/xdg"
	"github.com/stretchr/testify/suite"
)

type EnsuresTestSuite struct {
	suite.Suite

	fsRoot string
}

func (suite *EnsuresTestSuite) SetupTest() {
	t := suite.T()

	suite.fsRoot = t.TempDir()
	t.Setenv("TMPDIR", suite.fsRoot)

	oldConfigHome := xdg.ConfigHome
	oldDataHome := xdg.DataHome
	oldCacheHome := xdg.CacheHome
	oldStateHome := xdg.StateHome
	oldRuntimeDir := xdg.RuntimeDir

	t.Cleanup(func() {
		xdg.ConfigHome = oldConfigHome
		xdg.DataHome = oldDataHome
		xdg.CacheHome = oldCacheHome
		xdg.StateHome = oldStateHome
		xdg.RuntimeDir = oldRuntimeDir
	})

	xdg.ConfigHome = filepath.Join(suite.fsRoot, "config_home")
	xdg.DataHome = filepath.Join(suite.fsRoot, "data_home")
	xdg.CacheHome = filepath.Join(suite.fsRoot, "cache_home")
	xdg.StateHome = filepath.Join(suite.fsRoot, "state_home")
	xdg.RuntimeDir = filepath.Join(suite.fsRoot, "runtime_dir")
}

func (suite *EnsuresTestSuite) TestEnsureDir() {
	path := filepath.Join(paths.ConfigRoot(), "xx")

	suite.NoError(paths.EnsureDir(path))
	suite.DirExists(path)

	stat, err := os.Stat(path)
	suite.NoError(err)
	suite.Equal(paths.DirectoryPermission, stat.Mode().Perm())
}

func (suite *EnsuresTestSuite) TestEnsureFile() {
	path := filepath.Join(paths.ConfigRoot(), "xx")

	suite.NoError(paths.EnsureFile(path, "aaa"))
	suite.FileExists(path)

	stat, err := os.Stat(path)
	suite.NoError(err)
	suite.Equal(paths.FilePermission, stat.Mode().Perm())

	content, err := os.ReadFile(path)
	suite.NoError(err)
	suite.Equal("aaa", string(content))
}

func (suite *EnsuresTestSuite) TestEnsureRoots() {
	suite.NoError(paths.EnsureRoots("xx", "yy"))
	suite.DirExists(paths.ConfigNamespace("xx"))
	suite.NoDirExists(paths.ConfigNamespaceScript("xx", "yy"))
	suite.NoDirExists(paths.ConfigNamespaceScriptConfig("xx", "yy"))
	suite.DirExists(paths.StateNamespace("xx"))
	suite.DirExists(paths.StateNamespaceScript("xx", "yy"))
	suite.DirExists(paths.CacheNamespace("xx"))
	suite.DirExists(paths.CacheNamespaceScript("xx", "yy"))
	suite.DirExists(paths.DataNamespace("xx"))
	suite.DirExists(paths.DataNamespaceScript("xx", "yy"))

	content, err := os.ReadFile(filepath.Join(paths.CacheRoot(), paths.CacheDirTagName))
	suite.NoError(err)
	suite.Equal(paths.CacheDirTagContent, string(content))
}

func TestEnsures(t *testing.T) {
	suite.Run(t, &EnsuresTestSuite{})
}
