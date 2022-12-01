package testlib

import (
	"testing"

	"github.com/9seconds/chore/internal/script"
)

type ScriptTestSuite struct {
	t *testing.T
}

func (suite *ScriptTestSuite) Setup(t *testing.T) {
	t.Helper()

	suite.t = t
}

func (suite *ScriptTestSuite) NewScript(namespace, executable string) *script.Script {
	scr := &script.Script{
		Namespace:  namespace,
		Executable: executable,
	}

	suite.t.Cleanup(scr.Cleanup)

	return scr
}
