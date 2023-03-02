package testlib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/paths"
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

	require.NoError(suite.t, paths.EnsureDir(path))

	return path
}

func (suite *CustomRootTestSuite) EnsureFile(path, content string, mode os.FileMode) string {
	suite.t.Helper()

	require.NoError(suite.t, paths.EnsureFile(path, content))
	require.NoError(suite.t, os.Chmod(path, mode))

	return path
}

func (suite *CustomRootTestSuite) EnsureScript(namespace, executable, content string) string {
	suite.t.Helper()

	content = "#!/usr/bin/env bash\nset -eu -o pipefail\n" + content
	path := paths.ConfigNamespaceScript(namespace, executable)

	suite.EnsureFile(path, content, defaultScriptPermission)

	return path
}

func (suite *CustomRootTestSuite) EnsureScriptConfig(namespace, executable string, content interface{}) string {
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

	path := paths.ConfigNamespaceScriptConfig(namespace, executable)

	suite.EnsureFile(path, strContent, defaultConfigPermission)

	return path
}
