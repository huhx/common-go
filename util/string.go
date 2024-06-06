package util

import (
	"github.com/samber/lo"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomString(length int) string {
	return lo.RandomString(length, charset)
}
