package controlTelegram

import "io"

type File struct {
	caption  *string
	filename string
	reader   io.Reader
}

func (t *File) New(caption *string, filename string, reader io.Reader) {
	t.caption = caption
	t.filename = filename
	t.reader = reader
}

func (t *File) NeedsUpload() bool {
	return true
}

func (t *File) UploadData() (string, io.Reader, error) {
	return t.filename, t.reader, nil
}

func (t *File) SendData() string {
	return t.filename
}
