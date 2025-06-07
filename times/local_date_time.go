package times

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

var timezone = time.Local

type LocalDateTime struct {
	date LocalDate
	time LocalTime
}

func LocalDateTimeNow() LocalDateTime {
	now := time.Now().In(timezone)
	return DateTimeFromTime(now)
}

func (ldt LocalDateTime) Date() LocalDate {
	return ldt.date
}

func (ldt LocalDateTime) Time() LocalTime {
	return ldt.time
}

func DateTimeFromTime(datetime time.Time) LocalDateTime {
	return LocalDateTime{DateFromTime(datetime), localTimeFromTime(datetime)}
}

func DateTimeFromString(datetimeString string, pattern string) (*LocalDateTime, error) {
	if datetime, err := time.Parse(pattern, datetimeString); err != nil {
		return nil, err
	} else {
		fromTime := DateTimeFromTime(datetime)
		return &fromTime, nil
	}
}

func DateTimeFromDefaultString(datetimeString string) (*LocalDateTime, error) {
	if datetime, err := time.Parse(time.DateTime, datetimeString); err != nil {
		return nil, err
	} else {
		fromTime := DateTimeFromTime(datetime)
		return &fromTime, nil
	}
}

func (ldt LocalDateTime) StartOfToday() LocalDateTime {
	return LocalDateTime{ldt.date, NewZeroTime()}
}

func (ldt LocalDateTime) PassDays(dateTime LocalDateTime) int {
	return int(ldt.AsTime(timezone).Sub(dateTime.AsTime(timezone)).Hours() / 24)
}

func (ldt LocalDateTime) PassHours(dateTime LocalDateTime) int {
	return int(ldt.AsTime(timezone).Sub(dateTime.AsTime(timezone)).Hours())
}

func (ldt LocalDateTime) PassMinutes(dateTime LocalDateTime) int {
	return int(ldt.AsTime(timezone).Sub(dateTime.AsTime(timezone)).Minutes())
}

func (ldt LocalDateTime) PassSeconds(dateTime LocalDateTime) int {
	return int(ldt.AsTime(timezone).Sub(dateTime.AsTime(timezone)).Seconds())
}

func (ldt LocalDateTime) ToSolar() LocalDateTime {
	localDate := ldt.date.ToSolar()
	return LocalDateTime{localDate, ldt.time}
}

func (ldt LocalDateTime) ToLunar() LocalDateTime {
	localDate := ldt.date.ToLunar()
	return LocalDateTime{localDate, ldt.time}
}

func (ldt LocalDateTime) PlusYears(year int) LocalDateTime {
	localDate := ldt.date.PlusYear(year)
	return LocalDateTime{localDate, ldt.time}
}

func (ldt LocalDateTime) PlusMonths(month int) LocalDateTime {
	localDate := ldt.date.PlusMonth(month)
	return LocalDateTime{localDate, ldt.time}
}

func (ldt LocalDateTime) PlusWeeks(weeks int) LocalDateTime {
	localDate := ldt.date.PlusWeeks(weeks)
	return LocalDateTime{localDate, ldt.time}
}

func (ldt LocalDateTime) PlusDays(days int) LocalDateTime {
	localDate := ldt.date.PlusDays(days)
	return LocalDateTime{localDate, ldt.time}
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

func (ldt LocalDateTime) AsTime(zone *time.Location) time.Time {
	return time.Date(
		ldt.date.Year, time.Month(ldt.date.Month), ldt.date.Day, ldt.time.Hour,
		ldt.time.Minute, ldt.time.Second, ldt.time.Nanosecond, zone,
	)
}

func (ldt LocalDateTime) String() string {
	return ldt.date.String() + "T" + ldt.time.String()
}

func (ldt LocalDateTime) MarshalText() ([]byte, error) {
	return []byte(ldt.String()), nil
}

func (ldt *LocalDateTime) UnmarshalText(data []byte) error {
	res, left, err := parseLocalDateTime(data)
	if err == nil && len(left) != 0 {
		err = errors.New(string(left) + "extra characters")
	}
	if err != nil {
		return err
	}

	*ldt = res
	return nil
}

func (ldt LocalDateTime) Value() (driver.Value, error) {
	return ldt.AsTime(time.UTC), nil
}

func (ldt *LocalDateTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ldt.date = DateFromYMD(v.Year(), int(v.Month()), v.Day())
		ldt.time = localTimeFromTime(v)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into LocalDateTimeWrapper", value)
	}
}
