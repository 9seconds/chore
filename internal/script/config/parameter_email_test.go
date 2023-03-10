package config_test

import (
	"io"
	"net"
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ParameterEmailTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
	testlib.NetworkTestSuite
}

func (suite *ParameterEmailTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
	suite.NetworkTestSuite.Setup(suite.T())
}

func (suite *ParameterEmailTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewEmail("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterEmailTestSuite) TestType() {
	param, err := config.NewEmail("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterEmail, param.Type())
}

func (suite *ParameterEmailTestSuite) TestIncorrectParameter() {
	testTable := map[string][]string{
		"resolve":   {"xx", ""},
		"domain_re": {"["},
		"name_re":   {"["},
	}

	for testName, testValues := range testTable {
		testName := testName
		testValues := testValues

		suite.T().Run(testName, func(t *testing.T) {
			for _, testValue := range testValues {
				testValue := testValue

				t.Run(testValue, func(t *testing.T) {
					_, err := config.NewEmail("", false, map[string]string{
						testName: testValue,
					})
					assert.Error(t, err)
				})
			}
		})
	}
}

func (suite *ParameterEmailTestSuite) TestValidateEmail() {
	testTable := map[string]bool{
		"@":                false,
		"":                 false,
		"xx":               false,
		"xx@":              false,
		"xx@yy":            true,
		"xx@yy-xx":         true,
		"ss.fff@gmail.com": true,
		".@g.x":            false,
		"Name <xx@yy>":     false,
		"<xx@yy>":          false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewEmail("", false, nil)
			assert.NoError(t, err)

			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect email")
			}
		})
	}
}

func (suite *ParameterEmailTestSuite) TestValidateEmailDomain() {
	testTable := map[string]bool{
		"xx@gmail.com":  true,
		"xx@gmail.ru":   false,
		"yy@yandex.com": false,
		"yy@gmail.com":  true,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewEmail("", false, map[string]string{
				"domain_re": `^gmail\.com$`,
			})
			assert.NoError(t, err)

			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "does not match")
			}
		})
	}
}

func (suite *ParameterEmailTestSuite) TestValidateEmailName() {
	testTable := map[string]bool{
		"xx@gmail.com":   false,
		"xx@gmail.ru":    false,
		"yy@yandex.com":  true,
		"yy@gmail.com":   true,
		"yyyy@gmail.com": false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewEmail("", false, map[string]string{
				"name_re": `^yy$`,
			})
			assert.NoError(t, err)

			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "does not match")
			}
		})
	}
}

func (suite *ParameterEmailTestSuite) TestValidateCannotResolve() {
	suite.DNS().
		On("LookupMX", mock.Anything, "gmail.com").
		Once().
		Return([]*net.MX{}, io.EOF)

	param, err := config.NewEmail("", false, map[string]string{
		"resolve": "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "xx@gmail.com"),
		"cannot resolve MX records of the domain")
}

func (suite *ParameterEmailTestSuite) TestValidateResolveNoMXRecords() {
	suite.DNS().
		On("LookupMX", mock.Anything, "gmail.com").
		Once().
		Return([]*net.MX{}, nil)

	param, err := config.NewEmail("", false, map[string]string{
		"resolve": "true",
	})
	suite.NoError(err)
	suite.ErrorContains(
		param.Validate(suite.Context(), "xx@gmail.com"),
		"no mx records were found")
}

func (suite *ParameterEmailTestSuite) TestValidateResolveOk() {
	suite.DNS().
		On("LookupMX", mock.Anything, "gmail.com").
		Once().
		Return([]*net.MX{nil}, nil)

	param, err := config.NewEmail("", false, map[string]string{
		"resolve": "true",
	})
	suite.NoError(err)
	suite.NoError(param.Validate(suite.Context(), "xx@gmail.com"))
}

func TestParameterEmail(t *testing.T) {
	suite.Run(t, &ParameterEmailTestSuite{})
}
