package filer

import (
	"mime/multipart"
	"os"

	"errors"
)

var (
	ErrBadFileProvided = errors.New("[filer] bad file provided")
	ErrWrongHash       = errors.New("[filer] hash min length 6 bytes")
)

type processedFromForm struct {
	Hash      string
	Temp      *os.File
	Header    *multipart.FileHeader
	Extension string
}
