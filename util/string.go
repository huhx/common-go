package util

import "strconv"

func StringToInt64(text string) int64 {
	number, _ := strconv.ParseInt(text, 10, 64)
	return number
}

func Int64ToString(number int64) string {
	return strconv.FormatInt(number, 10)
}
