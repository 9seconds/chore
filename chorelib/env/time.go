package env

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/thlib/go-timezone-local/tzlocal"
)

func generateTime(ctx context.Context, values chan<- string) {
	now := time.Now()
	_, offset := now.Zone()
	utc := now.UTC()

	sendEnvValue(ctx, values,
		EnvStartedAtRFC3339, utc.Format(time.RFC3339))
	sendEnvValue(ctx, values,
		EnvStartedAtUnix, strconv.FormatInt(utc.Unix(), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtYear, strconv.FormatInt(int64(utc.Year()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtYearDay, strconv.FormatInt(int64(utc.YearDay()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtDay, strconv.FormatInt(int64(utc.Day()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtMonth, strconv.FormatInt(int64(utc.Month()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtMonthStr, utc.Month().String())
	sendEnvValue(ctx, values,
		EnvStartedAtHour, strconv.FormatInt(int64(utc.Hour()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtMinute, strconv.FormatInt(int64(utc.Minute()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtSecond, strconv.FormatInt(int64(utc.Second()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtNanosecond, strconv.FormatInt(int64(utc.Nanosecond()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtOffset, strconv.FormatInt(int64(offset), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtWeekday, strconv.FormatInt(int64(utc.Weekday()), 10))
	sendEnvValue(ctx, values,
		EnvStartedAtWeekdayStr, utc.Weekday().String())

	if tzname, err := tzlocal.RuntimeTZ(); err != nil {
		log.Printf("cannot obtain runtime timezone information: %v", err)
	} else {
		sendEnvValue(ctx, values, EnvStartedAtTimezone, tzname)
	}
}
