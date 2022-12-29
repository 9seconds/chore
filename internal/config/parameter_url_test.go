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

type ParameterURLTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
	testlib.NetworkTestSuite
}

func (suite *ParameterURLTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
	suite.NetworkTestSuite.Setup(suite.T())
}

func (suite *ParameterURLTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewURL("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterURLTestSuite) TestType() {
	param, err := config.NewURL("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterURL, param.Type())
}

func (suite *ParameterURLTestSuite) TestString() {
	param, err := config.NewURL("", false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterStringTestSuite) TestIncorrectParameter() {
	testTable := map[string][]string{
		"resolve":   {"xx", ""},
		"domain_re": {"["},
		"path_re":   {"["},
		"user_re":   {"["},
	}

	for testName, testValues := range testTable {
		testName := testName
		testValues := testValues

		suite.T().Run(testName, func(t *testing.T) {
			for _, testValue := range testValues {
				testValue := testValue

				t.Run(testValue, func(t *testing.T) {
					_, err := config.NewURL("", false, map[string]string{
						testName: testValue,
					})
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterURLTestSuite) TestValidateDomain() {
	testTable := map[string]bool{
		"https://amazon.de":     false,
		"https://google.ru.com": true,
		"https://google.com":    true,
		"https://amazon.com":    false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewURL("", false, map[string]string{
				"domain_re": `google\.(ru\.)?com$`,
			})
			assert.NoError(t, err)

			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect host")
			}
		})
	}
}

func (suite *ParameterURLTestSuite) TestValidateDomainPath() {
	testTable := map[string]bool{
		"https://amazon.de/path":   false,
		"https://google.ru.com/pp": true,
		"https://google.com/p":     true,
		"https://google.com/y":     false,
		"https://amazon.com":       false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewURL("", false, map[string]string{
				"domain_re": `google\.(ru\.)?com$`,
				"path_re":   `p+`,
			})
			assert.NoError(t, err)

			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect")
			}
		})
	}
}

func (suite *ParameterURLTestSuite) TestValidateUser() {
	testTable := map[string]bool{
		"https://amazon.de/path":       false,
		"https://aaa@amazon.de/path":   true,
		"https://aaa:c@amazon.de/path": true,
		"https://aa:c@amazon.de/path":  false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewURL("", false, map[string]string{
				"user_re": `aaa`,
			})
			assert.NoError(t, err)

			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect user")
			}
		})
	}
}

func (suite *ParameterURLTestSuite) TestCannotResolveHTTP() {
	suite.Dialer().
		On("DialContext", mock.Anything, "tcp", "amazon.com:80").
		Once().
		Return(suite.MakeNetConn(), io.EOF)

	param, err := config.NewURL("", false, map[string]string{
		"domain_re": `amazon\.com`,
		"resolve":   "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "http://amazon.com/ppp"),
		"cannot dial to")
}

func (suite *ParameterURLTestSuite) TestCannotResolveHTTPS() {
	suite.Dialer().
		On("DialContext", mock.Anything, "tcp", "amazon.com:443").
		Once().
		Return(suite.MakeNetConn(), io.EOF)

	param, err := config.NewURL("", false, map[string]string{
		"domain_re": `amazon\.com`,
		"resolve":   "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "https://amazon.com/ppp"),
		"cannot dial to")
}

func (suite *ParameterURLTestSuite) TestCannotResolveGivenPort() {
	suite.Dialer().
		On("DialContext", mock.Anything, "tcp", "amazon.com:4430").
		Once().
		Return(suite.MakeNetConn(), io.EOF)

	param, err := config.NewURL("", false, map[string]string{
		"domain_re": `amazon\.com`,
		"resolve":   "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "https://amazon.com:4430/ppp"),
		"cannot dial to")
}

func (suite *ParameterURLTestSuite) TestResolveOk() {
	connMock := suite.MakeNetConn()
	connMock.
		On("Close").
		Once().
		Return(nil)

	suite.Dialer().
		On("DialContext", mock.Anything, "tcp", "amazon.com:443").
		Once().
		Return(connMock, nil)

	param, err := config.NewURL("", false, map[string]string{
		"domain_re": `amazon\.com`,
		"resolve":   "true",
	})
	suite.NoError(err)
	suite.NoError(param.Validate(suite.Context(), "https://amazon.com/ppp"))
}

func TestParameterURL(t *testing.T) {
	suite.Run(t, &ParameterURLTestSuite{})
}
