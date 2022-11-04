package commands_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/script"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OSTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
	testlib.CustomRootTestSuite
	testlib.ScriptTestSuite

	s       script.Script
	args    []string
	environ []string

	stdin  io.Reader
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func (suite *OSTestSuite) SetupTest() {
	t := suite.T()

	suite.CtxTestSuite.Setup(t)
	suite.CustomRootTestSuite.Setup(t)
	suite.ScriptTestSuite.Setup(t)

	suite.EnsureScript("x", "y", "echo $CHORE_CALLER $1")

	scr, err := suite.NewScript("x", "y")
	require.NoError(t, err)

	parsedArgs := argparse.ParsedArgs{
		Keywords: map[string]string{
			"k": "v",
		},
		Positional: []string{"a", "b"},
	}

	suite.s = scr
	suite.environ = scr.Environ(suite.Context(), parsedArgs)
	suite.args = parsedArgs.Positional
	suite.stdout = &bytes.Buffer{}
	suite.stderr = &bytes.Buffer{}

	stdin, err := os.Open(os.DevNull)
	require.NoError(t, err)

	t.Cleanup(func() {
		stdin.Close()
	})

	suite.stdin = stdin
}

func (suite *OSTestSuite) TestExecuteCommand() {
	cmd := commands.NewOS(
		suite.Context(),
		suite.s,
		suite.environ,
		suite.args,
		suite.stdin,
		suite.stdout,
		suite.stderr)

	suite.Equal(0, cmd.Pid())

	suite.NoError(cmd.Start())
	suite.NotEqual(0, cmd.Pid())

	result, err := cmd.Wait()
	suite.NoError(err)
	suite.Equal(0, result.ExitCode)
	suite.Less(result.UserTime, time.Second)
	suite.Less(result.SystemTime, time.Second)
	suite.Less(result.ElapsedTime, time.Second)
	suite.Empty(suite.stderr.String())
	suite.Equal("y a\n", suite.stdout.String())
}

func (suite *OSTestSuite) TestExitCode() {
	cmd := commands.NewOS(
		suite.Context(),
		suite.s,
		suite.environ,
		suite.args,
		suite.stdin,
		suite.stdout,
		suite.stderr)

	suite.EnsureScript("x", "y", "exit 3")

	suite.NoError(cmd.Start())
	result, err := cmd.Wait()
	suite.NoError(err)

	suite.Equal(3, result.ExitCode)
}

func TestOs(t *testing.T) {
	suite.Run(t, &OSTestSuite{})
}
