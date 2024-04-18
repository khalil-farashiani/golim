package contract

import "io"

type Logger interface {
	Warn(msg io.Writer)
	Error(msg io.Writer)
	Info(msg io.Writer)
}
