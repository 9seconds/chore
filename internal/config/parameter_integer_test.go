package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterIntegerTestSuite struct {
	suite.Suite
}

func (suite *ParameterIntegerTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			p, err := config.NewInteger(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, p.Required())
		})
	}
}

func (suite *ParameterIntegerTestSuite) TestType() {
	p, err := config.NewInteger(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterInteger, p.Type())
}

func (suite *ParameterIntegerTestSuite) TestString() {
	p, err := config.NewInteger(false, nil)
	suite.NoError(err)
	suite.NotEmpty(p.String())
}

func (suite *ParameterIntegerTestSuite) TestIncorrectValues() {
	p, err := config.NewInteger(false, map[string]string{
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
			assert.Error(t, p.Validate(testValue))
		})
	}
}

func (suite *ParameterIntegerTestSuite) TestCorrectValues() {
	p, err := config.NewInteger(false, map[string]string{
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
			assert.NoError(t, p.Validate(testValue))
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

					_, err := config.NewInteger(false, spec)
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterIntegerTestSuite) TestIncorrectMinMax() {
	_, err := config.NewInteger(false, map[string]string{
		"min": "100",
		"max": "-100",
	})
	suite.Error(err)
}

func TestParameterInteger(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ParameterIntegerTestSuite{})
}
