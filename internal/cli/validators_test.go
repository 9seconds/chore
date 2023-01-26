package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ValidNamespaceTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	fn cobra.PositionalArgs
}

func (suite *ValidNamespaceTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.fn = validNamespace(0)

	suite.EnsureFile(suite.ConfigNamespacePath("xx"), "", 0o600)
}

func (suite *ValidNamespaceTestSuite) TestUnknownPath() {
	suite.ErrorContains(suite.fn(nil, []string{"aaa"}), "invalid namespace")
}

func (suite *ValidNamespaceTestSuite) TestNotDirectory() {
	suite.ErrorIs(suite.fn(nil, []string{"xx"}), ErrNamespaceIsNotDirectory)
}

type ValidScriptTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	fn cobra.PositionalArgs
}

func (suite *ValidScriptTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.fn = validScript(0, 1)

	suite.EnsureFile(suite.ConfigNamespacePath("xx"), "", 0o600)
	suite.EnsureDir(suite.ConfigNamespacePath("yy"))
	suite.EnsureDir(suite.ConfigScriptConfigPath("zz", "a"))
	suite.EnsureScript("aa", "a", "")
}

func (suite *ValidScriptTestSuite) TestUnknownNamespace() {
	suite.ErrorContains(
		suite.fn(nil, []string{"cc", "dd"}),
		"no such file or directory")
}

func (suite *ValidScriptTestSuite) TestUnknownScript() {
	suite.ErrorContains(
		suite.fn(nil, []string{"yy", "dd"}),
		"no such file or directory")
}

func (suite *ValidScriptTestSuite) TestNotScript() {
	suite.NoError(suite.fn(nil, []string{"aa", "a"}))
}

func (suite *ValidScriptTestSuite) TestOk() {
	suite.ErrorContains(
		suite.fn(nil, []string{"zz", "a"}),
		"no such file or directory")
}

func TestArgumentOptional(t *testing.T) {
	testable := argumentOptional(1, func(_ *cobra.Command, args []string) error {
		return fmt.Errorf("hello %s", args[1])
	})

	t.Run("empty", func(t *testing.T) {
		assert.NoError(t, testable(nil, nil))
	})

	t.Run("not-enough-arguments", func(t *testing.T) {
		assert.NoError(t, testable(nil, []string{"1"}))
	})

	t.Run("enough-arguments", func(t *testing.T) {
		assert.ErrorContains(t, testable(nil, []string{"1", "2"}), "hello 2")
	})
}

func TestValidAsciiName(t *testing.T) {
	err := errors.New("error!")
	testable := validAsciiName(0, err)

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

func TestValidNamespace(t *testing.T) {
	suite.Run(t, &ValidNamespaceTestSuite{})
}

func TestValidScript(t *testing.T) {
	suite.Run(t, &ValidScriptTestSuite{})
}
