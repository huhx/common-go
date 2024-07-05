package types

import (
	"database/sql/driver"
	"fmt"
	"github.com/goccy/go-json"
)

type Int64Array []int64

func (a Int64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *Int64Array) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &a)
}

func (a Int64Array) IndexOf(elem int64) int {
	for i, v := range a {
		if v == elem {
			return i
		}
	}
	return -1
}
