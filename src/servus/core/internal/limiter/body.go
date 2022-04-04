package limiter

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"path"
)

var (
	ErrConfigNil = errors.New("[limiter] body config nil pointer")
)

type BodyConfig struct {
	// limit body?
	Active bool `json:"active"`

	// max body size in MB.
	MaxSize int64 `json:"maxSize"`

	// bypass limit this paths.
	Except []string `json:"except"`
}

type Body struct {
	config       *BodyConfig
	maxSizeBytes int64
}

func NewBody(config *BodyConfig) (*Body, error) {
	if config == nil {
		return nil, ErrConfigNil
	}

	var except = config.Except
	for index := range except {
		except[index] = normalizePath(except[index])
	}

	var instance = &Body{}
	instance.config = config

	var mbToByte = config.MaxSize * 1000 * 1000
	instance.maxSizeBytes = mbToByte

	return instance, nil
}

func (i *Body) IsActive() bool {
	if i.config == nil {
		return false
	}
	return i.config.Active
}

func (i *Body) GetMaxSize() int64 {
	if i.config == nil {
		return 0
	}
	return i.config.MaxSize
}

func (i *Body) GetExcept() []string {
	if i.config == nil {
		return nil
	}
	return i.config.Except
}

func (i *Body) Middleware(next http.Handler) http.Handler {
	if i.config == nil {
		return nil
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !i.config.Active || i.checkExcept(request.URL.Path) {
			next.ServeHTTP(writer, request)
			return
		}
		var passed = i.check(writer, request)
		if !passed {
			writer.WriteHeader(413)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (i *Body) check(w http.ResponseWriter, r *http.Request) (passed bool) {
	limitedReader := http.MaxBytesReader(w, r.Body, i.maxSizeBytes)

	defer func() {
		_ = limitedReader.Close()
	}()

	body, err := io.ReadAll(limitedReader)
	if err == nil {
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return err == nil || err.Error() != "http: request body too large"
}

func (i *Body) checkExcept(path string) (skip bool) {
	path = normalizePath(path)
	for index := range i.config.Except {
		if i.config.Except[index] == path {
			return true
		}
	}
	return false
}

// from path like /hello or ///hello// make /hello/.
func normalizePath(to string) string {
	return path.Clean(to)
}
