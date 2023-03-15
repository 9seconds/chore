package edit

import (
	"os"
	"os/exec"
	"testing"

	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/suite"
)

type EnsureFileTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
}

func (suite *EnsureFileTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.EnsureDir(paths.ConfigRoot())
}

func (suite *EnsureFileTestSuite) TestCreate() {
	path := paths.ConfigNamespace("xx")

	suite.NoError(ensureFile(path, []byte{1, 2, 3}))

	data, err := os.ReadFile(path)
	suite.NoError(err)
	suite.Equal([]byte{1, 2, 3}, data)
}

func (suite *EnsureFileTestSuite) TestFileExists() {
	path := paths.ConfigNamespace("xx")
	suite.EnsureFile(path, "aaa", ConfigDefaultPermission)

	suite.NoError(ensureFile(path, []byte{1, 2, 3}))

	data, err := os.ReadFile(path)
	suite.NoError(err)
	suite.Equal("aaa", string(data))
}

func TestEnsureFile(t *testing.T) {
	suite.Run(t, &EnsureFileTestSuite{})
}

type OpenEditorTestSuite struct {
	suite.Suite

	testlib.CustomRootTestSuite
	testlib.CtxTestSuite
}

func (suite *OpenEditorTestSuite) SetupTest() {
	suite.CustomRootTestSuite.Setup(suite.T())
	suite.CtxTestSuite.Setup(suite.T())

	suite.EnsureFile(paths.ConfigNamespace("aa"), "aa", ConfigDefaultPermission)
}

func (suite *OpenEditorTestSuite) TestEditedOk() {
	path := paths.ConfigNamespace("aa")

	suite.NoError(openEditor(suite.Context(), "true", path))

	data, err := os.ReadFile(path)

	suite.NoError(err)
	suite.Equal("aa", string(data))
}

func (suite *OpenEditorTestSuite) TestEditedFailed() {
	path := paths.ConfigNamespace("aa")

	suite.ErrorContains(openEditor(suite.Context(), "false", path), "command exited with 1")

	data, err := os.ReadFile(path)

	suite.NoError(err)
	suite.Equal("aa", string(data))
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
