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
		suite.NoError(os.MkdirAll(testValue, 0500))

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
	suite.EnsureScript("xx", "1", "echo 1")

	s, err := script.New("xx", "1")
	suite.NoError(err)

	suite.T().Cleanup(func() {
		suite.NoError(os.RemoveAll(s.TempPath()))
	})

	suite.NotEmpty(s.String())
}

func (suite *ScriptTestSuite) TestEnviron() {
	httpmock.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(".*?"),
		httpmock.NewBytesResponder(http.StatusInternalServerError, nil))

	suite.EnsureScript("xx", "1", "echo 1")

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
