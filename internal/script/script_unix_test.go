//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package script_test

import (
	"os"

	"github.com/9seconds/chore/internal/script"
)

func (suite *ScriptTestSuite) TestScriptNoExecutableBit() {
	suite.EnsureScript("xx", "1", "echo 1")

	path := suite.ConfigScriptPath("xx", "1")
	suite.NoError(os.Chmod(path, 0o600))

	_, err := script.New("xx", "1")
	suite.ErrorContains(err, "permission denied")
}
