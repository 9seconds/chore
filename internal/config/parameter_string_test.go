package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterStringTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterStringTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterStringTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewString(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterStringTestSuite) TestType() {
	param, err := config.NewString(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterString, param.Type())
}

func (suite *ParameterStringTestSuite) TestString() {
	param, err := config.NewString(false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterStringTestSuite) TestIncorrectRegexp() {
	_, err := config.NewString(false, map[string]string{
		"regexp": "[",
	})
	suite.Error(err)
}

func (suite *ParameterStringTestSuite) TestIncorrectLength() {
	testNames := []string{
		"min_length",
		"max_length",
	}
	testValues := []string{
		"",
		"-1",
		"x",
		"wdfsladkf1111",
	}

	for _, testName := range testNames {
		testName := testName

		suite.T().Run(testName, func(t *testing.T) {
			for _, testValue := range testValues {
				testValue := testValue

				suite.T().Run(testValue, func(t *testing.T) {
					_, err := config.NewString(false, map[string]string{
						testName: testValue,
					})
					assert.ErrorContains(t, err, "incorrect "+testName+" value")
				})
			}
		})
	}
}

func (suite *ParameterStringTestSuite) TestInvalidValues() {
	param, err := config.NewString(false, map[string]string{
		"regexp": `^xx\w{2}\d`,
	})
	suite.NoError(err)

	testTable := []string{
		"xx",
		"",
		"xxaa",
		"xx11x",
		"xxaax",
		"yyaa1",
		"xxxaa1",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.Error(t, param.Validate(suite.Context(), testValue))
		})
	}
}

func (suite *ParameterStringTestSuite) TestValidValues() {
	param, err := config.NewString(false, map[string]string{
		"regexp": `xx\w{2}\d`,
	})
	suite.NoError(err)

	testTable := []string{
		"xxaa1",
		"xxba2",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.NoError(t, param.Validate(suite.Context(), testValue))
		})
	}
}

func (suite *ParameterStringTestSuite) TestStringLength() {
	testTable := map[string]bool{
		"":        false,
		"x":       false,
		"xx":      false,
		"xxx":     true,
		"xxxx":    true,
		"xxxxx":   true,
		"xxxxxxx": false,
	}

	param, err := config.NewString(false, map[string]string{
		"regexp":     `x*`,
		"min_length": "3",
		"max_length": "5",
	})
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "value length must be")
			}
		})
	}
}

func TestParameterString(t *testing.T) {
	suite.Run(t, &ParameterStringTestSuite{})
}
