package script_test

import (
	"context"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/git"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/Showmax/go-fqdn"
	"github.com/adrg/xdg"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ScriptTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	testlib.NetworkTestSuite
}

func (suite *ScriptTestSuite) SetupTest() {
	t := suite.T()

	suite.CustomRootTestSuite.Setup(t)
	suite.NetworkTestSuite.Setup(t)
}

func (suite *ScriptTestSuite) TestAbsentScript() {
	_, err := script.New("xx", "1")
	suite.Error(err)
}

func (suite *ScriptTestSuite) TestCannotCreatePath() {
	testTable := map[string]string{
		xdg.DataHome:  "data",
		xdg.CacheHome: "cache",
		xdg.StateHome: "state",
	}

	suite.EnsureScript("xx", "1", "echo 1")

	for testValue, testName := range testTable {
		suite.NoError(os.MkdirAll(testValue, 0o500))

		suite.T().Run(testName, func(t *testing.T) {
			scr, err := script.New("xx", "1")
			assert.NoError(t, err)

			assert.ErrorContains(t, scr.EnsureDirs(), "permission denied")
		})
	}
}

func (suite *ScriptTestSuite) TestCannotReadConfig() {
	suite.EnsureScript("xx", "1", "echo 1")
	suite.EnsureScriptConfig("xx", "1", "x")

	_, err := script.New("xx", "1")
	suite.ErrorContains(err, "cannot parse config file")
}

func (suite *ScriptTestSuite) TestEmptyScript() {
	suite.EnsureFile(
		paths.ConfigNamespaceScript("xx", "1"),
		"\t \r\n",
		0o700)

	_, err := script.New("xx", "1")
	suite.ErrorContains(err, "script is empty")
}

func (suite *ScriptTestSuite) TestDirsAreAvailable() {
	suite.EnsureScript("xx", "1", "echo 1")

	scr, err := script.New("xx", "1")
	suite.NoError(err)

	suite.NoDirExists(scr.DataPath())
	suite.NoDirExists(scr.CachePath())
	suite.NoDirExists(scr.StatePath())
	suite.NoDirExists(scr.TempPath())
	suite.NoFileExists(scr.ConfigPath())

	suite.NoError(scr.EnsureDirs())
	suite.DirExists(scr.DataPath())
	suite.DirExists(scr.CachePath())
	suite.DirExists(scr.StatePath())
	suite.DirExists(scr.TempPath())
	suite.NoFileExists(scr.ConfigPath())
}

func (suite *ScriptTestSuite) TestDoNotRecreateTempPath() {
	suite.EnsureScript("xx", "1", "echo 1")

	scr, err := script.New("xx", "1")
	suite.NoError(err)

	suite.NoError(scr.EnsureDirs())

	tmpPath1 := scr.TempPath()

	suite.NoError(scr.EnsureDirs())
	suite.Equal(tmpPath1, scr.TempPath())
}

func (suite *ScriptTestSuite) TestString() {
	suite.EnsureScript("xx", "1", "echo 1")

	scr, err := script.New("xx", "1")
	suite.NoError(err)
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

	scr.Config.Network = true
	scr.Config.Git = git.AccessModeIfUndefined

	environ := scr.Environ(context.Background(), argparse.ParsedArgs{
		Parameters: map[string]string{
			"k":  "v",
			"XX": "y",
		},
		Flags: map[string]argparse.FlagValue{
			"cleanup": argparse.FlagTrue,
			"welcome": argparse.FlagFalse,
		},
		Positional: []string{"a", "b", "c"},
	})

	data := map[string]string{}

	for _, v := range environ {
		name, value, found := strings.Cut(v, "=")
		data[name] = value

		require.True(suite.T(), found)
	}

	count := 45

	suite.Equal(scr.Namespace, data[env.EnvNamespace])
	suite.Equal(scr.Executable, data[env.EnvCaller])
	suite.Equal(scr.Path(), data[env.EnvPathCaller])
	suite.Equal(scr.DataPath(), data[env.EnvPathData])
	suite.Equal(scr.CachePath(), data[env.EnvPathCache])
	suite.Equal(scr.StatePath(), data[env.EnvPathState])
	suite.Equal(scr.TempPath(), data[env.EnvPathTemp])
	suite.Equal("v", data[env.ParameterName("k")])
	suite.Equal("y", data[env.ParameterName("XX")])
	suite.EqualValues(argparse.FlagTrue, data[env.FlagName("CLEANUP")])
	suite.EqualValues(argparse.FlagFalse, data[env.FlagName("WELCOME")])
	suite.Contains(data, env.EnvRecursion)
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
	suite.Contains(data, env.EnvGitReference)
	suite.Contains(data, env.EnvGitReferenceShort)
	suite.Contains(data, env.EnvGitReferenceType)
	suite.Contains(data, env.EnvGitCommitHash)
	suite.Contains(data, env.EnvGitCommitHashShort)
	suite.Contains(data, env.EnvGitIsDirty)

	if value, err := os.Hostname(); err == nil {
		suite.Equal(value, data[env.EnvHostname])

		count++
	}

	if value, err := fqdn.FqdnHostname(); err == nil {
		suite.Equal(value, data[env.EnvHostnameFQDN])

		count++
	}

	if value, err := user.Current(); err == nil {
		suite.Equal(value.Username, data[env.EnvUserName])

		count += 3
	}

	suite.Len(data, count)
}

func TestScript(t *testing.T) {
	suite.Run(t, &ScriptTestSuite{})
}
