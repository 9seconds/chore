package cli_test

import (
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/testlib"
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

func (suite *LockTestSuite) TestPath() {
	suite.NoError(suite.param.UnmarshalText([]byte(suite.path)))
	suite.Equal(suite.path, suite.param.Value("cc"))
}

func (suite *LockTestSuite) TestMagicValue() {
	suite.NoError(suite.param.UnmarshalText([]byte(cli.MagicValue)))
	suite.Equal("cc", suite.param.Value("cc"))
}

func TestLock(t *testing.T) {
	suite.Run(t, &LockTestSuite{})
}
