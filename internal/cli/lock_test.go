package cli_test

import (
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/filelock"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LockTestSuite struct {
	suite.Suite
	testlib.CustomRootTestSuite

	param cli.Lock
	path  string
}

func (suite *LockTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.path = suite.EnsureFile(
		filepath.Join(suite.RootPath(), "xx"),
		"a",
		0o600)
}

func (suite *LockTestSuite) TestUnknownPath() {
	suite.ErrorContains(
		suite.param.UnmarshalText([]byte("aa")),
		"cannot stat path")
}

func (suite *LockTestSuite) TestLockMode() {
	testTable := map[string]filelock.LockType{
		suite.path:            filelock.LockTypeExclusive,
		"x:" + suite.path:     filelock.LockTypeExclusive,
		"s:" + suite.path:     filelock.LockTypeShared,
		cli.MagicValue:        filelock.LockTypeExclusive,
		"x:" + cli.MagicValue: filelock.LockTypeExclusive,
		"s:" + cli.MagicValue: filelock.LockTypeShared,
	}

	for testName, lockType := range testTable {
		testName := testName
		lockType := lockType

		suite.T().Run(testName, func(t *testing.T) {
			param := cli.Lock{}
			assert.NoError(t, param.UnmarshalText([]byte(testName)))
			assert.Equal(t, lockType, param.LockMode())
		})
	}
}

func (suite *LockTestSuite) TestPath() {
	testTable := map[string]string{
		suite.path:            suite.path,
		"x:" + suite.path:     suite.path,
		"s:" + suite.path:     suite.path,
		cli.MagicValue:        "aaa",
		"x:" + cli.MagicValue: "aaa",
		"s:" + cli.MagicValue: "aaa",
	}

	for testName, expectedPath := range testTable {
		testName := testName
		expectedPath := expectedPath

		suite.T().Run(testName, func(t *testing.T) {
			param := cli.Lock{}
			assert.NoError(t, param.UnmarshalText([]byte(testName)))
			assert.Equal(t, expectedPath, param.Path("aaa"))
		})
	}
}

func (suite *LockTestSuite) TestDefault() {
	suite.Equal(filelock.LockTypeNo, suite.param.LockMode())
	suite.Equal("u", suite.param.Path("u"))
}

func TestLock(t *testing.T) {
	suite.Run(t, &LockTestSuite{})
}
