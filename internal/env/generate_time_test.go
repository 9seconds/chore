package env_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateTimeTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateTimeTestSuite) TestOk() {
	env.GenerateTime(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()
	suite.Len(data, 15)

	tme, err := time.Parse(time.RFC3339, data[env.EnvStartedAtRFC3339])
	suite.NoError(err)

	suite.WithinDuration(time.Now(), tme, time.Second)
	suite.Equal(
		strconv.FormatInt(tme.Unix(), 10),
		data[env.EnvStartedAtUnix])
	suite.Equal(
		strconv.FormatInt(int64(tme.Year()), 10),
		data[env.EnvStartedAtYear])
	suite.Equal(
		strconv.FormatInt(int64(tme.YearDay()), 10),
		data[env.EnvStartedAtYearDay])
	suite.Equal(
		strconv.FormatInt(int64(tme.Day()), 10),
		data[env.EnvStartedAtDay])
	suite.Equal(
		strconv.FormatInt(int64(tme.Month()), 10),
		data[env.EnvStartedAtMonth])
	suite.Equal(
		tme.Month().String(),
		data[env.EnvStartedAtMonthStr])
	suite.Equal(
		strconv.FormatInt(int64(tme.Hour()), 10),
		data[env.EnvStartedAtHour])
	suite.Equal(
		strconv.FormatInt(int64(tme.Minute()), 10),
		data[env.EnvStartedAtMinute])
	suite.Equal(
		strconv.FormatInt(int64(tme.Second()), 10),
		data[env.EnvStartedAtSecond])
	suite.NotEmpty(data[env.EnvStartedAtNanosecond])
	suite.Equal(
		strconv.FormatInt(int64(tme.Weekday()), 10),
		data[env.EnvStartedAtWeekday])
	suite.Equal(
		tme.Weekday().String(),
		data[env.EnvStartedAtWeekdayStr])
	suite.NotEmpty(data[env.EnvStartedAtTimezone])
	suite.NotEmpty(data[env.EnvStartedAtOffset])
}

func TestGenerateTime(t *testing.T) {
	suite.Run(t, &GenerateTimeTestSuite{})
}
