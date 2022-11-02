package script_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/script"
	"github.com/adrg/xdg"
	"github.com/jarcoal/httpmock"
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

func (suite *ScriptTestSuite) TestEnviron() {
	httpmock.ActivateNonDefault(env.HTTPClientV4)
	httpmock.ActivateNonDefault(env.HTTPClientV6)
	httpmock.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(".*?"),
		httpmock.NewBytesResponder(http.StatusInternalServerError, nil))
	suite.T().Cleanup(httpmock.DeactivateAndReset)

	suite.createScript("xx", "1", "echo 1")

	s, err := script.New("xx", "1")
	suite.NoError(err)

	suite.T().Cleanup(func() {
		suite.NoError(os.RemoveAll(s.TempPath()))
	})

	s.Config.Network = true
	environ := s.Environ(context.Background(), argparse.ParsedArgs{
		Keywords: map[string]string{
			"k":  "v",
			"XX": "y",
		},
		Positional: []string{"a", "b", "c"},
	})

	data := map[string]string{}

	for _, v := range environ {
		name, value, found := strings.Cut(v, "=")
		require.True(suite.T(), found)
		data[name] = value
	}

	suite.Len(data, 30)
	suite.Equal(s.Namespace, data[env.EnvNamespace])
	suite.Equal(s.Executable, data[env.EnvCaller])
	suite.Equal(s.Path(), data[env.EnvPathCaller])
	suite.Equal(s.DataPath(), data[env.EnvPathData])
	suite.Equal(s.CachePath(), data[env.EnvPathCache])
	suite.Equal(s.StatePath(), data[env.EnvPathState])
	suite.Equal(s.RuntimePath(), data[env.EnvPathRuntime])
	suite.Equal(s.TempPath(), data[env.EnvPathTemp])
	suite.Equal("v", data[env.EnvArgPrefix+"K"])
	suite.Equal("y", data[env.EnvArgPrefix+"XX"])
	suite.Contains(data, env.EnvIdUnique)
	suite.Contains(data, env.EnvIdChainUnique)
	suite.Contains(data, env.EnvIdIsolated)
	suite.Contains(data, env.EnvIdChainIsolated)
	suite.Contains(data, env.EnvMachineId)
	suite.Contains(data, env.EnvStartedAtRFC3339)
	suite.Contains(data, env.EnvStartedAtUnix)
	suite.Contains(data, env.EnvStartedAtYear)
	suite.Contains(data, env.EnvStartedAtYearDay)
	suite.Contains(data, env.EnvStartedAtDay)
	suite.Contains(data, env.EnvStartedAtMonth)
	suite.Contains(data, env.EnvStartedAtMonthStr)
	suite.Contains(data, env.EnvStartedAtHour)
	suite.Contains(data, env.EnvStartedAtMinute)
	suite.Contains(data, env.EnvStartedAtSecond)
	suite.Contains(data, env.EnvStartedAtNanosecond)
	suite.Contains(data, env.EnvStartedAtTimezone)
	suite.Contains(data, env.EnvStartedAtOffset)
	suite.Contains(data, env.EnvStartedAtWeekday)
	suite.Contains(data, env.EnvStartedAtWeekdayStr)
}

func TestScript(t *testing.T) {
	suite.Run(t, &ScriptTestSuite{})
}
