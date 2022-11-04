package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/alecthomas/assert/v2"
	"github.com/stretchr/testify/suite"
)

type ParameterFloatTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterFloatTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterFloatTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewFloat(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterFloatTestSuite) TestType() {
	param, err := config.NewFloat(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterFloat, param.Type())
}

func (suite *ParameterFloatTestSuite) TestString() {
	param, err := config.NewInteger(false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterFloatTestSuite) TestIncorrectValues() {
	param, err := config.NewFloat(false, map[string]string{
		"min": "10",
		"max": "20",
	})
	suite.NoError(err)

	testTable := []string{
		"-100",
		"xxx",
		"",
		"--",
		"200",
		"5",
		"0x00",
		"inf",
		"-Inf",
		"nan",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.Error(t, param.Validate(suite.Context(), testValue))
		})
	}
}

func (suite *ParameterFloatTestSuite) TestCorrectValues() {
	param, err := config.NewFloat(false, map[string]string{
		"min": "10",
		"max": "20",
	})
	suite.NoError(err)

	testTable := []string{
		"10",
		"11",
		"19",
		"20",
		"15.0",
		"15.45",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.NoError(t, param.Validate(suite.Context(), testValue))
		})
	}
}

func (suite *ParameterFloatTestSuite) TestIncorrectSpec() {
	testTable := []string{
		"min",
		"max",
	}
	testValues := []string{
		"xxx",
		"",
		"-",
		"nan",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			for _, tv := range testValues {
				tv := tv

				t.Run(tv, func(t *testing.T) {
					spec := map[string]string{}
					spec[testValue] = tv

					_, err := config.NewInteger(false, spec)
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterFloatTestSuite) TestIncorrectMinMax() {
	_, err := config.NewFloat(false, map[string]string{
		"min": "100",
		"max": "-100",
	})
	suite.Error(err)
}

func (suite *ParameterFloatTestSuite) TestAllowInf() {
	_, err := config.NewFloat(false, map[string]string{
		"min": "-inf",
		"max": "inf",
	})
	suite.NoError(err)
}

func TestParameterFloat(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ParameterFloatTestSuite{})
}
