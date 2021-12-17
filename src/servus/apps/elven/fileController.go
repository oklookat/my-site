package elven

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"servus/core/external/errorMan"
	"servus/core/external/filer"
	"strings"

	"github.com/gorilla/mux"
)

const filesPageSize = 2

// GET url/
// params:
// cursor = ULID
// by = created (uploaded)
// start = newest (DESC), oldest (ASC)
func (f *fileRoute) getAll(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// validate query params.
	val, em, _ := f.validate.getAll(request, true)
	if em.HasErrors() {
		h.Send(em.GetJSON(), 400, nil)
		return
	}
	// get files by query params.
	files, totalPages, err := val.getAll()
	if err != nil {
		call.Logger.Error(fmt.Sprintf("files get error: %v", err.Error()))
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	// generate response with pagination.
	var responseContent = ResponseContent{}
	responseContent.Meta.CurrentPage = val.page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = filesPageSize
	responseContent.Data = files
	// make json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// POST url/
func (f *fileRoute) createOne(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	auth := AuthPipe{}
	auth.get(request)
	em := errorMan.NewValidation()
	var tempDir = call.Config.Uploads.Temp
	// get file from form.
	processed, err := filer.ProcessFromForm(request, "file", tempDir)
	if err != nil {
		if err == filer.ErrBadFileProvided {
			em.Add("file", "bad file provided.")
			h.Send(em.GetJSON(), 400, err)
			return
		}
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	// check is file exists.
	var hash = processed.Hash
	var fileInDB = FileModel{Hash: hash}
	found, err := fileInDB.findByHash()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	if found {
		_ = os.Remove(processed.Temp.Name())
		fileJSON, err := json.Marshal(fileInDB)
		if err != nil {
			h.Send(errorMan.ThrowServer(), 500, err)
			return
		}
		h.Send(string(fileJSON), 200, err)
		return
	}
	var filename = processed.Header.Filename
	var extension = filepath.Ext(filename)
	var extensionWithoutDot = strings.Replace(extension, ".", "", -1)
	var dirsHash, _ = filer.GenerateDirsByHash(hash)
	var uploadTo = call.Config.Uploads.To
	// saveAt - full path where temp file will be moved.
	var saveAt = fmt.Sprintf("%v/%v/", uploadTo, dirsHash)
	// newFileName - ULIDSTRING.jpg
	filenameULID, err := generateULID()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	var newFileName = filenameULID + extension
	// newFilePathLocal - 1c/2d/1513cd2345e.jpg
	var newFilePathLocal = dirsHash + "/" + newFileName
	// newFilePath - full path to file like D:/images/files/1c/2d/1513cd2345e.jpg
	var newFilePath = saveAt + "/" + newFileName
	// make dir and move temp file to it
	err = os.MkdirAll(saveAt, os.ModePerm)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	defer func() {
		if err != nil {
			_ = filer.DeleteEmptyDirsRecursive(saveAt)
		}
	}()
	err = os.Rename(processed.Temp.Name(), newFilePath)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	defer func() {
		if err != nil {
			_ = os.Remove(newFilePath)
		}
	}()
	// create file model
	var size = processed.Header.Size
	fileInDB = FileModel{
		UserID:       auth.User.ID,
		Hash:         hash,
		Path:         newFilePathLocal,
		Name:         newFileName,
		OriginalName: filename,
		Extension:    extensionWithoutDot,
		Size:         size,
	}
	err = fileInDB.create()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	defer func() {
		if err != nil {
			_ = fileInDB.deleteByID()
		}
	}()
	// send file to user
	fileJSON, err := json.Marshal(&fileInDB)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(fileJSON), 200, err)
	return
}

// DELETE url/id
func (f *fileRoute) deleteOne(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	var params = mux.Vars(request)
	var id = params["id"]
	var file = FileModel{ID: id}
	found, err := file.findByID()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	if !found {
		h.Send(errorMan.ThrowNotFound(), 404, err)
		return
	}
	// delete dirs if empty
	var fullPath = call.Config.Uploads.To + "/" + file.Path
	err = os.Remove(fullPath)
	if err == nil {
		fullPath, _ = filepath.Split(fullPath)
		err = filer.DeleteEmptyDirsRecursive(fullPath)
	}
	err = file.deleteByID()
	h.Send("", 200, err)
	return
}
