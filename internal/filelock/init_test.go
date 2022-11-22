package filelock_test

import (
	"path/filepath"

	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/suite"
)

type BaseLockTestSuite struct {
	suite.Suite
	testlib.CustomRootTestSuite
	testlib.CtxTestSuite

	path string
}

func (suite *BaseLockTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())
	suite.CtxTestSuite.Setup(suite.T())

	suite.path = suite.EnsureFile(
		filepath.Join(suite.RootPath(), "t"),
		"xx",
		0o600)
}
