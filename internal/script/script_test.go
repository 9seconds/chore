package script_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/script"
	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ScriptTestSuite struct {
	suite.Suite

	fsRoot string
}

func (suite *ScriptTestSuite) SetupTest() {
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

func (suite *ScriptTestSuite) buildPath(root, namespace string) string {
	return filepath.Join(root, env.ChoreDir, namespace)
}

func (suite *ScriptTestSuite) ensureNamespace(namespace string) string {
	base := suite.buildPath(xdg.ConfigHome, namespace)
	require.NoError(suite.T(), os.MkdirAll(base, 0700))

	return base
}

func (suite *ScriptTestSuite) createScript(namespace, executable, content string) {
	base := suite.ensureNamespace(namespace)
	err := os.WriteFile(
		filepath.Join(base, executable),
		[]byte("#!/usr/bin/env sh\n"+content),
		0700)
	require.NoError(suite.T(), err)
}

func (suite *ScriptTestSuite) createConfig(namespace, executable string, content interface{}) {
	base := suite.ensureNamespace(namespace)
	data, err := json.Marshal(content)
	suite.NoError(err)

	err = os.WriteFile(filepath.Join(base, executable+".json"), data, 0600)
	require.NoError(suite.T(), err)
}

func (suite *ScriptTestSuite) TestAbsentScript() {
	_, err := script.New("xx", "1")
	suite.Error(err)
}

func (suite *ScriptTestSuite) TestCannotCreatePath() {
	testTable := map[string]string{
		xdg.DataHome:   "data",
		xdg.CacheHome:  "cache",
		xdg.StateHome:  "state",
		xdg.RuntimeDir: "runtime",
	}

	suite.createScript("xx", "1", "echo 1")

	for testValue, testName := range testTable {
		suite.NoError(os.MkdirAll(testValue, 0500))

		suite.T().Run(testName, func(t *testing.T) {
			_, err := script.New("xx", "1")
			assert.ErrorContains(t, err, "permission denied")
		})
	}
}

func (suite *ScriptTestSuite) TestCannotReadConfig() {
	suite.createScript("xx", "1", "echo 1")
	suite.createConfig("xx", "1", "x")

	_, err := script.New("xx", "1")
	suite.ErrorContains(err, "cannot parse config file")
}

func (suite *ScriptTestSuite) TestDirsAreAvailable() {
	suite.createScript("xx", "1", "echo 1")

	s, err := script.New("xx", "1")
	suite.NoError(err)

	suite.T().Cleanup(func() {
		suite.NoError(os.RemoveAll(s.TempPath()))
	})

	suite.DirExists(s.DataPath())
	suite.DirExists(s.CachePath())
	suite.DirExists(s.StatePath())
	suite.DirExists(s.RuntimePath())
	suite.DirExists(s.TempPath())
}

func (suite *ScriptTestSuite) TestString() {
	suite.createScript("xx", "1", "echo 1")

	s, err := script.New("xx", "1")
	suite.NoError(err)

	suite.T().Cleanup(func() {
		suite.NoError(os.RemoveAll(s.TempPath()))
	})

	suite.NotEmpty(s.String())
}

func TestScript(t *testing.T) {
	suite.Run(t, &ScriptTestSuite{})
}
