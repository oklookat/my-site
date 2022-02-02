package stacktracer

import (
	"bytes"
	"fmt"
	"io"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

// TODO: improve stacktrace.
// TODO: stacktrace not see caller except iHTTP
type Instance struct {
	timestamp string
	trace     string
}

func (i *Instance) Set(err error) {
	i.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	i.trace = fmt.Sprintf("Timestamp: %v\n", i.timestamp)
	if err == nil {
		i.trace = i.trace + "No error."
	} else {
		i.trace = i.trace + "Message:\n" + err.Error() + "\n\n"
	}
	i.trace = i.trace + i.GetStacktrace()
}

// returns trace string/io.Reader.
func (i *Instance) GetReader() io.Reader {
	return strings.NewReader(i.trace)
}

func (i *Instance) GetTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (i *Instance) GetStacktrace() string {
	buf := bytes.NewBufferString("TRACE:\n")
	pprof.Lookup("goroutine").WriteTo(buf, 1)
	return buf.String()
}
