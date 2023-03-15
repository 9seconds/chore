package validators_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NamespaceTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	fn cobra.PositionalArgs
}

func (suite *NamespaceTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.fn = validators.Namespace(0)

	suite.EnsureFile(paths.ConfigNamespace("xx"), "", 0o600)
}

func (suite *NamespaceTestSuite) TestUnknownPath() {
	suite.ErrorContains(suite.fn(nil, []string{"aaa"}), "invalid namespace")
}

func (suite *NamespaceTestSuite) TestNotDirectory() {
	suite.ErrorIs(suite.fn(nil, []string{"xx"}), validators.ErrNamespaceIsNotDirectory)
}

type ScriptTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	fn cobra.PositionalArgs
}

func (suite *ScriptTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.fn = validators.Script(0, 1)

	suite.EnsureFile(paths.ConfigNamespace("xx"), "", 0o600)
	suite.EnsureDir(paths.ConfigNamespace("yy"))
	suite.EnsureDir(paths.ConfigNamespaceScriptConfig("zz", "a"))
	suite.EnsureScript("aa", "a", "")
}

func (suite *ScriptTestSuite) TestUnknownNamespace() {
	suite.ErrorContains(
		suite.fn(nil, []string{"cc", "dd"}),
		"no such file or directory")
}

func (suite *ScriptTestSuite) TestUnknownScript() {
	suite.ErrorContains(
		suite.fn(nil, []string{"yy", "dd"}),
		"no such file or directory")
}

func (suite *ScriptTestSuite) TestNotScript() {
	suite.NoError(suite.fn(nil, []string{"aa", "a"}))
}

func (suite *ScriptTestSuite) TestOk() {
	suite.ErrorContains(
		suite.fn(nil, []string{"zz", "a"}),
		"no such file or directory")
}

type ArgumentOptionalTestSuite struct {
	suite.Suite

	fn cobra.PositionalArgs
}

func (suite *ArgumentOptionalTestSuite) SetupTest() {
	suite.fn = validators.ArgumentOptional(1, func(_ *cobra.Command, args []string) error {
		return fmt.Errorf("hello %s", args[1])
	})
}

func (suite *ArgumentOptionalTestSuite) TestEmpty() {
	suite.NoError(suite.fn(nil, nil))
}

func (suite *ArgumentOptionalTestSuite) TestNotEnoughArguments() {
	suite.NoError(suite.fn(nil, []string{"1"}))
}

func (suite *ArgumentOptionalTestSuite) TestEnoughArguments() {
	suite.ErrorContains(suite.fn(nil, []string{"1", "2"}), "hello 2")
}

func TestASCIIName(t *testing.T) {
	err := errors.New("error")
	testable := validators.ASCIIName(0, err)

	testTable := map[string]bool{
		"":       false,
		"a":      true,
		"aa-aa":  true,
		"b_c":    true,
		"привет": false,
	}

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		t.Run(testValue, func(t *testing.T) {
			got := testable(nil, []string{testValue})

			if isValid {
				assert.NoError(t, got)
			} else {
				assert.ErrorIs(t, got, err)
			}
		})
	}
}

func TestArgumentOptional(t *testing.T) {
	suite.Run(t, &ArgumentOptionalTestSuite{})
}

func TestNamespace(t *testing.T) {
	suite.Run(t, &NamespaceTestSuite{})
}

func TestScript(t *testing.T) {
	suite.Run(t, &ScriptTestSuite{})
}
