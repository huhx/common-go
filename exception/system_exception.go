package exception

type System struct {
	Content string
}

func (c System) Code() int {
	return 500
}

func (c System) Message() string {
	if c.Content != "" {
		return c.Content
	} else {
		return "Internal Error"
	}
}
