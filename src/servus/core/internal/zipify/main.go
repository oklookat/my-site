package zipify

import (
	"archive/zip"
	"bytes"
	"io"
)

type ZipFile struct {
	initialized bool
	buffer *bytes.Buffer
	writer *zip.Writer
}

func (z *ZipFile) AddFile(filename string, data io.Reader) error {
	if !z.initialized {
		// Create a buffer to write our archive to.
		z.buffer = new(bytes.Buffer)
		// Create a new zip archive.
		z.writer = zip.NewWriter(z.buffer)
		z.initialized = true
	}
	f, err := z.writer.Create(filename)
	if err != nil {
		return err
	}
	var body []byte
	_, err = data.Read(body)
	_, err = f.Write(body)
	if err != nil {
		return err
	}
	// Make sure to check the error on Close.
	err = z.writer.Close()
	if err != nil {
		return err
	}
	return err
}

func (z *ZipFile) GetFile() *bytes.Buffer {
	return z.buffer
}