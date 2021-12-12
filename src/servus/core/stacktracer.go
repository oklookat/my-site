package core

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type StackTrace struct {
	timestamp string
	trace string
}

func (s *StackTrace) Set(err error) {
	s.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	s.trace = fmt.Sprintf("TIMESTAMP: %v\n", s.timestamp)
	if err == nil {
		s.trace = s.trace + "Error is missing."
	}
	s.trace = s.trace + "TRACE:\n\n"
	s.trace = s.trace + fmt.Sprintf("%+v", err)
}

// GetReader - returns trace string/io.Reader.
func (s *StackTrace) GetReader() io.Reader {
	return strings.NewReader(s.trace)
}

func (s *StackTrace) GetTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}