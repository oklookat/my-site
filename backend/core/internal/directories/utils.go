package directories

import (
	"errors"
	"fmt"
)

func createError(message string) error {
	return errors.New("[core/directories] " + message)
}

func wrapError(err error, message string) error {
	return fmt.Errorf("[core/directories] %v. %w", message, err)
}
