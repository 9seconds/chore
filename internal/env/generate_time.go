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
			StartedAtRFC3339, utc.Format(time.RFC3339))
		sendValue(ctx, results,
			StartedAtUnix, strconv.FormatInt(utc.Unix(), 10))
		sendValue(ctx, results,
			StartedAtYear, strconv.FormatInt(int64(utc.Year()), 10))
		sendValue(ctx, results,
			StartedAtYearDay, strconv.FormatInt(int64(utc.YearDay()), 10))
		sendValue(ctx, results,
			StartedAtDay, strconv.FormatInt(int64(utc.Day()), 10))
		sendValue(ctx, results,
			StartedAtMonth, strconv.FormatInt(int64(utc.Month()), 10))
		sendValue(ctx, results,
			StartedAtMonthStr, utc.Month().String())
		sendValue(ctx, results,
			StartedAtHour, strconv.FormatInt(int64(utc.Hour()), 10))
		sendValue(ctx, results,
			StartedAtMinute, strconv.FormatInt(int64(utc.Minute()), 10))
		sendValue(ctx, results,
			StartedAtSecond, strconv.FormatInt(int64(utc.Second()), 10))
		sendValue(ctx, results,
			StartedAtNanosecond, strconv.FormatInt(int64(utc.Nanosecond()), 10))
		sendValue(ctx, results,
			StartedAtOffset, strconv.FormatInt(int64(offset), 10))
		sendValue(ctx, results,
			StartedAtWeekday, strconv.FormatInt(int64(utc.Weekday()), 10))
		sendValue(ctx, results,
			StartedAtWeekdayStr, utc.Weekday().String())

		if tzname, err := tzlocal.RuntimeTZ(); err != nil {
			log.Printf("cannot obtain runtime timezone information: %v", err)
		} else {
			sendValue(ctx, results, StartedAtTimezone, tzname)
		}
	}()
}
