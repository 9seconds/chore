package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterSemverTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterSemverTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterSemverTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewSemver("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterSemverTestSuite) TestType() {
	param, err := config.NewSemver("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterSemver, param.Type())
}

func (suite *ParameterSemverTestSuite) TestIncorrectSemver() {
	param, err := config.NewSemver("", false, nil)
	suite.NoError(err)

	testTable := map[string]bool{
		"x":      false,
		"":       false,
		"v":      false,
		"1.":     false,
		"1.0":    true,
		"1":      true,
		"v2.3":   true,
		"v2.3.4": true,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect semver")
			}
		})
	}
}

func (suite *ParameterSemverTestSuite) TestValidateConstraint() {
	param, err := config.NewSemver("", false, map[string]string{
		"constraint": "~1.2.3",
	})
	suite.NoError(err)

	testTable := map[string]bool{
		"1.0":    false,
		"v1.2.3": true,
		"v1.2.4": true,
		"1.4.0":  false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "invalid version")
			}
		})
	}
}

func TestParameterSemver(t *testing.T) {
	suite.Run(t, &ParameterSemverTestSuite{})
}
