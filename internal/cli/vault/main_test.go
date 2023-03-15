package vault_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli/edit"
	"github.com/9seconds/chore/internal/cli/vault"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type VaultTestSuite struct {
	suite.Suite

	testlib.CobraTestSuite
}

func (suite *VaultTestSuite) SetupTest() {
	suite.CobraTestSuite.Setup(suite.T(), "", func() *cobra.Command {
		cmd := &cobra.Command{}

		cmd.AddCommand(
			vault.NewList(),
			vault.NewDelete(),
			vault.NewSet(),
			vault.NewGet())

		return cmd
	})

	suite.EnsureScript("ns", "s", "echo 1")
	suite.EnsureScript("xx", "s", "echo 1")
	suite.EnsureScript("z", "s", "echo 1")
	suite.EnsureFile(paths.AppConfigPath(), `
[vault]
z = ""
ns = "xxx"`, edit.ConfigDefaultPermission)
}

func (suite *VaultTestSuite) TestUnknownPassword() {
	_, err := suite.ExecuteCommand("list", "xx")
	suite.ErrorContains(err, "cannot find out correct password")
}

func (suite *VaultTestSuite) TestEmptyPassword() {
	_, err := suite.ExecuteCommand("list", "z")
	suite.ErrorContains(err, "password is empty")
}

func (suite *VaultTestSuite) TestCRUD() {
	ctx, err := suite.ExecuteCommand("list", "ns")
	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	ctx, err = suite.ExecuteCommand("get", "ns", "k")
	suite.ErrorIs(err, vault.ErrKeyUnknown)
	suite.Contains(ctx.Stderr.String(), err.Error())

	ctx, err = suite.ExecuteCommand("set", "ns", "k", "v")
	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	ctx, err = suite.ExecuteCommand("get", "ns", "k")
	suite.NoError(err)
	suite.Contains(ctx.StdoutLines(), "v")
	suite.Empty(ctx.StderrLines())

	ctx, err = suite.ExecuteCommand("delete", "ns", "k", "k", "k2")
	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	ctx, err = suite.ExecuteCommand("list", "ns")
	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	ctx, err = suite.ExecuteCommand("get", "ns", "k")
	suite.ErrorIs(err, vault.ErrKeyUnknown)
	suite.Contains(ctx.Stderr.String(), err.Error())
}

func TestVault(t *testing.T) {
	suite.Run(t, &VaultTestSuite{})
}
