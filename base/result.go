package base

type Result[R any] struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    R      `json:"data"`
}

type ListResult[R any] struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    []R    `json:"data"`
}

type PageableResult[R any] struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    R      `json:"data"`
}
