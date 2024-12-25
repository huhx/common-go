package times

import (
	"database/sql/driver"
	"errors"
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

func LocalDateNow() LocalDate {
	now := time.Now().In(timezone)
	return DateFromTime(now)
}

func (ldt LocalDate) PassDays(date LocalDate) int {
	return int(ldt.AsTime(timezone).Sub(date.AsTime(timezone)).Hours() / 24)
}

func (ldt LocalDate) ToSolar() LocalDate {
	lunar := calendar.NewLunarFromYmd(ldt.Year, ldt.Month, ldt.Day)
	solar := lunar.GetSolar()
	return DateFromYMD(solar.GetYear(), solar.GetMonth(), solar.GetDay())
}

// ToLunar todo：可以做优化，GetLunar方法里面计算了天干地支、节气之类的。我们这里不需要
func (ldt LocalDate) ToLunar() LocalDate {
	solar := calendar.NewSolarFromYmd(ldt.Year, ldt.Month, ldt.Day)
	lunar := solar.GetLunar()
	return DateFromYMD(lunar.GetYear(), lunar.GetMonth(), lunar.GetDay())
}

func (ldt LocalDate) PlusYear(year int) LocalDate {
	newTime := ldt.AsTime(timezone).AddDate(year, 0, 0)
	return DateFromTime(newTime)
}

func (ldt LocalDate) PlusMonth(month int) LocalDate {
	newTime := ldt.AsTime(timezone).AddDate(0, month, 0)
	return DateFromTime(newTime)
}

func (ldt LocalDate) PlusWeeks(weeks int) LocalDate {
	newTime := ldt.AsTime(timezone).AddDate(0, 0, 7*weeks)
	return DateFromTime(newTime)
}

func (ldt LocalDate) PlusDays(days int) LocalDate {
	newTime := ldt.AsTime(timezone).AddDate(0, 0, days)
	return DateFromTime(newTime)
}

func (ldt LocalDate) Compare(date LocalDate) int {
	return ldt.AsTime(timezone).Compare(date.AsTime(timezone))
}

func (ldt LocalDate) Before(date LocalDate) bool {
	return ldt.Compare(date) < 0
}

func (ldt LocalDate) After(date LocalDate) bool {
	return ldt.Compare(date) > 0
}

func (ldt LocalDate) Equal(date LocalDate) bool {
	return ldt.Compare(date) == 0
}

func (ldt LocalDate) CopyYear(year int) LocalDate {
	return DateFromYMD(year, ldt.Month, ldt.Day)
}

// AsTime converts d into a specific times instance at midnight in zone.
func (ldt LocalDate) AsTime(zone *time.Location) time.Time {
	return time.Date(ldt.Year, time.Month(ldt.Month), ldt.Day, 0, 0, 0, 0, zone)
}

// String returns RFC 3339 representation of d.
func (ldt LocalDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", ldt.Year, ldt.Month, ldt.Day)
}

// MarshalText returns RFC 3339 representation of d.
func (ldt LocalDate) MarshalText() ([]byte, error) {
	return []byte(ldt.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (ldt *LocalDate) UnmarshalText(b []byte) error {
	res, err := parseLocalDate(b)
	if err != nil {
		return err
	}
	*ldt = res
	return nil
}

func (ldt LocalDate) Value() (driver.Value, error) {
	t := time.Date(ldt.Year, time.Month(ldt.Month), ldt.Day, 0, 0, 0, 0, time.UTC)
	return t, nil
}

func (ldt *LocalDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ldt.Year = v.Year()
		ldt.Month = int(v.Month())
		ldt.Day = v.Day()
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into LocalDateWrapper", value)
	}
}

func parseLocalDate(b []byte) (LocalDate, error) {
	// full-date      = date-fullyear "-" date-month "-" date-mday
	// date-fullyear  = 4DIGIT
	// date-month     = 2DIGIT  ; 01-12
	// date-mday      = 2DIGIT  ; 01-28, 01-29, 01-30, 01-31 based on month/year
	var date LocalDate

	if len(b) != 10 || b[4] != '-' || b[7] != '-' {
		return date, errors.New(string(b) + "dates are expected to have the format YYYY-MM-DD")
	}

	var err error

	date.Year, err = parseDecimalDigits(b[0:4])
	if err != nil {
		return LocalDate{}, err
	}

	date.Month, err = parseDecimalDigits(b[5:7])
	if err != nil {
		return LocalDate{}, err
	}

	date.Day, err = parseDecimalDigits(b[8:10])
	if err != nil {
		return LocalDate{}, err
	}

	if !isValidDate(date.Year, date.Month, date.Day) {
		return LocalDate{}, errors.New(string(b) + "impossible date")
	}

	return date, nil
}

func isValidDate(year int, month int, day int) bool {
	return month > 0 && month < 13 && day > 0 && day <= daysIn(month, year)
}

var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

func daysIn(m int, year int) int {
	if m == 2 && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func parseDecimalDigits(b []byte) (int, error) {
	v := 0

	for i, c := range b {
		if c < '0' || c > '9' {
			return 0, errors.New(string(b[i:i+1]) + "expected digit (0-9)")
		}
		v *= 10
		v += int(c - '0')
	}

	return v, nil
}
