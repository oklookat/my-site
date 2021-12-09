package core

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"time"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type StackTrace struct {
	timestamp string
	trace string
}

func (s *StackTrace) Read(p []byte) (n int, err error) {
	p = []byte(s.trace)
	return len(s.trace), io.EOF
}

func (u *Utils) GetStackTrace(err error) (t StackTrace) {
	t = StackTrace{}
	t.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	stack, ok := errors.Cause(err).(stackTracer)
	if !ok {
		t.trace = "stacktrace not available."
		return t
	}
	var st = stack.StackTrace()
	t.trace = fmt.Sprintf("%+v", st[0:])
	return t
}
