package limiter

import (
	"io"
	"net/http"
	"regexp"
)

type Instance struct {
	active bool
	maxSize int64
	except []string
}

func New(active bool, maxSizeMB int64, except []string) *Instance {
	for index := range except {
		except[index] = normalizePath(except[index])
	}
	return &Instance{active: active, maxSize: maxSizeMB, except: except}
}

func (i *Instance) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !i.active || i.checkExcept(request.URL.Path) {
			next.ServeHTTP(writer, request)
		}
		var passed = i.Check(writer, request.Body)
		if !passed {
			writer.WriteHeader(413)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (i *Instance) Check(w http.ResponseWriter, r io.ReadCloser) (passed bool) {
	var mbToByte = i.maxSize << 4
	limitedReader := http.MaxBytesReader(w, r, mbToByte)
	_, err := io.ReadAll(limitedReader)
	return err.Error() != "http: request body too large"
}

func (i *Instance) checkExcept(path string) (skip bool) {
	path = normalizePath(path)
	for index := range i.except {
		if i.except[index] == path {
			return true
		}
	}
	return false
}

// normalizePath - from path like /hello or ///hello// make /hello/.
func normalizePath(path string) string {
	regex := regexp.MustCompile(`//+`)
	path = regex.ReplaceAllString(path, "/")
	return path
}