package testlib

import (
	"testing"

	"github.com/9seconds/chore/internal/script"
	"github.com/alecthomas/assert/v2"
)

type ScriptTestSuite struct {
	t *testing.T
}

func (suite *ScriptTestSuite) Setup(t *testing.T) {
	t.Helper()

	suite.t = t
}

func (suite *ScriptTestSuite) NewScript(namespace, executable string) (*script.Script, error) {
	scr, err := script.New(namespace, executable)

	if err == nil {
		suite.t.Cleanup(func() {
			assert.NoError(suite.t, scr.Cleanup())
		})
	}

	return scr, err
}
