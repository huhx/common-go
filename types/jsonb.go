package types

import (
	"database/sql/driver"
	"errors"
	"github.com/goccy/go-json"
)

type JSONB map[string]interface{}

func NewJSONB(data map[string]interface{}) JSONB {
	return data
}

func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
