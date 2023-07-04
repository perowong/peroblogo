package controller

type ErrCodeType int

const (
	ErrCodeOk    ErrCodeType = 0
	ErrCodeParam ErrCodeType = iota + 1000
	ErrCodeServer
	ErrCodeToken
)

var CodeMap = map[ErrCodeType]string{
	ErrCodeParam:  "param error",
	ErrCodeServer: "server error",
}
