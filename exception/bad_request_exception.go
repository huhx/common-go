package exception

type BadRequest struct {
	Content string
}

func (c BadRequest) Code() int {
	return 400
}

func (c BadRequest) Message() string {
	if c.Content != "" {
		return c.Content
	} else {
		return "Bad Request"
	}
}
