package osversion_test

import (
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/env/osversion"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MacOSPlistTestSuite struct {
	suite.Suite

	testlib.FixturesTestSuite
}

func (suite *MacOSPlistTestSuite) SetupTest() {
	suite.FixturesTestSuite.Setup(suite.T())
}

func (suite *MacOSPlistTestSuite) FixturePath(path string) string {
	return suite.FixturesTestSuite.FixturePath(filepath.Join("macos", path+".xml"))
}

func (suite *MacOSPlistTestSuite) TestOk() {
	version, err := osversion.ParseMacOSPlist(suite.FixturePath("ok"))

	suite.NoError(err)
	suite.Equal("macos", version.ID)
	suite.Equal("monterey", version.Codename)
	suite.Equal("12.6", version.Version)
	suite.EqualValues(12, version.Major)
	suite.EqualValues(6, version.Minor)
}

func (suite *MacOSPlistTestSuite) TestFail() {
	testTable := []string{
		"no-codename",
		"no-version",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			_, err := osversion.ParseMacOSPlist(suite.FixturePath(testValue))
			assert.ErrorContains(t, err, "cannot find out")
		})
	}
}

func (suite *MacOSPlistTestSuite) TestUnknownFile() {
	_, err := osversion.ParseMacOSPlist("xxx")
	suite.ErrorContains(err, "cannot open")
}

func (suite *MacOSPlistTestSuite) TestUnknownFormat() {
	_, err := osversion.ParseMacOSPlist(suite.FixturePath("unknown-format"))
	suite.ErrorContains(err, "cannot parse plist")
}

func (suite *MacOSPlistTestSuite) TestIncorrectSemver() {
	_, err := osversion.ParseMacOSPlist(suite.FixturePath("incorrect-semver"))
	suite.ErrorContains(err, "cannot parse version")
}

func (suite *MacOSPlistTestSuite) TestUnknownVersion() {
	_, err := osversion.ParseMacOSPlist(suite.FixturePath("unknown-version"))
	suite.ErrorContains(err, "unknown major mac os version")
}

func TestMacOSPlist(t *testing.T) {
	suite.Run(t, &MacOSPlistTestSuite{})
}
