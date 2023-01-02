package config_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/alecthomas/assert/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ParameterIPTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
	testlib.NetworkTestSuite
}

func (suite *ParameterIPTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
	suite.NetworkTestSuite.Setup(suite.T())
}

func (suite *ParameterIPTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewIP("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterIPTestSuite) TestType() {
	param, err := config.NewIP("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterIP, param.Type())
}

func (suite *ParameterIPTestSuite) TestIncorrectParameter() {
	testTable := map[string][]string{
		"resolve":           {"xx", ""},
		"allowed_subnets":   {",,,xx", "::/256", ",,127.0.0.5", "127.0.0.0/8,11"},
		"forbidden_subnets": {",,,xx", "::/256", ",,127.0.0.5", "127.0.0.0/8,11"},
	}

	for testName, testValues := range testTable {
		testName := testName
		testValues := testValues

		suite.T().Run(testName, func(t *testing.T) {
			for _, testValue := range testValues {
				testValue := testValue

				t.Run(testValue, func(t *testing.T) {
					_, err := config.NewIP("", false, map[string]string{
						testName: testValue,
					})
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterIPTestSuite) TestValidateForbiddenSubnet() {
	testTable := map[string]bool{
		"127.0.0.1":  false,
		"127.0.1.10": false,
		"cafe::2":    true,
		"10.0.0.0":   true,
		"xxx":        false,
		"":           false,
	}

	param, err := config.NewIP("", false, map[string]string{
		"forbidden_subnets": "127.0.0.0/8,127.0.1.0/24",
	})
	suite.NoError(err)

	for testName, testValue := range testTable {
		testName := testName
		testValue := testValue

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if testValue {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func (suite *ParameterIPTestSuite) TestValidateAllowedSubnet() {
	testTable := map[string]bool{
		"127.0.0.1":  true,
		"127.0.1.10": true,
		"cafe::2":    false,
		"10.0.0.0":   false,
		"10.0.1.10":  true,
		"xxx":        false,
		"":           false,
	}

	param, err := config.NewIP("", false, map[string]string{
		"allowed_subnets": "127.0.0.0/8,10.0.1.0/24",
	})
	suite.NoError(err)

	for testName, testValue := range testTable {
		testName := testName
		testValue := testValue

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if testValue {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func (suite *ParameterIPTestSuite) TestValidateResolveFailed() {
	suite.DNS().
		On("LookupAddr", mock.Anything, "10.0.1.10").
		Once().
		Return([]string(nil), io.EOF)

	param, err := config.NewIP("", false, map[string]string{
		"allowed_subnets": "127.0.0.0/8,10.0.1.0/24",
		"resolve":         "true",
	})
	suite.NoError(err)

	suite.ErrorContains(
		param.Validate(suite.Context(), "10.0.1.10"),
		"cannot do reverse lookup")
}

func (suite *ParameterIPTestSuite) TestValidateResolveOk() {
	suite.DNS().
		On("LookupAddr", mock.Anything, "10.0.1.10").
		Once().
		Return([]string{"xxx"}, nil)

	param, err := config.NewIP("", false, map[string]string{
		"allowed_subnets": "127.0.0.0/8,10.0.1.0/24",
		"resolve":         "true",
	})
	suite.NoError(err)
	suite.NoError(param.Validate(suite.Context(), "10.0.1.10"))
}

func TestParameterIP(t *testing.T) {
	suite.Run(t, &ParameterIPTestSuite{})
}
