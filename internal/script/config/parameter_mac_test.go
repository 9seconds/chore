package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterMacTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterMacTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterMacTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewMac("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterMacTestSuite) TestType() {
	param, err := config.NewMac("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterMac, param.Type())
}

func (suite *ParameterMacTestSuite) TestValidaton() {
	testTable := map[string]bool{
		"X":   false,
		"XX":  false,
		"1":   false,
		"AB:": false,
		"":    false,

		"00:00:5e:00:53:01":       true,
		"02:00:5e:10:00:00:00:01": true,
		"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01": true,
		"00-00-5e-00-53-01":       true,
		"02-00-5e-10-00-00-00-01": true,
		"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01": true,
		"0000.5e00.5301":      true,
		"0200.5e10.0000.0001": true,
		"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001": true,
	}

	param, err := config.NewMac("", false, nil)
	suite.NoError(err)

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err := param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect mac address")
			}
		})
	}
}

func TestParameterMac(t *testing.T) {
	suite.Run(t, &ParameterMacTestSuite{})
}
