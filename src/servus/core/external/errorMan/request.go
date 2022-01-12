package errorMan

import "encoding/json"

type RequestError struct {
}

// ThrowServer - server error JSON (500).
func (r RequestError) Server() string {
	err := ePrimitive{
		eError{StatusCode: 500, ErrorCode: "E_SERVER"},
		"Server error.",
	}
	var bytes, _ = json.Marshal(err)
	return string(bytes)
}

// ThrowNotAuthorized - unauthorized error JSON (401).
func (r RequestError) NotAuthorized() string {
	err := ePrimitive{
		eError{StatusCode: 401, ErrorCode: "E_UNAUTHORIZED"},
		"You not authorized.",
	}
	var bytes, _ = json.Marshal(err)
	return string(bytes)
}

// ThrowForbidden - forbidden error JSON (403). Uses when wrong credentials or if user authorized, but cannot access to secure resource.
func (r RequestError) Forbidden() string {
	err := ePrimitive{
		eError{StatusCode: 403, ErrorCode: "E_ACCESS_DENIED"},
		"You cannot access this resource.",
	}
	var bytes, _ = json.Marshal(err)
	return string(bytes)
}

// ThrowNotFound - not found error JSON (404).
func (r RequestError) NotFound() string {
	err := ePrimitive{
		eError{StatusCode: 404, ErrorCode: "E_NOTFOUND"},
		"Not found.",
	}
	var bytes, _ = json.Marshal(err)
	return string(bytes)
}