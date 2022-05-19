package zipify

import "errors"

var (
	ErrNotInit = errors.New("[zipify] not initialized. Maybe before you not called New() or called GetBytesAndClose()?")
	ErrDataNil = errors.New("[zipify] data has nil pointer")
)
