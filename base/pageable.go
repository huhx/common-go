package base

import (
	"fmt"
	"github.com/huhx/common-go/util"
)

type Pageable struct {
	PageIndex int `json:"pageIndex" example:"0"`
	PageSize  int `json:"pageSize" example:"20"`
}

func NewPageable(pageIndex, pageSize string) Pageable {
	return Pageable{util.StringToInt(pageIndex), util.StringToInt(pageSize)}
}

func (p Pageable) Limit() int {
	return p.PageSize
}

func (p Pageable) Offset() int {
	return p.PageIndex * p.PageSize
}

func (p Pageable) String() string {
	return fmt.Sprintf("PageIndex: %d, PageSize: %d", p.PageIndex, p.PageSize)
}
