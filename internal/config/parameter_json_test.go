package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterJSONTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterJSONTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterJSONTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewJSON("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterJSONTestSuite) TestType() {
	param, err := config.NewJSON("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterJSON, param.Type())
}

func (suite *ParameterJSONTestSuite) TestValidaton() {
	testTable := map[string]bool{
		"":           false,
		"{":          false,
		"[":          false,
		":":          false,
		"{'xx': 1,}": false,
		`{"xx": 1,}`: false,
		"{'xx': 1}":  false,

		"[]":             true,
		`{"x": [1,2,3]}`: true,
	}

	param, err := config.NewJSON("", false, nil)
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "invalid json")
			}
		})
	}
}

func TestParameterJSON(t *testing.T) {
	suite.Run(t, &ParameterJSONTestSuite{})
}
