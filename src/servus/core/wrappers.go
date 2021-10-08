package core

import (
	"net/http"
	"servus/core/modules/logger"
)

// NewBaseController - create BaseController.
func NewBaseController() *BaseController {
	return &BaseController{Logger: Logger}
}

// BaseController - template for controller with cool features.
type BaseController struct {
	// Logger - its a logger.
	Logger *logger.Logger
}

// Send - wrapper for http.ResponseWriter. Sends response and clear errorMan errors.
func (b *BaseController) Send(response http.ResponseWriter, body string, statusCode int) {
	response.WriteHeader(statusCode)
	_, err := response.Write([]byte(body))
	if err != nil {
		b.Logger.Error("baseController: failed to send response. Error: %v", err.Error())
	}
}
