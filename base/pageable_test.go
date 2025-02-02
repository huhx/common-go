package base

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPageableToStringMethod(t *testing.T) {
	pageable := Pageable{
		PageIndex: 1,
		PageSize:  20,
	}
	result := fmt.Sprint(pageable)

	assert.Equal(t, "PageIndex: 1, PageSize: 20", result)
}
