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
	// 处理 value 为 nil 的情况
	if value == nil {
		*a = nil
		return nil
	}

	// 处理 []byte 类型的值
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	// 解析 JSON 数据到 []int64 类型
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
