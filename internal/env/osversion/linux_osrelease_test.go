package osversion_test

import (
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/env/osversion"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LinuxOSReleaseTestSuite struct {
	suite.Suite

	testlib.FixturesTestSuite
}

func (suite *LinuxOSReleaseTestSuite) SetupTest() {
	suite.FixturesTestSuite.Setup(suite.T())
}

func (suite *LinuxOSReleaseTestSuite) FixturePath(path string) string {
	return suite.FixturesTestSuite.FixturePath(filepath.Join("linux", path+".ini"))
}

func (suite *LinuxOSReleaseTestSuite) TestOk() {
	version, err := osversion.ParseLinuxOSRelease(suite.FixturePath("ok"))

	suite.NoError(err)
	suite.Equal("ubuntu", version.ID)
	suite.Equal("focal", version.Codename)
	suite.Equal("20.04", version.Version)
	suite.EqualValues(20, version.Major)
	suite.EqualValues(4, version.Minor)
}

func (suite *LinuxOSReleaseTestSuite) TestFail() {
	testTable := []string{
		"no-codename",
		"no-id",
		"no-version-id",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			_, err := osversion.ParseLinuxOSRelease(suite.FixturePath(testValue))
			assert.ErrorContains(t, err, "cannot find out")
		})
	}
}

func (suite *LinuxOSReleaseTestSuite) TestIncorrectVersion() {
	_, err := osversion.ParseLinuxOSRelease(suite.FixturePath("incorrect-version"))
	suite.ErrorContains(err, "cannot parse version")
}

func (suite *LinuxOSReleaseTestSuite) TestUnknownFile() {
	_, err := osversion.ParseLinuxOSRelease("xxx")
	suite.ErrorContains(err, "cannot load ini file")
}

func TestLinuxOSRelease(t *testing.T) {
	suite.Run(t, &LinuxOSReleaseTestSuite{})
}
