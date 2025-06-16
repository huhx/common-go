package times

import (
	"database/sql/driver"
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"time"
)

type LocalDate struct {
	Year  int
	Month int
	Day   int
}

func DateFromYMD(year, month, day int) LocalDate {
	return LocalDate{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

func DateFromTime(t time.Time) LocalDate {
	return LocalDate{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}
}

func DateFromString(dateString string, pattern string) (*LocalDate, error) {
	if date, err := time.Parse(pattern, dateString); err != nil {
		return nil, err
	} else {
		fromTime := DateFromTime(date)
		return &fromTime, nil
	}
}

func DateFromDefaultPattern(dateString string) (*LocalDate, error) {
	if date, err := time.Parse(time.DateOnly, dateString); err != nil {
		return nil, err
	} else {
		fromTime := DateFromTime(date)
		return &fromTime, nil
	}
}

func (ld LocalDate) Weekday() int {
	return int(ld.AsTime(timezone).Weekday())
}

func LocalDateNow() LocalDate {
	now := time.Now().In(timezone)
	return DateFromTime(now)
}

func (ld LocalDate) Age() int {
	now := LocalDateNow()
	age := now.Year - ld.Year
	if now.Month < ld.Month || (now.Month == ld.Month && now.Day < ld.Day) {
		age--
	}
	return age
}

func (ld LocalDate) StartOfMonth() LocalDate {
	return DateFromYMD(ld.Year, ld.Month, 1)
}

func (ld LocalDate) EndOfMonth() LocalDate {
	return DateFromYMD(ld.Year, ld.Month, daysIn(ld.Month, ld.Year))
}

func (ld LocalDate) PassDays(date LocalDate) int {
	return int(ld.AsTime(timezone).Sub(date.AsTime(timezone)) / (24 * time.Hour))
}

func (ld LocalDate) ToSolar() LocalDate {
	lunar := calendar.NewLunarFromYmd(ld.Year, ld.Month, ld.Day)
	solar := lunar.GetSolar()
	return DateFromYMD(solar.GetYear(), solar.GetMonth(), solar.GetDay())
}

// ToLunar todo：可以做优化，GetLunar方法里面计算了天干地支、节气之类的。我们这里不需要
func (ld LocalDate) ToLunar() LocalDate {
	solar := calendar.NewSolarFromYmd(ld.Year, ld.Month, ld.Day)
	lunar := solar.GetLunar()
	return DateFromYMD(lunar.GetYear(), lunar.GetMonth(), lunar.GetDay())
}

func (ld LocalDate) PlusYear(year int) LocalDate {
	newTime := ld.AsTime(timezone).AddDate(year, 0, 0)
	return DateFromTime(newTime)
}

func (ld LocalDate) PlusMonth(month int) LocalDate {
	newTime := ld.AsTime(timezone).AddDate(0, month, 0)
	return DateFromTime(newTime)
}

func (ld LocalDate) PlusWeeks(weeks int) LocalDate {
	newTime := ld.AsTime(timezone).AddDate(0, 0, 7*weeks)
	return DateFromTime(newTime)
}

func (ld LocalDate) PlusDays(days int) LocalDate {
	newTime := ld.AsTime(timezone).AddDate(0, 0, days)
	return DateFromTime(newTime)
}

func (ld LocalDate) Compare(date LocalDate) int {
	return ld.AsTime(timezone).Compare(date.AsTime(timezone))
}

func (ld LocalDate) Before(date LocalDate) bool {
	return ld.Compare(date) < 0
}

func (ld LocalDate) After(date LocalDate) bool {
	return ld.Compare(date) > 0
}

func (ld LocalDate) Equal(date LocalDate) bool {
	return ld.Compare(date) == 0
}

func (ld LocalDate) CopyYear(year int) LocalDate {
	return DateFromYMD(year, ld.Month, ld.Day)
}

// AsTime converts d into a specific times instance at midnight in zone.
func (ld LocalDate) AsTime(zone *time.Location) time.Time {
	return time.Date(ld.Year, time.Month(ld.Month), ld.Day, 0, 0, 0, 0, zone)
}

// String returns RFC 3339 representation of d.
func (ld LocalDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", ld.Year, ld.Month, ld.Day)
}

// MarshalText returns RFC 3339 representation of d.
func (ld LocalDate) MarshalText() ([]byte, error) {
	return []byte(ld.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (ld *LocalDate) UnmarshalText(b []byte) error {
	res, err := parseLocalDate(b)
	if err != nil {
		return err
	}
	*ld = res
	return nil
}

func (ld LocalDate) Value() (driver.Value, error) {
	t := time.Date(ld.Year, time.Month(ld.Month), ld.Day, 0, 0, 0, 0, time.UTC)
	return t, nil
}

func (ld *LocalDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ld.Year = v.Year()
		ld.Month = int(v.Month())
		ld.Day = v.Day()
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into LocalDateWrapper", value)
	}
}
