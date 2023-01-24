package cli

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FlagEditorTestSuite struct {
	suite.Suite

	flag *flagEditor
}

func (suite *FlagEditorTestSuite) SetupTest() {
	suite.flag = &flagEditor{}
}

func (suite *FlagEditorTestSuite) Setenv(key, value string) {
	suite.T().Setenv(key, value)
}

func (suite *FlagEditorTestSuite) TestType() {
	suite.Equal("executable", suite.flag.Type())
}

func (suite *FlagEditorTestSuite) TestOverride() {
	edPath, err := exec.LookPath("ed")
	if err != nil {
		suite.T().Skip("cannot detect ed")
	}

	suite.NoError(suite.flag.Set("ed"))

	value, err := suite.flag.Get()

	suite.NoError(err)
	suite.Equal(edPath, value)
}

func (suite *FlagEditorTestSuite) TestVisualPriority() {
	suite.Setenv("VISUAL", "visual")
	suite.Setenv("EDITOR", "edit")

	value, err := suite.flag.Get()

	suite.NoError(err)
	suite.Equal("visual", value)
}

func (suite *FlagEditorTestSuite) TestEditor() {
	suite.Setenv("EDITOR", "edit")

	value, err := suite.flag.Get()

	suite.NoError(err)
	suite.Equal("edit", value)
}

func (suite *FlagEditorTestSuite) TestFallback() {
	value, err := suite.flag.Get()

	suite.NoError(err)
	suite.NotEmpty(value)
}

func (suite *FlagEditorTestSuite) TestSetUnknownEditor() {
	suite.Error(suite.flag.Set("xxxxxxxxxxxxxx"))
}

func TestFlagEditor(t *testing.T) {
	suite.Run(t, &FlagEditorTestSuite{})
}
