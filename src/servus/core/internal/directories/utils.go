package directories

import "github.com/pkg/errors"

func createError(message string) error {
	return errors.New("[core/directories] " + message)
}

func wrapError(err error, message string) error {
	return errors.Wrap(err, "[core/directories] "+message)
}
