package testlib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/adrg/xdg"
	"github.com/stretchr/testify/require"
)

const (
	defaultDirPermission    = 0o700
	defaultScriptPermission = 0o700
	defaultConfigPermission = 0o600
)

type CustomRootTestSuite struct {
	fsRoot string
	t      *testing.T
}

func (suite *CustomRootTestSuite) Setup(t *testing.T) {
	t.Helper()

	suite.t = t
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

func (suite *CustomRootTestSuite) RootPath() string {
	return suite.fsRoot
}

func (suite *CustomRootTestSuite) EnsureDir(path string) string {
	suite.t.Helper()

	require.NoError(suite.t, os.MkdirAll(path, defaultDirPermission))

	return path
}

func (suite *CustomRootTestSuite) EnsureFile(
	path, content string, mode os.FileMode,
) string {
	suite.t.Helper()

	suite.EnsureDir(filepath.Dir(path))
	require.NoError(suite.t, os.WriteFile(path, []byte(content), mode))

	return path
}

func (suite *CustomRootTestSuite) ConfigNamespacePath(namespace string) string {
	suite.t.Helper()

	return filepath.Join(env.RootPathConfig(), namespace)
}

func (suite *CustomRootTestSuite) ConfigScriptPath(namespace, executable string) string {
	suite.t.Helper()

	return filepath.Join(suite.ConfigNamespacePath(namespace), executable)
}

func (suite *CustomRootTestSuite) ConfigScriptConfigPath(namespace, executable string) string {
	suite.t.Helper()

	return suite.ConfigScriptPath(namespace, executable) + ".toml"
}

func (suite *CustomRootTestSuite) DataNamespacePath(namespace string) string {
	suite.t.Helper()

	return filepath.Join(env.RootPathData(), namespace)
}

func (suite *CustomRootTestSuite) DataScriptPath(namespace, executable string) string {
	suite.t.Helper()

	return filepath.Join(suite.DataNamespacePath(namespace), executable)
}

func (suite *CustomRootTestSuite) CacheNamespacePath(namespace string) string {
	suite.t.Helper()

	return filepath.Join(env.RootPathCache(), namespace)
}

func (suite *CustomRootTestSuite) CacheScriptPath(namespace, executable string) string {
	suite.t.Helper()

	return filepath.Join(suite.CacheNamespacePath(namespace), executable)
}

func (suite *CustomRootTestSuite) StateNamespacePath(namespace string) string {
	suite.t.Helper()

	return filepath.Join(env.RootPathState(), namespace)
}

func (suite *CustomRootTestSuite) StateScriptPath(namespace, executable string) string {
	suite.t.Helper()

	return filepath.Join(suite.StateNamespacePath(namespace), executable)
}

func (suite *CustomRootTestSuite) RuntimeNamespacePath(namespace string) string {
	suite.t.Helper()

	return filepath.Join(env.RootPathRuntime(), namespace)
}

func (suite *CustomRootTestSuite) RuntimeScriptPath(namespace, executable string) string {
	suite.t.Helper()

	return filepath.Join(suite.RuntimeNamespacePath(namespace), executable)
}

func (suite *CustomRootTestSuite) EnsureScript(namespace, executable, content string) string {
	suite.t.Helper()

	content = "#!/usr/bin/env bash\nset -eu -o pipefail\n" + content
	path := suite.ConfigScriptPath(namespace, executable)

	suite.EnsureFile(
		suite.ConfigScriptPath(namespace, executable),
		content,
		defaultScriptPermission)

	return path
}

func (suite *CustomRootTestSuite) EnsureScriptConfig(namespace, executable string, content interface{}) {
	suite.t.Helper()

	strContent := ""

	switch val := content.(type) {
	case string:
		strContent = val
	case []byte:
		strContent = string(val)
	default:
		data, err := json.Marshal(content)
		strContent = string(data)

		require.NoError(suite.t, err)
	}

	suite.EnsureFile(
		suite.ConfigScriptConfigPath(namespace, executable),
		strContent,
		defaultConfigPermission)
}
