package errorMan

type eError struct {
	StatusCode int    `json:"statusCode"`
	ErrorCode  string `json:"errorCode"`
}

// basic error without unnecessary information.
type ePrimitive struct {
	eError
	Message string `json:"message"`
}
