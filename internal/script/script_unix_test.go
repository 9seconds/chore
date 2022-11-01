//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package script_test

import (
	"os"
	"path/filepath"

	"github.com/9seconds/chore/internal/script"
)

func (suite *ScriptTestSuite) TestScriptNoExecutableBit() {
	suite.createScript("xx", "1", "echo 1")

	path := filepath.Join(suite.ensureNamespace("xx"), "1")
	suite.NoError(os.Chmod(path, 0600))

	_, err := script.New("xx", "1")
	suite.ErrorContains(err, "permission denied")
}
