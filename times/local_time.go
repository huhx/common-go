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

func TimeFromString(timeString string, pattern string) (*LocalTime, error) {
	if datetime, err := time.Parse(pattern, timeString); err != nil {
		return nil, err
	} else {
		fromTime := localTimeFromTime(datetime)
		return &fromTime, nil
	}
}

func TimeFromDefaultString(timeString string) (*LocalTime, error) {
	if datetime, err := time.Parse(time.TimeOnly, timeString); err != nil {
		return nil, err
	} else {
		fromTime := localTimeFromTime(datetime)
		return &fromTime, nil
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

func (lt LocalTime) String() string {
	s := fmt.Sprintf("%02d:%02d:%02d", lt.Hour, lt.Minute, lt.Second)

	if lt.Precision > 0 {
		s += fmt.Sprintf(".%09d", lt.Nanosecond)[:lt.Precision+1]
	} else if lt.Nanosecond > 0 {
		// Nanoseconds are specified, but precision is not provided. Use the
		// minimum.
		s += strings.Trim(fmt.Sprintf(".%09d", lt.Nanosecond), "0")
	}

	return s
}

// MarshalText returns RFC 3339 representation of d.
func (lt LocalTime) MarshalText() ([]byte, error) {
	return []byte(lt.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (lt *LocalTime) UnmarshalText(b []byte) error {
	res, left, err := parseLocalTime(b)
	if err == nil && len(left) != 0 {
		err = errors.New(string(left) + "extra characters")
	}
	if err != nil {
		return err
	}
	*lt = res
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
