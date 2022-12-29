package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterXMLTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterXMLTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterXMLTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewXML("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterXMLTestSuite) TestType() {
	param, err := config.NewXML("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterXML, param.Type())
}

func (suite *ParameterXMLTestSuite) TestString() {
	param, err := config.NewXML("", false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterXMLTestSuite) TestValidaton() {
	testTable := map[string]bool{
		"":               false,
		"{":              false,
		"[":              false,
		":":              false,
		"{'xx': 1,}":     false,
		`{"xx": 1,}`:     false,
		"{'xx': 1}":      false,
		"[]":             false,
		`{"x": [1,2,3]}`: false,
		"<":              false,
		"<>":             false,
		"<xxx>":          false,
		"<></>":          false,

		"<xxx/>":                true,
		"<xxx></xxx>":           true,
		"<xxx>111</xxx>":        true,
		"<xxx aa='1'>111</xxx>": true,
	}

	param, err := config.NewXML("", false, nil)
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "invalid xml")
			}
		})
	}
}

func TestParameterXML(t *testing.T) {
	suite.Run(t, &ParameterXMLTestSuite{})
}
