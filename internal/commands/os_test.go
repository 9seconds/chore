package commands_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/script"
	"github.com/adrg/xdg"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OSTestSuite struct {
	suite.Suite

	ctx       context.Context
	ctxCancel context.CancelFunc
	s         script.Script
	args      argparse.ParsedArgs

	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (suite *OSTestSuite) SetupTest() {
	t := suite.T()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	suite.ctx = ctx
	suite.ctxCancel = cancel
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

	outR, outW, err := os.Pipe()
	require.NoError(t, err)

	errR, errW, err := os.Pipe()
	require.NoError(t, err)

	go func() {
		io.Copy(&suite.stdout, outR)
	}()

	go func() {
		io.Copy(&suite.stderr, errR)
	}()

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	t.Cleanup(func() {
		outW.Close()
		outR.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	})

	os.Stdout = outW
	os.Stderr = errW

	dir := t.TempDir()

	xdg.ConfigHome = filepath.Join(dir, "config_home")
	xdg.DataHome = filepath.Join(dir, "data_home")
	xdg.CacheHome = filepath.Join(dir, "cache_home")
	xdg.StateHome = filepath.Join(dir, "state_home")
	xdg.RuntimeDir = filepath.Join(dir, "runtime_dir")

	s := script.Script{
		Namespace:  "x",
		Executable: "y",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(s.Path()), 0700))
	require.NoError(t, os.WriteFile(
		s.Path(),
		[]byte("#!/bin/sh\necho $CHORE_CALLER $1"),
		0700))

	s, err = script.New("x", "y")
	require.NoError(t, err)

	suite.s = s
	suite.args = argparse.ParsedArgs{
		Keywords: map[string]string{
			"k": "v",
		},
		Positional: []string{"a", "b"},
	}
}

func (suite *OSTestSuite) WriteScriptContent(content string) {
	content = "#!/usr/bin/env bash\n" + content

	require.NoError(suite.T(), os.WriteFile(
		suite.s.Path(),
		[]byte(content),
		0700))
}

func (suite *OSTestSuite) TestExecuteCommand() {
	cmd := commands.NewOS(suite.ctx, suite.s, suite.args)

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
	cmd := commands.NewOS(suite.ctx, suite.s, suite.args)
	suite.WriteScriptContent("exit 3")

	suite.NoError(cmd.Start())
	result, err := cmd.Wait()
	suite.NoError(err)

	suite.Equal(3, result.ExitCode)
}

func TestOs(t *testing.T) {
	suite.Run(t, &OSTestSuite{})
}
