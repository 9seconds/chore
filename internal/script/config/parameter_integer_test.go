package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterIntegerTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterIntegerTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterIntegerTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewInteger("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterIntegerTestSuite) TestType() {
	param, err := config.NewInteger("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterInteger, param.Type())
}

func (suite *ParameterIntegerTestSuite) TestIncorrectValues() {
	param, err := config.NewInteger("", false, map[string]string{
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
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.Error(t, param.Validate(suite.Context(), testValue))
		})
	}
}

func (suite *ParameterIntegerTestSuite) TestCorrectValues() {
	param, err := config.NewInteger("", false, map[string]string{
		"min": "10",
		"max": "20",
	})
	suite.NoError(err)

	testTable := []string{
		"10",
		"11",
		"19",
		"20",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.NoError(t, param.Validate(suite.Context(), testValue))
		})
	}
}

func (suite *ParameterIntegerTestSuite) TestIncorrectSpec() {
	testTable := []string{
		"min",
		"max",
	}
	testValues := []string{
		"xxx",
		"",
		"-",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			for _, tv := range testValues {
				tv := tv

				t.Run(tv, func(t *testing.T) {
					spec := map[string]string{}
					spec[testValue] = tv

					_, err := config.NewInteger("", false, spec)
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterIntegerTestSuite) TestIncorrectMinMax() {
	_, err := config.NewInteger("", false, map[string]string{
		"min": "100",
		"max": "-100",
	})
	suite.Error(err)
}

func TestParameterInteger(t *testing.T) {
	suite.Run(t, &ParameterIntegerTestSuite{})
}
