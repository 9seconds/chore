package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterHexTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterHexTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterHexTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewHex(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterHexTestSuite) TestType() {
	param, err := config.NewHex(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterHex, param.Type())
}

func (suite *ParameterHexTestSuite) TestString() {
	param, err := config.NewHex(false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterHexTestSuite) TestIncorrectLength() {
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
					_, err := config.NewHex(false, map[string]string{
						testName: testValue,
					})
					assert.ErrorContains(t, err, "incorrect '"+testName+"' value")
				})
			}
		})
	}
}

func (suite *ParameterHexTestSuite) TestValidaton() {
	testTable := map[string]bool{
		"X":  false,
		"XX": false,
		"1":  false,
		"11": true,
		"":   true,
		"AB": true,
		"A":  false,
		"a":  false,
		"AX": false,
		"aa": true,
		"BB": true,
	}

	param, err := config.NewHex(false, nil)
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrectly encoded hex value")
			}
		})
	}
}

func (suite *ParameterHexTestSuite) TestStringLengthValidation() {
	testTable := map[string]bool{
		"":       false,
		"AA":     true,
		"CCCC":   true,
		"ABCDEF": false,
	}

	param, err := config.NewHex(false, map[string]string{
		"min_length": "1",
		"max_length": "4",
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

func TestParameterHex(t *testing.T) {
	suite.Run(t, &ParameterHexTestSuite{})
}
