package cli

import (
	"os"
	"os/exec"
	"testing"

	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/suite"
)

type OpenEditorTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	testlib.CtxTestSuite
}

func (suite *OpenEditorTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())
	suite.CtxTestSuite.Setup(suite.T())

	suite.EnsureDir(paths.ConfigNamespace("xx"))
}

func (suite *OpenEditorTestSuite) TestCannotWriteDirectory() {
	suite.ErrorIs(
		openEditor(
			suite.Context(),
			"true",
			paths.ConfigNamespace("xx"),
			[]byte{1, 2, 3}),
		ErrExpectedFile)
}

func (suite *OpenEditorTestSuite) TestEditedOk() {
	path := paths.ConfigNamespace("aa")

	suite.NoError(openEditor(suite.Context(), "true", path, []byte{1, 2, 3}))

	data, err := os.ReadFile(path)

	suite.NoError(err)
	suite.Equal([]byte{1, 2, 3}, data)
}

func (suite *OpenEditorTestSuite) TestEditedFailed() {
	path := paths.ConfigNamespace("aa")

	suite.ErrorContains(
		openEditor(suite.Context(), "false", path, []byte{1, 2, 3}),
		"command exited with 1")

	data, err := os.ReadFile(path)

	suite.NoError(err)
	suite.Equal([]byte{1, 2, 3}, data)
}

func TestOpenEditor(t *testing.T) {
	if _, err := exec.LookPath("true"); err != nil {
		t.Skipf("do not have true in /bin: %v", err)
	}

	if _, err := exec.LookPath("false"); err != nil {
		t.Skipf("do not have false in /bin: %v", err)
	}

	suite.Run(t, &OpenEditorTestSuite{})
}
