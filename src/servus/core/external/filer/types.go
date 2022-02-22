package filer

import (
	"mime/multipart"
	"os"

	"github.com/pkg/errors"
)

var ErrBadFileProvided = errors.New("bad file provided")

type processedFromForm struct {
	Hash      string
	Temp      *os.File
	Header    *multipart.FileHeader
	Extension string
}
