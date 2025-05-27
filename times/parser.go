package times

import (
	"errors"
)

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
	dt.date = date

	sep := b[10]
	if sep != 'T' && sep != ' ' && sep != 't' {
		return dt, nil, errors.New(string(b[10:11]) + "datetime separator is expected to be T or a space")
	}

	t, rest, err := parseLocalTime(b[11:])
	if err != nil {
		return dt, nil, err
	}
	dt.time = t

	return dt, rest, nil
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
