package config_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParameterDatetimeTestSuite struct {
	suite.Suite
	testlib.CtxTestSuite

	now time.Time
}

func (suite *ParameterDatetimeTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())

	suite.now = time.Now()
}

func (suite *ParameterDatetimeTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewDatetime(testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterDatetimeTestSuite) TestType() {
	param, err := config.NewDatetime(false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterDatetime, param.Type())
}

func (suite *ParameterDatetimeTestSuite) TestString() {
	param, err := config.NewDatetime(false, nil)
	suite.NoError(err)
	suite.NotEmpty(param.String())
}

func (suite *ParameterDatetimeTestSuite) TestIncorrectParameters() {
	for _, testName := range []string{"future_delta", "past_delta", "rounded_to"} {
		testName := testName

		suite.T().Run(testName, func(t *testing.T) {
			for _, choice := range []string{"-1", "1", "", "aa"} {
				choice := choice

				t.Run(choice, func(t *testing.T) {
					_, err := config.NewDatetime(false, map[string]string{
						testName: choice,
					})
					assert.ErrorContains(t, err, "incorrect duration")
				})
			}

			t.Run("-1s", func(t *testing.T) {
				_, err := config.NewDatetime(false, map[string]string{
					testName: "-1s",
				})
				assert.ErrorContains(t, err, "should be >=0")
			})
		})
	}
}

func (suite *ParameterDatetimeTestSuite) TestValidateUnixes() {
	testTable := []string{
		"unix",
		"unix_ms",
		"unix_us",
	}

	invalidChoices := []string{
		"-1",
		"",
		"x",
		"1a",
		"22_",
	}

	validChoices := []string{
		"0",
		"100",
	}

	for _, testName := range testTable {
		testName := testName

		suite.T().Run(testName, func(t *testing.T) {
			param, err := config.NewDatetime(false, map[string]string{
				"layout": testName,
			})
			assert.NoError(t, err)

			for _, choice := range invalidChoices {
				choice := choice

				t.Run(choice, func(t *testing.T) {
					assert.ErrorContains(
						t,
						param.Validate(suite.Context(), choice),
						"incorrect timestamp")
				})
			}

			for _, choice := range validChoices {
				choice := choice

				t.Run(choice, func(t *testing.T) {
					assert.NoError(
						t,
						param.Validate(suite.Context(), choice))
				})
			}
		})
	}
}

func (suite *ParameterDatetimeTestSuite) TestDelta() {
	unixNow := time.Now().Unix()
	stamps := map[string]int64{
		"past_before":   unixNow - 5,
		"past_after":    unixNow - 15,
		"future_before": unixNow + 5,
		"future_after":  unixNow + 15,
	}
	futures := map[string]map[string]bool{
		"past_delta": {
			"past_before":   true,
			"past_after":    false,
			"future_before": true,
			"future_after":  true,
		},
		"future_delta": {
			"past_before":   true,
			"past_after":    true,
			"future_before": true,
			"future_after":  false,
		},
	}

	for deltaName, deltaParams := range futures {
		deltaName := deltaName
		deltaParams := deltaParams

		suite.T().Run(deltaName, func(t *testing.T) {
			validator, err := config.NewDatetime(false, map[string]string{
				"layout":  "unix",
				deltaName: "10s",
			})
			assert.NoError(t, err)

			for paramName, isValid := range deltaParams {
				paramName := paramName
				isValid := isValid

				t.Run(paramName, func(t *testing.T) {
					err := validator.Validate(
						suite.Context(),
						strconv.FormatInt(stamps[paramName], 10))

					if isValid {
						assert.NoError(t, err)
					} else {
						assert.ErrorContains(t, err, "timestamp is too far")
					}
				})
			}
		})
	}
}

func (suite *ParameterDatetimeTestSuite) TestLayouts() {
	tme := time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)

	testTable := map[string]map[string]bool{
		"unix": {
			tme.Format(time.RFC3339):          false,
			tme.Format(time.RFC3339Nano):      false,
			tme.Format(time.RFC1123):          false,
			tme.Format(time.RFC1123Z):         false,
			tme.Format(time.RFC850):           false,
			tme.Format(time.RFC822):           false,
			tme.Format(time.RFC822Z):          false,
			tme.Format(time.RubyDate):         false,
			tme.Format(time.UnixDate):         false,
			strconv.FormatInt(tme.Unix(), 10): true,
		},
		"unix_ms": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMilli(), 10): true,
		},
		"unix_us": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): true,
		},
		"": {
			tme.Format(time.RFC3339):               true,
			tme.Format(time.RFC3339Nano):           true,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"rfc3339": {
			tme.Format(time.RFC3339):               true,
			tme.Format(time.RFC3339Nano):           true,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"rfc3339_nano": {
			tme.Format(time.RFC3339):               true,
			tme.Format(time.RFC3339Nano):           true,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"iso": {
			tme.Format(time.RFC3339):               true,
			tme.Format(time.RFC3339Nano):           true,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"iso8601": {
			tme.Format(time.RFC3339):               true,
			tme.Format(time.RFC3339Nano):           true,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"rfc1123": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               true,
			tme.Format(time.RFC1123Z):              true,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"rfc1123z": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              true,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"rfc850": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                true,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"rfc822": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                true,
			tme.Format(time.RFC822Z):               true,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"rfc822z": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               true,
			tme.Format(time.RubyDate):              false,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"ruby_date": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              true,
			tme.Format(time.UnixDate):              false,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"unix_date": {
			tme.Format(time.RFC3339):               false,
			tme.Format(time.RFC3339Nano):           false,
			tme.Format(time.RFC1123):               false,
			tme.Format(time.RFC1123Z):              false,
			tme.Format(time.RFC850):                false,
			tme.Format(time.RFC822):                false,
			tme.Format(time.RFC822Z):               false,
			tme.Format(time.RubyDate):              true,
			tme.Format(time.UnixDate):              true,
			strconv.FormatInt(tme.UnixMicro(), 10): false,
		},
		"2006": {
			"1992": true,
			"x":    false,
		},
		"u 5": {
			"u 6": true,
			"u":   false,
		},
	}

	for layoutName, tests := range testTable {
		layoutName := layoutName
		tests := tests

		suite.T().Run(layoutName, func(t *testing.T) {
			validator, err := config.NewDatetime(false, map[string]string{
				"layout": layoutName,
			})
			assert.NoError(t, err)

			for testName, isValid := range tests {
				testName := testName
				isValid := isValid

				t.Run(testName, func(t *testing.T) {
					err = validator.Validate(suite.Context(), testName)

					if isValid {
						assert.NoError(t, err)
					} else {
						assert.ErrorContains(t, err, "incorrect timestamp in layout")
					}
				})
			}
		})
	}

	for _, choice := range []string{"u", "u 2"} {
		choice := choice

		suite.T().Run(choice, func(t *testing.T) {
			_, err := config.NewDatetime(false, map[string]string{
				"layout": choice,
			})
			assert.ErrorContains(t, err, "incorrect layout")
		})
	}
}

func (suite *ParameterDatetimeTestSuite) TestRounding() {
	param, err := config.NewDatetime(false, map[string]string{
		"layout":     "3:04:05PM",
		"rounded_to": "5s",
	})
	suite.NoError(err)

	testTable := map[string]bool{
		"1:05:00AM": true,
		"1:05:02AM": false,
		"1:05:04AM": false,
		"1:05:05AM": true,
		"1:05:06AM": false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "timestamp is not rounded to")
			}
		})
	}
}

func (suite *ParameterDatetimeTestSuite) TestLocations() {
	tme := time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)

	param, err := config.NewDatetime(false, map[string]string{
		"layout":   time.RFC3339,
		"location": "UTC",
	})
	suite.NoError(err)

	berlinTimezone, err := time.LoadLocation("Europe/Berlin")
	require.NoError(suite.T(), err)

	moscowTimezone, err := time.LoadLocation("Europe/Moscow")
	require.NoError(suite.T(), err)

	testTable := map[string]bool{
		tme.Format(time.RFC3339):                    true,
		tme.In(berlinTimezone).Format(time.RFC3339): false,
		tme.In(moscowTimezone).Format(time.RFC3339): false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			err = param.Validate(suite.Context(), testName)

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "location is not matched")
			}
		})
	}
}

func TestParameterDatetime(t *testing.T) {
	suite.Run(t, &ParameterDatetimeTestSuite{})
}
