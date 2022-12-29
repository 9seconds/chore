package config_test

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParameterFileTestSuite struct {
	suite.Suite

	path string

	testlib.CtxTestSuite
	testlib.CustomRootTestSuite
}

func (suite *ParameterFileTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
	suite.CustomRootTestSuite.Setup(suite.T())

	suite.path = suite.EnsureFile(
		filepath.Join(suite.RootPath(), "file"),
		"<html><head></head><body></body></html>",
		0o600)
}

func (suite *ParameterFileTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewFile("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterFileTestSuite) TestType() {
	param, err := config.NewFile("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterFile, param.Type())
}

func (suite *ParameterFileTestSuite) TestString() {
	param, err := config.NewDirectory("", false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterFileTestSuite) TestIncorrectParameter() {
	testTable := map[string][]string{
		"exists":     {"xx", ""},
		"readable":   {"xx", ""},
		"writable":   {"xx", ""},
		"executable": {"xx", ""},
		"mode":       {"[", "xx", "", "-1", "101000000000000000000000"},
		"mimetypes":  {"xx", "xx/yy", "application/???"},
	}

	for testName, testValues := range testTable {
		testName := testName
		testValues := testValues

		suite.T().Run(testName, func(t *testing.T) {
			for _, testValue := range testValues {
				testValue := testValue

				t.Run(testValue, func(t *testing.T) {
					_, err := config.NewFile("", false, map[string]string{
						testName: testValue,
					})
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterFileTestSuite) TestRequireExistButAbsent() {
	param, err := config.NewFile("", false, map[string]string{
		"exists": "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "xx"),
		"does not exist")
}

func (suite *ParameterFileTestSuite) TestDirectory() {
	param, err := config.NewFile("", false, nil)
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), suite.RootPath()),
		"is not a file")
}

func (suite *ParameterFileTestSuite) TestAbsent() {
	param, err := config.NewFile("", false, nil)
	suite.NoError(err)
	suite.NoError(param.Validate(suite.Context(), "aaa"))
}

func (suite *ParameterFileTestSuite) TestWrongMode() {
	param, err := config.NewFile("", false, map[string]string{
		"mode": "055",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), suite.path),
		"incorrect mode")
}

func (suite *ParameterFileTestSuite) TestMimetype() {
	testTable := map[string]bool{
		"application/json":             false,
		"text/plain":                   false,
		"application/json,text/plain":  false,
		"application/x-rar-compressed": false,
		"text/html":                    true,
		"application/json,text/html":   true,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewFile("", false, map[string]string{
				"mimetypes": testName,
			})
			require.NoError(t, err)

			err = param.Validate(suite.Context(), suite.path)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "unexpected mimetype")
			}
		})
	}
}

func TestParameterFile(t *testing.T) {
	suite.Run(t, &ParameterFileTestSuite{})
}
