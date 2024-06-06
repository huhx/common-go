package exception

type Validation struct {
	Content string
}

func (c Validation) Code() int {
	return 400
}

func (c Validation) Message() string {
	if c.Content != "" {
		return c.Content
	} else {
		return "Validation Failed"
	}
}
