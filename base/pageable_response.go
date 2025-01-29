package base

type PageableResponse[T any] struct {
	Data      []T   `json:"data"`      // 当前页的数据列表，使用泛型类型 T
	Total     int64 `json:"total"`     // 总记录数
	TotalPage int64 `json:"totalPage"` // 总页数
	PageIndex int   `json:"pageIndex"` // 当前页码
	PageSize  int   `json:"pageSize"`  // 每页记录数
}

func NewPageableResponse[T any](data []T, total int64, page Pageable) PageableResponse[T] {
	return PageableResponse[T]{
		Data:      data,
		Total:     total,
		TotalPage: (total + int64(page.PageSize) - 1) / int64(page.PageSize),
		PageIndex: page.PageIndex,
		PageSize:  page.PageSize,
	}
}

func (resp PageableResponse[T]) IsFirst() bool {
	return resp.PageIndex == 0
}

func (resp PageableResponse[T]) IsLast() bool {
	return int64(resp.PageIndex)+1 == resp.TotalPage
}
