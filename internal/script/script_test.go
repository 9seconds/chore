package script_test

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/script"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/adrg/xdg"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ScriptTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	testlib.ScriptTestSuite
	testlib.NetworkTestSuite
}

func (suite *ScriptTestSuite) SetupTest() {
	t := suite.T()
	suite.CustomRootTestSuite.Setup(t)
	suite.ScriptTestSuite.Setup(t)
	suite.NetworkTestSuite.Setup(t)
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

	suite.EnsureScript("xx", "1", "echo 1")

	for testValue, testName := range testTable {
		suite.NoError(os.MkdirAll(testValue, 0o500))

		suite.T().Run(testName, func(t *testing.T) {
			_, err := script.New("xx", "1")
			assert.ErrorContains(t, err, "permission denied")
		})
	}
}

func (suite *ScriptTestSuite) TestCannotReadConfig() {
	suite.EnsureScript("xx", "1", "echo 1")
	suite.EnsureScriptConfig("xx", "1", "x")

	_, err := script.New("xx", "1")
	suite.ErrorContains(err, "cannot parse config file")
}

func (suite *ScriptTestSuite) TestDirsAreAvailable() {
	suite.EnsureScript("xx", "1", "echo 1")

	scr, err := script.New("xx", "1")
	suite.NoError(err)

	suite.T().Cleanup(func() {
		suite.NoError(os.RemoveAll(scr.TempPath()))
	})

	suite.DirExists(scr.DataPath())
	suite.DirExists(scr.CachePath())
	suite.DirExists(scr.StatePath())
	suite.DirExists(scr.RuntimePath())
	suite.DirExists(scr.TempPath())
}

func (suite *ScriptTestSuite) TestString() {
	suite.EnsureScript("xx", "1", "echo 1")

	scr, err := script.New("xx", "1")
	suite.NoError(err)

	suite.T().Cleanup(func() {
		suite.NoError(os.RemoveAll(scr.TempPath()))
	})

	suite.NotEmpty(scr.String())
}

func (suite *ScriptTestSuite) TestEnviron() {
	httpmock.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(".*?"),
		httpmock.NewBytesResponder(http.StatusInternalServerError, nil))

	suite.EnsureScript("xx", "1", "echo 1")

	scr, err := script.New("xx", "1")
	suite.NoError(err)

	suite.T().Cleanup(func() {
		suite.NoError(os.RemoveAll(scr.TempPath()))
	})

	scr.Config.Network = true
	environ := scr.Environ(context.Background(), argparse.ParsedArgs{
		Keywords: map[string]string{
			"k":  "v",
			"XX": "y",
		},
		Positional: []string{"a", "b", "c"},
	})

	data := map[string]string{}

	for _, v := range environ {
		name, value, found := strings.Cut(v, "=")
		data[name] = value

		require.True(suite.T(), found)
	}

	suite.Len(data, 39)
	suite.Equal(scr.Namespace, data[env.EnvNamespace])
	suite.Equal(scr.Executable, data[env.EnvCaller])
	suite.Equal(scr.Path(), data[env.EnvPathCaller])
	suite.Equal(scr.DataPath(), data[env.EnvPathData])
	suite.Equal(scr.CachePath(), data[env.EnvPathCache])
	suite.Equal(scr.StatePath(), data[env.EnvPathState])
	suite.Equal(scr.RuntimePath(), data[env.EnvPathRuntime])
	suite.Equal(scr.TempPath(), data[env.EnvPathTemp])
	suite.Equal("v", data[env.EnvArgPrefix+"K"])
	suite.Equal("y", data[env.EnvArgPrefix+"XX"])
	suite.Contains(data, env.EnvIDUnique)
	suite.Contains(data, env.EnvIDChainUnique)
	suite.Contains(data, env.EnvIDIsolated)
	suite.Contains(data, env.EnvIDChainIsolated)
	suite.Contains(data, env.EnvMachineID)
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
	suite.Contains(data, env.EnvHostname)
	suite.Contains(data, env.EnvHostnameFQDN)
}

func TestScript(t *testing.T) {
	suite.Run(t, &ScriptTestSuite{})
}
