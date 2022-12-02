package config

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const ParameterDatetime = "datetime"

type parameterDatetime struct {
	required    bool
	futureDelta time.Duration
	pastDelta   time.Duration
	roundedTo   time.Duration
	location    *time.Location
	layout      string
}

func (p parameterDatetime) Required() bool {
	return p.required
}

func (p parameterDatetime) Type() string {
	return ParameterDatetime
}

func (p parameterDatetime) String() string {
	return fmt.Sprintf(
		"required=%t, future_delta=%s, past_delta=%s, rounded_to=%s, location=%s, layout=%s",
		p.required,
		p.futureDelta,
		p.pastDelta,
		p.roundedTo,
		p.location,
		p.layout)
}

func (p parameterDatetime) Validate(_ context.Context, value string) error { //nolint: cyclop
	var (
		sec uint64
		tme time.Time
		err error
	)

	switch p.layout {
	case "unix", "unix_ms", "unix_us":
		sec, err = strconv.ParseUint(value, 10, 64)

		switch p.layout {
		case "unix":
			tme = time.Unix(int64(sec), 0)
		case "unix_ms":
			tme = time.UnixMilli(int64(sec))
		default:
			tme = time.UnixMicro(int64(sec))
		}
	default:
		tme, err = time.Parse(p.layout, value)
	}

	delta := time.Since(tme)

	switch {
	case err != nil:
		return fmt.Errorf("incorrect timestamp in layout %s: %w", p.layout, err)
	case p.futureDelta >= 0 && delta < -p.futureDelta:
		return fmt.Errorf("timestamp is too far in future: delta=%s, expected=%s", delta, p.futureDelta)
	case p.pastDelta >= 0 && delta > p.pastDelta:
		return fmt.Errorf("timestamp is too far in past: delta=%s, expected=%s", delta, p.pastDelta)
	case p.roundedTo > 0 && !tme.Equal(tme.Round(p.roundedTo)):
		return fmt.Errorf("timestamp is not rounded to %s", p.roundedTo)
	case p.location != nil:
		expectedZone, expectedOffset := tme.In(p.location).Zone()
		actualZone, actualOffset := tme.Zone()

		if expectedOffset != actualOffset {
			return fmt.Errorf(
				"location is not matched: expected=%s(%d), got=%s(%d)",
				expectedZone,
				expectedOffset,
				actualZone,
				actualOffset)
		}
	}

	return nil
}

func NewDatetime(required bool, spec map[string]string) (Parameter, error) { //nolint: cyclop
	param := parameterDatetime{
		required:    required,
		futureDelta: -1,
		pastDelta:   -1,
	}

	if value, err := parseDurationNegative(spec, "future_delta"); err == nil {
		param.futureDelta = value
	} else {
		return nil, err
	}

	if value, err := parseDurationNegative(spec, "past_delta"); err == nil {
		param.pastDelta = value
	} else {
		return nil, err
	}

	if value, err := parseDurationNegative(spec, "rounded_to"); err == nil {
		param.roundedTo = value
	} else {
		return nil, err
	}

	if value, ok := spec["location"]; ok {
		parsed, err := time.LoadLocation(value)
		if err != nil {
			return nil, fmt.Errorf("cannot load location: %w", err)
		}

		param.location = parsed
	}

	switch strings.ToLower(spec["layout"]) {
	case "unix", "unix_ms", "unix_us":
		param.layout = spec["layout"]
	case "", "rfc3339", "iso", "iso8601":
		param.layout = time.RFC3339
	case "rfc3339_nano":
		param.layout = time.RFC3339Nano
	case "rfc1123":
		param.layout = time.RFC1123
	case "rfc1123z":
		param.layout = time.RFC1123Z
	case "rfc850":
		param.layout = time.RFC850
	case "rfc822":
		param.layout = time.RFC822
	case "rfc822z":
		param.layout = time.RFC822Z
	case "ruby_date":
		param.layout = time.RubyDate
	case "unix_date":
		param.layout = time.UnixDate
	default:
		if time.Now().Format(spec["layout"]) == spec["layout"] {
			return nil, fmt.Errorf("incorrect layout: %s", spec["layout"])
		}

		param.layout = spec["layout"]
	}

	return param, nil
}
