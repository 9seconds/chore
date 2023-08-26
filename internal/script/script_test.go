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
		Parameters: map[string][]string{
			"k":  {"v", "u"},
			"XX": {"y"},
		},
		Flags: map[string]bool{
			"cleanup": true,
			"welcome": false,
		},
		Positional: []string{"a", "b", "c"},
	})

	data := map[string]string{}

	for _, v := range environ {
		name, value, found := strings.Cut(v, "=")
		data[name] = value

		require.True(suite.T(), found)
	}

	count := 48

	suite.Equal(scr.Namespace, data[env.Namespace])
	suite.Equal(scr.Executable, data[env.Caller])
	suite.Contains(data, env.Bin)
	suite.Equal(scr.Path(), data[env.PathCaller])
	suite.Equal(scr.DataPath(), data[env.PathData])
	suite.Equal(scr.CachePath(), data[env.PathCache])
	suite.Equal(scr.StatePath(), data[env.PathState])
	suite.Equal(scr.TempPath(), data[env.PathTemp])
	suite.Equal("v\nu", data[env.ParameterNameList("k")])
	suite.Equal("y", data[env.ParameterNameList("XX")])
	suite.Equal("u", data[env.ParameterName("k")])
	suite.Equal("y", data[env.ParameterName("XX")])
	suite.EqualValues(argparse.FlagEnabled, data[env.FlagName("CLEANUP")])
	suite.NotContains(data, env.FlagName("WELCOME"))
	suite.Contains(data, env.Self)
	suite.Contains(data, env.Slug)
	suite.Contains(data, env.IDRun)
	suite.Contains(data, env.IDChainRun)
	suite.Contains(data, env.IDIsolated)
	suite.Contains(data, env.IDChainIsolated)
	suite.Contains(data, env.MachineID)
	suite.Contains(data, env.StartedAtRFC3339)
	suite.Contains(data, env.StartedAtUnix)
	suite.Contains(data, env.StartedAtYear)
	suite.Contains(data, env.StartedAtYearDay)
	suite.Contains(data, env.StartedAtDay)
	suite.Contains(data, env.StartedAtMonth)
	suite.Contains(data, env.StartedAtMonthStr)
	suite.Contains(data, env.StartedAtHour)
	suite.Contains(data, env.StartedAtMinute)
	suite.Contains(data, env.StartedAtSecond)
	suite.Contains(data, env.StartedAtNanosecond)
	suite.Contains(data, env.StartedAtTimezone)
	suite.Contains(data, env.StartedAtOffset)
	suite.Contains(data, env.StartedAtWeekday)
	suite.Contains(data, env.StartedAtWeekdayStr)
	suite.Contains(data, env.GitReference)
	suite.Contains(data, env.GitReferenceShort)
	suite.Contains(data, env.GitReferenceType)
	suite.Contains(data, env.GitCommitHash)
	suite.Contains(data, env.GitCommitHashShort)
	suite.Contains(data, env.GitIsDirty)

	if value, err := os.Hostname(); err == nil {
		suite.Equal(value, data[env.Hostname])

		count++
	}

	if value, err := fqdn.FqdnHostname(); err == nil {
		suite.Equal(value, data[env.HostnameFQDN])

		count++
	}

	if value, err := user.Current(); err == nil {
		suite.Equal(value.Username, data[env.UserName])

		count += 3
	}

	suite.Len(data, count)
}

func TestScript(t *testing.T) {
	suite.Run(t, &ScriptTestSuite{})
}
