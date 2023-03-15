package testlib

import (
	"bufio"
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"
)

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

	subcommand  string
	makeCommand func() *cobra.Command
}

func (suite *CobraTestSuite) Setup(t *testing.T, subcommand string, makeCommand func() *cobra.Command) {
	t.Helper()
	suite.CtxTestSuite.Setup(t)
	suite.CustomRootTestSuite.Setup(t)

	suite.subcommand = subcommand
	suite.makeCommand = makeCommand
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
