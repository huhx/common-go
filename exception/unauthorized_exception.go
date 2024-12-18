package exception

type Unauthorized struct {
	Content string
}

func (c Unauthorized) Code() int {
	return 403
}

func (c Unauthorized) Message() string {
	if c.Content != "" {
		return c.Content
	} else {
		return "Unauthorized"
	}
}
