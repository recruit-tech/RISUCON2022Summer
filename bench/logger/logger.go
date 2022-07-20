package logger

import (
	"bytes"
	"errors"
	"io"
	"log"
)

var (
	Public  *log.Logger
	Private *log.Logger

	buf *bytes.Buffer
)

func Init(public, private io.Writer, debug bool) {
	f := log.Lmicroseconds
	if debug {
		f = log.Lmicroseconds | log.Lshortfile
	}

	buf = new(bytes.Buffer)
	pri := io.MultiWriter(private, buf)
	pub := io.MultiWriter(public, pri)

	Public = log.New(pub, "", 0)
	Private = log.New(pri, "", f)
}

func GetPrivateLogs() (string, error) {
	if buf == nil {
		return "", errors.New("private log buffer was not found")
	}
	return buf.String(), nil
}
