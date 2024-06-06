package exception

type Unauthenticated struct {
	Content string
}

func (c Unauthenticated) Code() int {
	return 401
}

func (c Unauthenticated) Message() string {
	if c.Content != "" {
		return c.Content
	} else {
		return "Unauthenticated"
	}
}
