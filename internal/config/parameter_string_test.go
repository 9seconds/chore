package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterStringTestSuite struct {
	suite.Suite
}

func (suite *ParameterStringTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			p, err := config.NewString(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, p.Required())
		})
	}
}

func (suite *ParameterStringTestSuite) TestType() {
	p, err := config.NewString(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterString, p.Type())
}

func (suite *ParameterStringTestSuite) TestString() {
	p, err := config.NewString(false, nil)
	suite.NoError(err)
	suite.NotEmpty(p.String())
}

func (suite *ParameterStringTestSuite) TestIncorrectRegexp() {
	_, err := config.NewString(false, map[string]string{
		"regexp": "[",
	})
	suite.Error(err)
}

func (suite *ParameterStringTestSuite) TestInvalidValues() {
	p, err := config.NewString(false, map[string]string{
		"regexp": `xx\w{2}\d`,
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
			assert.Error(t, p.Validate(testValue))
		})
	}
}

func (suite *ParameterStringTestSuite) TestValidValues() {
	p, err := config.NewString(false, map[string]string{
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
			assert.NoError(t, p.Validate(testValue))
		})
	}
}

func TestParameterString(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ParameterStringTestSuite{})
}
