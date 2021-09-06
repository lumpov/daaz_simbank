package log

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

// Writer struct to r/w and log
type Writer struct {
	rw     io.ReadWriter
	prefix string
}

// NewWriter instance
func NewWriter(rw io.ReadWriter, prefix string) *Writer {
	return &Writer{
		rw:     rw,
		prefix: prefix,
	}
}

func (l *Writer) Read(p []byte) (n int, err error) {
	n, err = l.rw.Read(p)
	if n > 0 {
		logrus.Debugf(InfoColor, fmt.Sprintf("[%s] Reading: %s", l.prefix, strings.ReplaceAll(string(p[:n]), "\r\n", " ")))
	}
	return n, err
}

func (l *Writer) Write(p []byte) (n int, err error) {
	n, err = l.rw.Write(p)
	if n > 0 {
		logrus.Debugf(GreenColor, fmt.Sprintf("[%s] Writing: %s", l.prefix, strings.ReplaceAll(string(p[:n]), "\r\n", " ")))
	}
	return n, err
}
