package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterEnumTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterEnumTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterEnumTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewEnum("", testValue, map[string]string{
				"choices": "xxx",
			})
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterEnumTestSuite) TestType() {
	param, err := config.NewEnum("", false, map[string]string{
		"choices": "xxx",
	})
	suite.NoError(err)
	suite.Equal(config.ParameterEnum, param.Type())
}

func (suite *ParameterEnumTestSuite) TestString() {
	param, err := config.NewEnum("", false, map[string]string{
		"choices": "xxx",
	})
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterEnumTestSuite) TestNoChoices() {
	_, err := config.NewEnum("", false, map[string]string{
		"choices": ",,,,",
	})
	suite.ErrorContains(err, "no choices are prodvided")
}

func (suite *ParameterEnumTestSuite) TestValidate() {
	testTable := map[string]bool{
		"a":        true,
		"b":        true,
		"c":        true,
		"d":        false,
		"dd":       true,
		"ddd":      false,
		"a,b":      false,
		"a,b,c,dd": false,
	}

	param, err := config.NewEnum("", false, map[string]string{
		"choices": "a,b,c,dd",
	})
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "invalid choice")
			}
		})
	}
}

func TestParameterEnum(t *testing.T) {
	suite.Run(t, &ParameterEnumTestSuite{})
}
