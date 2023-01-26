package cli_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
	"github.com/stretchr/testify/suite"
)

type CmdShowTestSuite struct {
	CmdTestSuite
}

func (suite *CmdShowTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("show", cli.NewShow)

	suite.EnsureScriptConfig("ns", "s", `
description = "ZZY"
git = "always"
network = true

[flags.flag1]
description = "This is a description for flag1"
required = false  # default value

[parameters.param]
description = "Never knows best"
type = "string"
required = false

[paramters.param.spec]
ascii = true
regexp = '^\d\w+$'`)
	suite.EnsureScript("ns", "s", "")
	suite.EnsureScript("ns", "s2", "")
	suite.EnsureScript("xx", "s3", "")
	suite.EnsureFile(
		filepath.Join(suite.CacheScriptPath("ns", "s"), "c"),
		strings.Repeat("a", 50*1024*1024),
		0o600)
}

func (suite *CmdShowTestSuite) TestNoArguments() {
	ctx, err := suite.ExecuteCommand(nil)

	suite.NoError(err)
	suite.Empty(ctx.StderrLines())
	suite.Equal([]string{"ns", "xx"}, ctx.StdoutLines())
}

func (suite *CmdShowTestSuite) TestListNamespace() {
	ctx, err := suite.ExecuteCommand([]string{"ns"})

	suite.NoError(err)
	suite.Empty(ctx.StderrLines())
	suite.Equal([]string{"s", "s2"}, ctx.StdoutLines())
}

func (suite *CmdShowTestSuite) TestShow() {
	ctx, err := suite.ExecuteCommand([]string{"ns", "s"})

	scr := &script.Script{
		Namespace:  "ns",
		Executable: "s",
	}

	seen := map[string]bool{
		scr.Path():        true,
		scr.ConfigPath():  true,
		scr.DataPath():    true,
		scr.CachePath():   true,
		scr.StatePath():   true,
		scr.RuntimePath(): true,
		"network":         true,
		"git":             true,
		"param":           true,
		"flag1":           true,
	}

	suite.NoError(err)
	suite.Empty(ctx.StderrLines())
	suite.Contains(ctx.Stdout.String(), "ZZY")

	for _, line := range ctx.StdoutLines() {
		switch {
		case strings.Contains(line, scr.ConfigPath()):
			suite.Contains(seen, scr.ConfigPath())
			delete(seen, scr.ConfigPath())
		case strings.Contains(line, scr.Path()):
			suite.Contains(seen, scr.Path())
			delete(seen, scr.Path())
		case strings.Contains(line, scr.DataPath()):
			suite.Contains(seen, scr.DataPath())
			delete(seen, scr.DataPath())
			suite.Contains(line, "4KB")
		case strings.Contains(line, scr.CachePath()):
			suite.Contains(seen, scr.CachePath())
			delete(seen, scr.CachePath())
			suite.Contains(line, "50MB")
		case strings.Contains(line, scr.StatePath()):
			suite.Contains(seen, scr.StatePath())
			delete(seen, scr.StatePath())
			suite.Contains(line, "4KB")
		case strings.Contains(line, scr.RuntimePath()):
			suite.Contains(seen, scr.RuntimePath())
			delete(seen, scr.RuntimePath())
			suite.Contains(line, "4KB")
		case strings.Contains(strings.ToLower(line), "network"):
			suite.Contains(seen, "network")
			delete(seen, "network")
			suite.Contains(line, "true")
		case strings.Contains(strings.ToLower(line), "git"):
			suite.Contains(seen, "git")
			delete(seen, "git")
			suite.Contains(line, "always")
		case strings.Contains(line, "param"):
			suite.Contains(seen, "param")
			delete(seen, "param")
			suite.Contains(line, "Never knows best")
			suite.Contains(line, "string")
		case strings.Contains(line, "flag1"):
			suite.Contains(seen, "flag1")
			delete(seen, "flag1")
			suite.Contains(line, "This is a description for flag1")
		}
	}

	suite.Empty(seen)
}

func TestCmdShow(t *testing.T) {
	suite.Run(t, &CmdShowTestSuite{})
}
