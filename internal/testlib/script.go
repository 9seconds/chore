package testlib

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/script"
	"github.com/stretchr/testify/require"
)

type ScriptTestSuite struct {
	t *testing.T
}

func (suite *ScriptTestSuite) Setup(t *testing.T) {
	suite.t = t
}

func (suite *ScriptTestSuite) NewScript(namespace, executable string) (script.Script, error) {
	s, err := script.New(namespace, executable)

	if err == nil {
		suite.t.Cleanup(func() {
			require.NoError(suite.t, os.RemoveAll(s.TempPath()))
		})
	}

	return s, err
}
