package zipify

import (
	"archive/zip"
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type ZipFile struct {
	initialized bool
	buffer      *bytes.Buffer
	writer      *zip.Writer
}

// New - creates new zip file.
func New() *ZipFile {
	var archive = &ZipFile{}
	// create a buffer to write our archive to.
	archive.buffer = new(bytes.Buffer)
	// create a new zip archive.
	archive.writer = zip.NewWriter(archive.buffer)
	archive.initialized = true
	return archive
}

// AddFile - add file to archive. After you added all files you MUST call GetRAW() for closing archive.
func (z *ZipFile) AddFile(filename string, data io.Reader) error {
	if !z.initialized || z.writer == nil {
		return errors.New("[zipify]: not initialized. Maybe before you called GetRAW() or not called New()?")
	}
	if data == nil {
		return errors.New("[zipify]: data nil pointer")
	}
	// add file.
	f, err := z.writer.Create(filename)
	if err != nil {
		return err
	}
	// write content to file.
	buf := make([]byte, 1024)
	for {
		n, err := data.Read(buf)
		if err == io.EOF {
			// there is no more data to read.
			break
		}
		if err != nil {
			return err
		}
		if n > 0 {
			_, err = f.Write(buf[:n])
			if err != nil {
				return err
			}
		}
	}
	return err
}

// GetRAW - get archive in bytes.Buffer. Also closes archive.
func (z *ZipFile) GetRAW() *bytes.Buffer {
	z.initialized = false
	if z.writer == nil {
		return nil
	}
	_ = z.writer.Close()
	return z.buffer
}
