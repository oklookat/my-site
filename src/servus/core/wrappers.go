package core

import (
	"net/http"
	"servus/core/modules/errorCollector"
	"servus/core/modules/logger"
)

// NewBaseController - create BaseController.
func NewBaseController() *BaseController {
	return &BaseController{Logger: Logger, EC: errorCollector.New()}
}

// BaseController - template for controller with cool features.
type BaseController struct {
	// Logger - its a logger.
	Logger *logger.Logger
	// EC - errorCollector. Used to collect errors and send it in JSON.
	EC *errorCollector.ErrorCollector
}

// Send - wrapper for http.ResponseWriter. Sends response and clear errorCollector errors.
func (b *BaseController) Send(response http.ResponseWriter, body string, statusCode int) {
	b.EC.Clear()
	response.WriteHeader(statusCode)
	_, err := response.Write([]byte(body))
	if err != nil {
		b.Logger.Error("baseController: failed to send response. Error: %v", err.Error())
	}
}
