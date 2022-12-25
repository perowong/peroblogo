package errcode

type ErrCodeType int

const (
	ErrCodeOk    ErrCodeType = 0
	ErrCodeParam ErrCodeType = iota + 1001
	ErrCodeServer
)

var CodeMap = map[ErrCodeType]string{
	ErrCodeParam:  "param error",
	ErrCodeServer: "server error",
}
