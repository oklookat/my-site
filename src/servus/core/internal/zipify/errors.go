package zipify

import "errors"

var ERR_NOT_INIT = errors.New("[zipify]: not initialized. Maybe before you not called New() or called GetBytesAndClose()?")
var ERR_NIL_POINTER = errors.New("[zipify]: argument has nil pointer")
