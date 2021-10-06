package elven

import (
	"net/http"
	"servus/core/modules/logger"
)

type baseController struct {
	logger *logger.Logger
}

func (b *baseController) Send(w http.ResponseWriter, body string, statusCode int) {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(body))
	if err != nil {
		b.logger.Error("controller: failed to send response. Error: %v", err.Error())
	}
}

type controllerAuth struct {
	*baseController
}

type controllerArticles struct {
	*baseController
}

type controllerFiles struct {
	*baseController
}

