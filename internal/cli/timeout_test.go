package cli_test

import (
	"testing"
	"time"

	"github.com/9seconds/chore/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TimeoutTestSuite struct {
	suite.Suite

	param cli.Timeout
}

func (suite *TimeoutTestSuite) TestParseNegative() {
	suite.ErrorContains(
		suite.param.UnmarshalText([]byte("-10s")),
		"duration should be positive")
}

func (suite *TimeoutTestSuite) TestParseFail() {
	testTable := []string{
		"",
		"xxx",
		"1i9s0dfzx",
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.ErrorContains(
				t,
				suite.param.UnmarshalText([]byte(testValue)),
				"cannot parse duration")
		})
	}
}

func (suite *TimeoutTestSuite) TestParse() {
	testTable := map[string]time.Duration{
		"10":    10 * time.Second,
		"100ms": 100 * time.Millisecond,
		"1h":    time.Hour,
		"1h5m":  time.Hour + 5*time.Minute,
	}

	for testName, expectedValue := range testTable {
		testName := testName
		expectedValue := expectedValue

		suite.T().Run(testName, func(t *testing.T) {
			var param cli.Timeout

			assert.NoError(t, param.UnmarshalText([]byte(testName)))
			assert.Equal(t, expectedValue, param.Value())
		})
	}
}

func TestTimeout(t *testing.T) {
	suite.Run(t, &TimeoutTestSuite{})
}
