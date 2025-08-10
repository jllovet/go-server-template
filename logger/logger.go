package logger

import (
	"io"
	"log"
	"os"
)

func New(out io.Writer, prefix string) *log.Logger {
	if out == nil {
		out = os.Stdout
	}
	return log.New(
		out,
		prefix,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile|log.Lmsgprefix,
	)
}

func Stdout(prefix string) *log.Logger {
	return New(os.Stdout, prefix)
}
