package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterBoolTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterBoolTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterBoolTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewBool(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterBoolTestSuite) TestType() {
	param, err := config.NewBool(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterBool, param.Type())
}

func (suite *ParameterBoolTestSuite) TestString() {
	param, err := config.NewBool(false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterBoolTestSuite) TestValues() {
	testTable := map[string]bool{
		"":      false,
		"True":  true,
		"true":  true,
		"1":     true,
		"yes":   false,
		"0":     true,
		"no":    false,
		"False": true,
		"false": true,
		"xx":    false,
		"2":     false,
	}

	param, err := config.NewBool(false, nil)
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "invalid syntax")
			}
		})
	}
}

func TestParameterBool(t *testing.T) {
	suite.Run(t, &ParameterBoolTestSuite{})
}
