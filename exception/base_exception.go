package exception

type Exception interface {
	Code() int
	Message() string
}
