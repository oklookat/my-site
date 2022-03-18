package banhammer

import "errors"

func createError(message string) error {
	return errors.New("[banhammer] " + message)
}
