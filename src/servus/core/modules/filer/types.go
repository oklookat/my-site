package filer

import (
	"github.com/pkg/errors"
	"mime/multipart"
	"os"
)

var ErrBadFileProvided = errors.New("bad file provided")

type processedFromForm struct {
	Hash string
	Temp *os.File
	Header *multipart.FileHeader
}
