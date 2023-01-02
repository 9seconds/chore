package config_test

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/alecthomas/assert/v2"
	"github.com/stretchr/testify/suite"
)

type ParameterDirectoryTestSuite struct {
	suite.Suite

	dir string

	testlib.CtxTestSuite
	testlib.CustomRootTestSuite
}

func (suite *ParameterDirectoryTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.dir = suite.EnsureDir(filepath.Join(suite.RootPath(), "dir"))
}

func (suite *ParameterDirectoryTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewDirectory("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterDirectoryTestSuite) TestType() {
	param, err := config.NewDirectory("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterDirectory, param.Type())
}

func (suite *ParameterDirectoryTestSuite) TestIncorrectParameter() {
	testTable := map[string][]string{
		"exists":     {"xx", ""},
		"readable":   {"xx", ""},
		"writable":   {"xx", ""},
		"executable": {"xx", ""},
		"mode":       {"[", "xx", "", "-1", "101000000000000000000000"},
	}

	for testName, testValues := range testTable {
		testName := testName
		testValues := testValues

		suite.T().Run(testName, func(t *testing.T) {
			for _, testValue := range testValues {
				testValue := testValue

				t.Run(testValue, func(t *testing.T) {
					_, err := config.NewDirectory("", false, map[string]string{
						testName: testValue,
					})
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterDirectoryTestSuite) TestExistRequiredButAbsent() {
	param, err := config.NewDirectory("", false, map[string]string{
		"exists": "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "xxx"),
		"does not exist")
}

func (suite *ParameterDirectoryTestSuite) TestAbsent() {
	param, err := config.NewDirectory("", false, map[string]string{
		"exists": "false",
	})
	suite.NoError(err)
	suite.NoError(param.Validate(suite.Context(), "xxx"))
}

func (suite *ParameterDirectoryTestSuite) TestIsNotDir() {
	param, err := config.NewDirectory("", false, map[string]string{
		"exists": "false",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), suite.EnsureScript("x", "y", "")),
		"is not a directory")
}

func (suite *ParameterDirectoryTestSuite) TestWrongMode() {
	param, err := config.NewDirectory("", false, map[string]string{
		"mode": "055",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), suite.dir),
		"incorrect mode")
}

func TestParameterDirectory(t *testing.T) {
	suite.Run(t, &ParameterDirectoryTestSuite{})
}
