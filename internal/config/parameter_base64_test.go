package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterBase64TestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterBase64TestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterBase64TestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewBase64(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterBase64TestSuite) TestType() {
	param, err := config.NewBase64(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterBase64, param.Type())
}

func (suite *ParameterBase64TestSuite) TestString() {
	param, err := config.NewBase64(false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterBase64TestSuite) TestValidation() {
	testTable := map[string]bool{
		"/":    false,
		"QQ+":  false,
		"QQ=":  false,
		"UVFX": true,
	}

	param, err := config.NewBase64(false, map[string]string{
		"encoding": config.Base64EncRawURL,
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
				assert.ErrorContains(t, err, "incorrectly encoded value")
			}
		})
	}
}

func TestParameterBase64(t *testing.T) {
	suite.Run(t, &ParameterBase64TestSuite{})
}
