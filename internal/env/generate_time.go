package env

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/thlib/go-timezone-local/tzlocal"
)

func GenerateTime(ctx context.Context, results chan<- string, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		now := time.Now()
		_, offset := now.Zone()
		utc := now.UTC()

		sendValue(ctx, results,
			EnvStartedAtRFC3339, utc.Format(time.RFC3339))
		sendValue(ctx, results,
			EnvStartedAtUnix, strconv.FormatInt(utc.Unix(), 10))
		sendValue(ctx, results,
			EnvStartedAtYear, strconv.FormatInt(int64(utc.Year()), 10))
		sendValue(ctx, results,
			EnvStartedAtYearDay, strconv.FormatInt(int64(utc.YearDay()), 10))
		sendValue(ctx, results,
			EnvStartedAtDay, strconv.FormatInt(int64(utc.Day()), 10))
		sendValue(ctx, results,
			EnvStartedAtMonth, strconv.FormatInt(int64(utc.Month()), 10))
		sendValue(ctx, results,
			EnvStartedAtMonthStr, utc.Month().String())
		sendValue(ctx, results,
			EnvStartedAtHour, strconv.FormatInt(int64(utc.Hour()), 10))
		sendValue(ctx, results,
			EnvStartedAtMinute, strconv.FormatInt(int64(utc.Minute()), 10))
		sendValue(ctx, results,
			EnvStartedAtSecond, strconv.FormatInt(int64(utc.Second()), 10))
		sendValue(ctx, results,
			EnvStartedAtNanosecond, strconv.FormatInt(int64(utc.Nanosecond()), 10))
		sendValue(ctx, results,
			EnvStartedAtOffset, strconv.FormatInt(int64(offset), 10))
		sendValue(ctx, results,
			EnvStartedAtWeekday, strconv.FormatInt(int64(utc.Weekday()), 10))
		sendValue(ctx, results,
			EnvStartedAtWeekdayStr, utc.Weekday().String())

		if tzname, err := tzlocal.RuntimeTZ(); err != nil {
			log.Printf("cannot obtain runtime timezone information: %v", err)
		} else {
			sendValue(ctx, results, EnvStartedAtTimezone, tzname)
		}
	}()
}
