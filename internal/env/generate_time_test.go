package env_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateTimeTestSuite struct {
	BaseTestSuite
}

func (suite *GenerateTimeTestSuite) TestOk() {
	env.GenerateTime(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()
	suite.Len(data, 15)

	tme, err := time.Parse(time.RFC3339, data[env.StartedAtRFC3339])
	suite.NoError(err)

	suite.WithinDuration(time.Now(), tme, time.Second)
	suite.Equal(
		strconv.FormatInt(tme.Unix(), 10),
		data[env.StartedAtUnix])
	suite.Equal(
		strconv.FormatInt(int64(tme.Year()), 10),
		data[env.StartedAtYear])
	suite.Equal(
		strconv.FormatInt(int64(tme.YearDay()), 10),
		data[env.StartedAtYearDay])
	suite.Equal(
		strconv.FormatInt(int64(tme.Day()), 10),
		data[env.StartedAtDay])
	suite.Equal(
		strconv.FormatInt(int64(tme.Month()), 10),
		data[env.StartedAtMonth])
	suite.Equal(
		tme.Month().String(),
		data[env.StartedAtMonthStr])
	suite.Equal(
		strconv.FormatInt(int64(tme.Hour()), 10),
		data[env.StartedAtHour])
	suite.Equal(
		strconv.FormatInt(int64(tme.Minute()), 10),
		data[env.StartedAtMinute])
	suite.Equal(
		strconv.FormatInt(int64(tme.Second()), 10),
		data[env.StartedAtSecond])
	suite.NotEmpty(data[env.StartedAtNanosecond])
	suite.Equal(
		strconv.FormatInt(int64(tme.Weekday()), 10),
		data[env.StartedAtWeekday])
	suite.Equal(
		tme.Weekday().String(),
		data[env.StartedAtWeekdayStr])
	suite.NotEmpty(data[env.StartedAtTimezone])
	suite.NotEmpty(data[env.StartedAtOffset])
}

func TestGenerateTime(t *testing.T) {
	suite.Run(t, &GenerateTimeTestSuite{})
}
