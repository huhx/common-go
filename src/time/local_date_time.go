package local

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

var timezone = time.Local

type LocalDateTime struct {
	LocalDate
	LocalTime
}

// AsTime converts d into a specific time instance in zone.
func (d LocalDateTime) AsTime(zone *time.Location) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, d.Hour, d.Minute, d.Second, d.Nanosecond, zone)
}

// String returns RFC 3339 representation of d.
func (d LocalDateTime) String() string {
	return d.LocalDate.String() + "T" + d.LocalTime.String()
}

// MarshalText returns RFC 3339 representation of d.
func (d LocalDateTime) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (d *LocalDateTime) UnmarshalText(data []byte) error {
	res, left, err := parseLocalDateTime(data)
	if err == nil && len(left) != 0 {
		err = errors.New(string(left) + "extra characters")
	}
	if err != nil {
		return err
	}

	*d = res
	return nil
}

func parseLocalDateTime(b []byte) (LocalDateTime, []byte, error) {
	var dt LocalDateTime

	const localDateTimeByteMinLen = 11
	if len(b) < localDateTimeByteMinLen {
		return dt, nil, errors.New(string(b) + "local datetimes are expected to have the format YYYY-MM-DDTHH:MM:SS[.NNNNNNNNN]")
	}

	date, err := parseLocalDate(b[:10])
	if err != nil {
		return dt, nil, err
	}
	dt.LocalDate = date

	sep := b[10]
	if sep != 'T' && sep != ' ' && sep != 't' {
		return dt, nil, errors.New(string(b[10:11]) + "datetime separator is expected to be T or a space")
	}

	t, rest, err := parseLocalTime(b[11:])
	if err != nil {
		return dt, nil, err
	}
	dt.LocalTime = t

	return dt, rest, nil
}

func (ldt LocalDateTime) Value() (driver.Value, error) {
	t := time.Date(ldt.Year, time.Month(ldt.Month), ldt.Day, ldt.Hour, ldt.Minute, ldt.Second, 0, time.UTC)
	return t, nil
}

func (ldt *LocalDateTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ldt.Year = v.Year()
		ldt.Month = int(v.Month())
		ldt.Day = v.Day()
		ldt.Hour = v.Hour()
		ldt.Minute = v.Minute()
		ldt.Second = v.Second()
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into LocalDateTimeWrapper", value)
	}
}

func LocalDateTimeNow() LocalDateTime {
	now := time.Now().In(timezone)
	localTime := LocalTime{
		Hour:       now.Hour(),
		Minute:     now.Minute(),
		Second:     now.Second(),
		Nanosecond: now.Nanosecond(),
	}
	localDate := LocalDate{
		Year:  now.Year(),
		Month: int(now.Month()),
		Day:   now.Day(),
	}
	return LocalDateTime{localDate, localTime}
}

func (ldt LocalDateTime) PassDays(dateTime LocalDateTime) int {
	return int(ldt.AsTime(timezone).Sub(dateTime.AsTime(timezone)).Hours() / 24)
}

func (ldt LocalDateTime) ToSolar() LocalDateTime {
	localDate := ldt.LocalDate.ToSolar()
	return LocalDateTime{localDate, ldt.LocalTime}
}

func (ldt LocalDateTime) ToLunar() LocalDateTime {
	localDate := ldt.LocalDate.ToLunar()
	return LocalDateTime{localDate, ldt.LocalTime}
}

func (ldt LocalDateTime) PlusYears(year int) LocalDateTime {
	localDate := ldt.LocalDate.PlusYear(year)
	return LocalDateTime{localDate, ldt.LocalTime}
}

func (ldt LocalDateTime) PlusMonths(month int) LocalDateTime {
	localDate := ldt.LocalDate.PlusMonth(month)
	return LocalDateTime{localDate, ldt.LocalTime}
}

func (ldt LocalDateTime) PlusWeeks(weeks int) LocalDateTime {
	localDate := ldt.LocalDate.PlusWeeks(weeks)
	return LocalDateTime{localDate, ldt.LocalTime}
}

func (ldt LocalDateTime) PlusDays(days int) LocalDateTime {
	localDate := ldt.LocalDate.PlusDays(days)
	return LocalDateTime{localDate, ldt.LocalTime}
}

func (ldt LocalDateTime) Compare(dateTim LocalDateTime) int {
	return ldt.AsTime(timezone).Compare(dateTim.AsTime(timezone))
}

func (ldt LocalDateTime) Before(dateTim LocalDateTime) bool {
	return ldt.Compare(dateTim) < 0
}

func (ldt LocalDateTime) After(dateTim LocalDateTime) bool {
	return ldt.Compare(dateTim) > 0
}

func (ldt LocalDateTime) Equal(dateTim LocalDateTime) bool {
	return ldt.Compare(dateTim) == 0
}
