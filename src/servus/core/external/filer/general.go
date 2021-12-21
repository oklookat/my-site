package filer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// DeleteEmptyDirsRecursive - delete empty dirs
//
// starts from path, and goes up deleting dirs (if empty) along the way.
//
// params:
//
// path: path like 'D:\Test\123\456\789\' or 'D:\Test\123\456\789\file.txt'
func DeleteEmptyDirsRecursive(path string) (err error) {
	path = filepath.ToSlash(path)
	inf, err := os.Stat(path)
	if err != nil {
		return
	}
	if !inf.IsDir() {
		path, _ = filepath.Split(path)
	}
	var deleteDirIfEmpty = func(path string) (resume bool, err error) {
		resume = false
		entry, err := os.ReadDir(path)
		if len(entry) == 0 {
			resume = true
			err = os.Remove(path)
		}
		return
	}
	var resume = true
	for resume {
		// delete dirs. Example:
		// it.1: delete D:\Test\123\456\789\ if empty
		// it.2: delete D:\Test\123\456\ if empty
		// it.3: delete D:\Test\123\ if empty
		// it.4: delete D:\Test\ if empty
		path = filepath.Dir(path)
		resume, err = deleteDirIfEmpty(path)
		if err != nil {
			break
		}
	}
	return
}

// GenerateDirsByHash - generate folders struct from hash like. Returns string like 1d/2c and error, if hash length less than 6 bytes.
func GenerateDirsByHash(hash string) (result string, err error) {
	if len(hash) < 6 {
		return "", errors.New("generateDirByHash: hash min length is 6 bytes")
	}
	var hashFirstTwo1 = hash[0:2]
	var hashFirstTwo2 = hash[2:4]
	result = fmt.Sprintf("%v/%v", hashFirstTwo1, hashFirstTwo2)
	return
}

// ProcessFromForm -
//
// 1. get file from request by key
//
// 2. move to temp folder
//
// 3. get MD5 hash
//
// if error - deletes temp file
//
// returns closed temp file, file header from request, and hash.
func ProcessFromForm(request *http.Request, formKey string, tempDir string) (data *processedFromForm, err error) {
	tempDir = filepath.ToSlash(tempDir)
	// get and validate
	fileFromForm, header, err := request.FormFile(formKey)
	if err != nil {
		err = ErrBadFileProvided
		return
	}
	defer func() {
		_ = fileFromForm.Close()
	}()
	if header == nil {
		err = ErrBadFileProvided
		return
	}
	// create temp dir.
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return
	}
	// create temp file.
	tempFile, err := os.CreateTemp(tempDir, "tmp-*")
	if err != nil {
		return
	}
	defer func() {
		// always close.
		_ = tempFile.Close()
		if err != nil {
			// delete if something goes wrong.
			_ = os.Remove(tempFile.Name())
		}
	}()
	// generate hash and write data to temp file.
	md5hr := md5.New()
	buf := make([]byte, 8192)
	for {
		var n int
		n, err = fileFromForm.Read(buf)
		if err == io.EOF {
			err = nil
			break
		}
		_, err = tempFile.Write(buf[:n])
		if err != nil {
			return
		}
		_, err = md5hr.Write(buf[:n])
		if err != nil {
			return
		}
	}
	data = &processedFromForm{}
	data.Temp = tempFile
	data.Header = header
	data.Hash = hex.EncodeToString(md5hr.Sum(nil))
	return
}
