package zipify

import (
	"archive/zip"
	"bytes"
	"io"
)

type ZipFile struct {
	isInit bool
	buffer *bytes.Buffer
	writer *zip.Writer
}

// create new zip file.
func New() *ZipFile {
	var archive = &ZipFile{}

	// create buffer where .zip be.
	archive.buffer = new(bytes.Buffer)

	// create zip archive.
	archive.writer = zip.NewWriter(archive.buffer)
	archive.isInit = true
	return archive
}

// add file to archive. After you added all files you MUST call GetRAW() for closing archive.
//
// filename: filename with extension like 'myArchive.zip'
func (z *ZipFile) AddFile(filename string, data io.Reader) error {
	if !z.isInit || z.writer == nil {
		return ErrNotInit
	}
	if data == nil {
		return ErrNilPointer
	}

	// add file.
	var zipWriter, err = z.writer.Create(filename)
	if err != nil {
		return err
	}

	// write content to file.
	var buf = make([]byte, 1024)
	for {
		n, err := data.Read(buf)
		if err != nil {
			if err == io.EOF {
				// there is no more data to read.
				break
			}
			return err
		}
		if n < 1 {
			continue
		}
		if _, err = zipWriter.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// get archive in bytes.Buffer. Also closes archive.
func (z *ZipFile) GetBytesAndClose() *bytes.Buffer {
	z.isInit = false
	if z.writer == nil {
		return nil
	}
	_ = z.writer.Close()
	return z.buffer
}
