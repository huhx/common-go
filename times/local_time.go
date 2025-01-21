package times

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type LocalTime struct {
	Hour       int // Hour of the day: [0; 24[
	Minute     int // Minute of the hour: [0; 60[
	Second     int // Second of the minute: [0; 60[
	Nanosecond int // Nanoseconds within the second:  [0, 1000000000[
	Precision  int // Number of digits to display for Nanosecond.
}

func localTimeFromTime(t time.Time) LocalTime {
	return LocalTime{
		Hour:       t.Hour(),
		Minute:     t.Minute(),
		Second:     t.Second(),
		Nanosecond: t.Nanosecond(),
	}
}

func NewZeroTime() LocalTime {
	return LocalTime{
		Hour:       0,
		Minute:     0,
		Second:     0,
		Nanosecond: 0,
		Precision:  0,
	}
}

func (d LocalTime) String() string {
	s := fmt.Sprintf("%02d:%02d:%02d", d.Hour, d.Minute, d.Second)

	if d.Precision > 0 {
		s += fmt.Sprintf(".%09d", d.Nanosecond)[:d.Precision+1]
	} else if d.Nanosecond > 0 {
		// Nanoseconds are specified, but precision is not provided. Use the
		// minimum.
		s += strings.Trim(fmt.Sprintf(".%09d", d.Nanosecond), "0")
	}

	return s
}

// MarshalText returns RFC 3339 representation of d.
func (d LocalTime) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (d *LocalTime) UnmarshalText(b []byte) error {
	res, left, err := parseLocalTime(b)
	if err == nil && len(left) != 0 {
		err = errors.New(string(left) + "extra characters")
	}
	if err != nil {
		return err
	}
	*d = res
	return nil
}

func (lt LocalTime) Value() (driver.Value, error) {
	t := time.Date(0, 1, 1, lt.Hour, lt.Minute, lt.Second, 0, time.UTC)
	return t.Format("15:04:05"), nil
}

func (lt *LocalTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		lt.Hour = v.Hour()
		lt.Minute = v.Minute()
		lt.Second = v.Second()
		return nil
	case string:
		parsedTime, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		lt.Hour = parsedTime.Hour()
		lt.Minute = parsedTime.Minute()
		lt.Second = parsedTime.Second()
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into LocalTimeWrapper", value)
	}
}

// parseLocalTime is a bit different because it also returns the remaining
// []byte that is didn't need. This is to allow parseDateTime to parse those
// remaining bytes as a timezone.
func parseLocalTime(b []byte) (LocalTime, []byte, error) {
	var (
		nspow = [10]int{0, 1e8, 1e7, 1e6, 1e5, 1e4, 1e3, 1e2, 1e1, 1e0}
		t     LocalTime
	)

	// check if b matches to have expected format HH:MM:SS[.NNNNNN]
	const localTimeByteLen = 8
	if len(b) < localTimeByteLen {
		return t, nil, errors.New(string(b) + "times are expected to have the format HH:MM:SS[.NNNNNN]")
	}

	var err error

	t.Hour, err = parseDecimalDigits(b[0:2])
	if err != nil {
		return t, nil, err
	}

	if t.Hour > 23 {
		return t, nil, errors.New(string(b[0:2]) + "hour cannot be greater 23")
	}
	if b[2] != ':' {
		return t, nil, errors.New(string(b[2:3]) + "expecting colon between hours and minutes")
	}

	t.Minute, err = parseDecimalDigits(b[3:5])
	if err != nil {
		return t, nil, err
	}
	if t.Minute > 59 {
		return t, nil, errors.New(string(b[3:5]) + "minutes cannot be greater 59")
	}
	if b[5] != ':' {
		return t, nil, errors.New(string(b[5:6]) + "expecting colon between minutes and seconds")
	}

	t.Second, err = parseDecimalDigits(b[6:8])
	if err != nil {
		return t, nil, err
	}

	if t.Second > 60 {
		return t, nil, errors.New(string(b[6:8]) + "seconds cannot be greater 60")
	}

	b = b[8:]

	if len(b) >= 1 && b[0] == '.' {
		frac := 0
		precision := 0
		digits := 0

		for i, c := range b[1:] {
			if !isDigit(c) {
				if i == 0 {
					return t, nil, errors.New(string(b[0:1]) + "need at least one digit after fraction point")
				}
				break
			}
			digits++

			const maxFracPrecision = 9
			if i >= maxFracPrecision {
				// go-toml allows decoding fractional seconds
				// beyond the supported precision of 9
				// digits. It truncates the fractional component
				// to the supported precision and ignores the
				// remaining digits.
				//
				// https://github.com/pelletier/go-toml/discussions/707
				continue
			}

			frac *= 10
			frac += int(c - '0')
			precision++
		}

		if precision == 0 {
			return t, nil, errors.New(string(b[:1]) + "nanoseconds need at least one digit")
		}

		t.Nanosecond = frac * nspow[precision]
		t.Precision = precision

		return t, b[1+digits:], nil
	}
	return t, b, nil
}

func isDigit(r byte) bool {
	return r >= '0' && r <= '9'
}
