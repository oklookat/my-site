package stacktracer

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// TODO: improve stacktrace.
type Instance struct {
	timestamp string
	trace     string
}

func (i *Instance) Set(err error) {
	i.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	i.trace = fmt.Sprintf("TIMESTAMP: %v\n", i.timestamp)
	if err == nil {
		i.trace = i.trace + "Error is missing."
	}
	i.trace = i.trace + "TRACE:\n\n"
	i.trace = i.trace + fmt.Sprintf("%+v", err)
}

// returns trace string/io.Reader.
func (i *Instance) GetReader() io.Reader {
	return strings.NewReader(i.trace)
}

func (i *Instance) GetTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
