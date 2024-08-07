package types

import (
	"database/sql/driver"
	"fmt"
	"github.com/goccy/go-json"
)

type StringArray []string

func NewStringArray(element string) *StringArray {
	return &StringArray{element}
}

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
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

func (a StringArray) IndexOf(elem string) int {
	for i, v := range a {
		if v == elem {
			return i
		}
	}
	return -1
}
