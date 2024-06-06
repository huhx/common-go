package exception

type NotFound struct {
	Content string
}

func (c NotFound) Code() int {
	return 404
}

func (c NotFound) Message() string {
	if c.Content != "" {
		return c.Content
	} else {
		return "Resource not found"
	}
}
