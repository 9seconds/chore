//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package script_test

import (
	"os"

	"github.com/9seconds/chore/internal/paths"
)

func (suite *ScriptTestSuite) TestScriptNoExecutableBit() {
	suite.EnsureScript("xx", "1", "echo 1")

	path := paths.ConfigNamespaceScript("xx", "1")
	suite.NoError(os.Chmod(path, 0o600))

	scr := suite.NewScript("xx", "1")
	suite.ErrorContains(scr.Init(), "permission denied")
}
