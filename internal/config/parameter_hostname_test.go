package config_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ParameterHostnameTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
	testlib.NetworkTestSuite
}

func (suite *ParameterHostnameTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
	suite.NetworkTestSuite.Setup(suite.T())
}

func (suite *ParameterHostnameTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewHostname(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterHostnameTestSuite) TestType() {
	param, err := config.NewHostname(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterHostname, param.Type())
}

func (suite *ParameterHostnameTestSuite) TestString() {
	param, err := config.NewHostname(false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterHostnameTestSuite) TestIncorrectParameter() {
	testTable := map[string][]string{
		"is_fqdn":    {"xx", ""},
		"resolve":    {"xx", ""},
		"regexp":     {"["},
		"min_length": {"-1", "xx", ""},
		"max_length": {"-1", "xx", ""},
	}

	for testName, testValues := range testTable {
		testName := testName
		testValues := testValues

		suite.T().Run(testName, func(t *testing.T) {
			for _, testValue := range testValues {
				testValue := testValue

				t.Run(testValue, func(t *testing.T) {
					_, err := config.NewHostname(false, map[string]string{
						testName: testValue,
					})
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterHostnameTestSuite) TestValidateHostname() {
	testTable := map[string]bool{
		"*":         false,
		"test.test": true,
		"xx":        true,
		"localhost": true,
		"xx-xx/xx":  false,
	}

	param, err := config.NewHostname(false, nil)
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect hostname")
			}
		})
	}
}

func (suite *ParameterHostnameTestSuite) TestValidateFQDN() {
	testTable := map[string]bool{
		"*":           false,
		"test.test":   true,
		"xx":          false,
		"localhost":   false,
		"xx-xx/xx":    false,
		"localhost.":  false,
		"company.com": true,
	}

	param, err := config.NewHostname(false, map[string]string{
		"is_fqdn": "true",
	})
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect hostname")
			}
		})
	}
}

func (suite *ParameterHostnameTestSuite) TestValidateRegexp() {
	testTable := map[string]bool{
		"*":           false,
		"test.test":   false,
		"xx":          true,
		"localhost":   false,
		"xx-xx/xx":    false,
		"localhost.":  false,
		"company.com": false,
	}

	param, err := config.NewHostname(false, map[string]string{
		"regexp": "x+",
	})
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func (suite *ParameterHostnameTestSuite) TestResolveDNSFailure() {
	suite.DNS().
		On("LookupHost", mock.Anything, "xx").
		Once().
		Return([]string{}, io.EOF)

	param, err := config.NewHostname(false, map[string]string{
		"resolve": "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "xx"),
		"cannot resolve dns records")
}

func (suite *ParameterHostnameTestSuite) TestResolveDNSNoRecords() {
	suite.DNS().
		On("LookupHost", mock.Anything, "xx").
		Once().
		Return([]string{}, nil)

	param, err := config.NewHostname(false, map[string]string{
		"resolve": "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "xx"),
		"no dns records")
}

func (suite *ParameterHostnameTestSuite) TestResolveOk() {
	suite.DNS().
		On("LookupHost", mock.Anything, "xx").
		Once().
		Return([]string{"xx"}, nil)

	param, err := config.NewHostname(false, map[string]string{
		"resolve": "true",
	})
	suite.NoError(err)
	suite.NoError(param.Validate(suite.Context(), "xx"))
}

func TestParameterHostname(t *testing.T) {
	suite.Run(t, &ParameterHostnameTestSuite{})
}
