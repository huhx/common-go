package times

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDateFromString(t *testing.T) {
	type args struct {
		dateString string
		pattern    string
	}
	tests := []struct {
		name     string
		args     args
		want     *LocalDate
		hasError bool
	}{
		{
			name:     "test the pattern is match the date string",
			args:     args{dateString: "2024-12-12", pattern: "2006-01-02"},
			want:     &LocalDate{Year: 2024, Month: 12, Day: 12},
			hasError: false,
		},
		{
			name:     "test the pattern is not matched the date string",
			args:     args{dateString: "2024-12-12", pattern: "20060102"},
			want:     nil,
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DateFromString(tt.args.dateString, tt.args.pattern)

			assert.Equalf(t, tt.hasError, err != nil, "DateFromString(%v, %v)", tt.args.dateString, tt.args.pattern)
			assert.Equalf(t, tt.want, result, "DateFromString(%v, %v)", tt.args.dateString, tt.args.pattern)
		})
	}
}

func TestLocalDate_ToSolar(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   LocalDate
	}{
		{
			name:   "test the to solar conversion",
			fields: fields{2024, 9, 22},
			want:   LocalDate{2024, 10, 24},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ldt := LocalDate{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			assert.Equalf(t, tt.want, ldt.ToSolar(), "ToSolar()")
		})
	}
}

func TestLocalDate_ToLunar(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   LocalDate
	}{
		{
			name:   "test the to lunar conversion",
			fields: fields{2024, 10, 24},
			want:   LocalDate{2024, 9, 22},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ldt := LocalDate{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			assert.Equalf(t, tt.want, ldt.ToLunar(), "ToLunar()")
		})
	}
}

func TestStartOfMonth(t *testing.T) {
	type LocalDateTest struct {
		name     string
		input    LocalDate
		expected LocalDate
	}
	tests := []LocalDateTest{
		{
			name:     "MidMonth",
			input:    DateFromYMD(2024, 5, 15),
			expected: DateFromYMD(2024, 5, 1),
		},
		{
			name:     "EndOfMonth",
			input:    DateFromYMD(2023, 12, 31),
			expected: DateFromYMD(2023, 12, 1),
		},
		{
			name:     "StartOfMonth",
			input:    DateFromYMD(1999, 1, 1),
			expected: DateFromYMD(1999, 1, 1),
		},
		{
			name:     "LeapYearFebruary",
			input:    DateFromYMD(2020, 2, 29),
			expected: DateFromYMD(2020, 2, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.StartOfMonth()
			assert.Equalf(t, tt.expected, result, "StartOfMonth()")
		})
	}
}

func TestLocalDate_EndOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		input    LocalDate
		expected LocalDate
	}{
		{
			name:     "TC01 - 闰年2月",
			input:    LocalDate{Year: 2024, Month: 2},
			expected: LocalDate{Year: 2024, Month: 2, Day: 29},
		},
		{
			name:     "TC02 - 非闰年2月",
			input:    LocalDate{Year: 2023, Month: 2},
			expected: LocalDate{Year: 2023, Month: 2, Day: 28},
		},
		{
			name:     "TC03 - 小月 April",
			input:    LocalDate{Year: 2023, Month: 4},
			expected: LocalDate{Year: 2023, Month: 4, Day: 30},
		},
		{
			name:     "TC04 - 大月 July",
			input:    LocalDate{Year: 2023, Month: 7},
			expected: LocalDate{Year: 2023, Month: 7, Day: 31},
		},
		{
			name:     "TC05 - 本身就是月末",
			input:    LocalDate{Year: 2023, Month: 12, Day: 31},
			expected: LocalDate{Year: 2023, Month: 12, Day: 31},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.EndOfMonth()
			assert.Equalf(t, tt.expected, result, "EndOfMonth()")
		})
	}
}

func TestLocalDate_PassDays(t *testing.T) {
	tests := []struct {
		name     string
		ld       LocalDate
		date     LocalDate
		expected int
	}{
		{
			name:     "4 days difference",
			ld:       LocalDate{Year: 2025, Month: 4, Day: 5},
			date:     LocalDate{Year: 2025, Month: 4, Day: 1},
			expected: 4,
		},
		{
			name:     "Same day",
			ld:       LocalDate{Year: 2025, Month: 4, Day: 1},
			date:     LocalDate{Year: 2025, Month: 4, Day: 1},
			expected: 0,
		},
		{
			name:     "Negative days",
			ld:       LocalDate{Year: 2025, Month: 4, Day: 1},
			date:     LocalDate{Year: 2025, Month: 4, Day: 5},
			expected: -4,
		},
		{
			name:     "Across leap year",
			ld:       LocalDate{Year: 2024, Month: 2, Day: 29},
			date:     LocalDate{Year: 2020, Month: 2, Day: 28},
			expected: 1462,
		},
		{
			name:     "Cross year boundary",
			ld:       LocalDate{Year: 2025, Month: 1, Day: 1},
			date:     LocalDate{Year: 2024, Month: 12, Day: 31},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ld.PassDays(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}
