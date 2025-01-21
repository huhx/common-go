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
			got, err := DateFromString(tt.args.dateString, tt.args.pattern)

			assert.Equalf(t, tt.hasError, err != nil, "DateFromString(%v, %v)", tt.args.dateString, tt.args.pattern)
			assert.Equalf(t, tt.want, got, "DateFromString(%v, %v)", tt.args.dateString, tt.args.pattern)
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
