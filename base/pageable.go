package base

type Pageable struct {
	PageIndex int `json:"pageIndex" example:"0"`
	PageSize  int `json:"pageSize" example:"20"`
}

func NewPageable(pageIndex, pageSize int) Pageable {
	return Pageable{pageIndex, pageSize}
}
