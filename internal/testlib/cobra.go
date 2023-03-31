package testlib

import (
	"bufio"
	"bytes"
	"context"
	"testing"

	"github.com/9seconds/chore/internal/cli/base"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

type CobraExitMock struct {
	mock.Mock
}

func (m *CobraExitMock) Exit(code int) {
	m.Called(code)
}

func (m *CobraExitMock) ExitMock(arguments ...any) *mock.Call {
	return m.On("Exit", arguments...)
}

type CobraCommandContext struct {
	context.Context

	Stdout bytes.Buffer
	Stderr bytes.Buffer
	Stdin  bytes.Buffer
}

func (c *CobraCommandContext) StdoutLines() []string {
	return c.scanLines(c.Stdout.Bytes())
}

func (c *CobraCommandContext) StderrLines() []string {
	return c.scanLines(c.Stderr.Bytes())
}

func (c *CobraCommandContext) StdinLines() []string {
	return c.scanLines(c.Stdin.Bytes())
}

func (c *CobraCommandContext) scanLines(buf []byte) []string {
	lines := []string{}
	scanner := bufio.NewScanner(bytes.NewBuffer(buf))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

type CobraTestSuite struct {
	CtxTestSuite
	CustomRootTestSuite

	exitMock    *CobraExitMock
	subcommand  string
	makeCommand func() *cobra.Command
}

func (suite *CobraTestSuite) Setup(t *testing.T, subcommand string, makeCommand func() *cobra.Command) {
	t.Helper()
	suite.CtxTestSuite.Setup(t)
	suite.CustomRootTestSuite.Setup(t)

	suite.exitMock = &CobraExitMock{}
	suite.subcommand = subcommand
	suite.makeCommand = makeCommand

	base.ExitFunc = suite.exitMock.Exit

	t.Cleanup(func() {
		suite.exitMock.AssertExpectations(t)
	})
}

func (suite *CobraTestSuite) ExitMock(arguments ...any) *mock.Call {
	return suite.exitMock.ExitMock(arguments...)
}

func (suite *CobraTestSuite) ExecuteCommand(args ...string) (*CobraCommandContext, error) {
	suite.t.Helper()

	cmd := suite.makeCommand()
	ctx := &CobraCommandContext{
		Context: suite.Context(),
	}

	if suite.subcommand != "" {
		args = append([]string{suite.subcommand}, args...)
	}

	cmd.SetIn(&ctx.Stdin)
	cmd.SetOut(&ctx.Stdout)
	cmd.SetErr(&ctx.Stderr)
	cmd.SetContext(ctx)
	cmd.SetArgs(args)

	return ctx, cmd.Execute()
}
