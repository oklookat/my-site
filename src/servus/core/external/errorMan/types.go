package errorMan

// eError - must have for errors.
type eError struct {
	StatusCode int    `json:"statusCode"`
	ErrorCode  string `json:"errorCode"`
}

// ePrimitive - basic error without unnecessary information.
type ePrimitive struct {
	eError
	Message string `json:"message"`
}
