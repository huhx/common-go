package base

type PageableResponse[T any] struct {
	Data      []T   `json:"data"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"totalPage"`
	PageIndex int   `json:"pageIndex"`
	PageSize  int   `json:"pageSize"`
}

func EmptyPageableResponse[T any]() PageableResponse[T] {
	return PageableResponse[T]{}
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

func (resp PageableResponse[T]) HasNext() bool {
	return int64(resp.PageIndex)+1 < resp.TotalPage
}

func (resp PageableResponse[T]) HasPrevious() bool {
	return resp.PageIndex > 0
}

func (resp PageableResponse[T]) NextPage() Pageable {
	return Pageable{resp.PageIndex + 1, resp.PageSize}
}

func (resp PageableResponse[T]) IsFirst() bool {
	return resp.PageIndex == 0
}

func (resp PageableResponse[T]) IsLast() bool {
	return int64(resp.PageIndex)+1 == resp.TotalPage
}
