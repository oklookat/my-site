package zipify

import "errors"

var (
	ErrNotInit    = errors.New("[zipify]: not initialized. Maybe before you not called New() or called GetBytesAndClose()?")
	ErrNilPointer = errors.New("[zipify]: argument has nil pointer")
)
