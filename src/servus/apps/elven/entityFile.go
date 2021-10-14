package elven

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"servus/core"
	"servus/core/modules/errorMan"
	"servus/core/modules/filer"
	"strings"
)

const filesPageSize = 2

// entityFile - manage files.
type entityFile struct {
	*entityBase
}

// GET url/
// params:
// cursor = ULID
// by = created (uploaded)
// start = newest (DESC), oldest (ASC)
func (f *entityFile) controllerGetAll(response http.ResponseWriter, request *http.Request) {
	val, em, _ := f.validatorControllerGetAll(request, true)
	if em.HasErrors() {
		f.Send(response, em.GetJSON(), 400)
		return
	}
	files, err := val.getAll()
	if err != nil {
		f.Logger.Error(fmt.Sprintf("files get error: %v", err.Error()))
		f.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	// generate response with pagination
	var responseContent = ResponseContent{}
	responseContent.Meta.PerPage = filesPageSize
	//if len(files) >= filesPageSize {
	//	var lastElement = len(files) - 1
	//	responseContent.Meta.Next = files[lastElement].ID
	//	files = files[:lastElement]
	//}
	responseContent.Data = files
	// make json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		f.Logger.Error(fmt.Sprintf("files response json marshal error: %v", err.Error()))
		f.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	f.Send(response, string(jsonResponse), 200)
}

// POST url/
func (f *entityFile) controllerCreateOne(response http.ResponseWriter, request *http.Request) {
	auth := PipeAuth{}
	auth.get(request)
	em := errorMan.NewValidation()
	var tempDir = core.Config.Uploads.Temp
	// get file from form.
	processed, err := filer.ProcessFromForm(request, "file", tempDir)
	if err != nil {
		if err == filer.ErrBadFileProvided {
			em.Add("file", "bad file provided.")
			f.Send(response, em.GetJSON(), 400)
			return
		}
		f.err500(response, request, err)
		return
	}
	// check is file exists.
	var hash = processed.Hash
	var fileInDB = ModelFile{Hash: hash}
	found, err := fileInDB.findByHash()
	if err != nil {
		f.err500(response, request, err)
		return
	}
	if found {
		fileJSON, err := json.Marshal(fileInDB)
		if err != nil {
			f.err500(response, request, err)
			return
		}
		f.Send(response, string(fileJSON), 200)
		return
	}
	var filename = processed.Header.Filename
	var extension = filepath.Ext(filename)
	var extensionWithoutDot = strings.Replace(extension, ".", "", -1)
	var dirsHash, _ = filer.GenerateDirsByHash(hash)
	var uploadTo = core.Config.Uploads.To
	// saveAt - full path where temp file will be moved.
	var saveAt = fmt.Sprintf("%v/%v/", uploadTo, dirsHash)
	// newFileName - ULIDSTRING.jpg
	filenameULID, err := oUtils.generateULID()
	if err != nil {
		f.err500(response, request, err)
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
		f.err500(response, request, err)
		return
	}
	defer func() {
		if err != nil {
			_ = filer.DeleteEmptyDirsRecursive(saveAt)
		}
	}()
	err = os.Rename(processed.Temp.Name(), newFilePath)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	defer func() {
		if err != nil {
			_ = os.Remove(newFilePath)
		}
	}()
	// create file model
	var size = processed.Header.Size
	fileInDB = ModelFile{
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
		f.err500(response, request, err)
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
		f.err500(response, request, err)
		return
	}
	f.Send(response, string(fileJSON), 200)
	return
}

// DELETE url/id
func (f *entityFile) controllerDeleteOne(response http.ResponseWriter, request *http.Request) {
	var params = mux.Vars(request)
	var id = params["id"]
	var file = ModelFile{ID: id}
	found, err := file.findByID()
	if err != nil {
		f.err500(response, request, err)
		return
	}
	if !found {
		f.Send(response, errorMan.ThrowNotFound(), 404)
		return
	}
	// delete dirs if empty
	var fullPath = core.Config.Uploads.To + "/" + file.Path
	err = os.Remove(fullPath)
	if err == nil {
		fullPath, _ = filepath.Split(fullPath)
		err = filer.DeleteEmptyDirsRecursive(fullPath)
	}
	err = file.deleteByID()
	f.Send(response, "", 200)
	return
}

// err403 - send an error if the user is not allowed to do something.
func (f *entityFile) err500(response http.ResponseWriter, request *http.Request, err error) {
	f.Logger.Warn("entityFile code 500 at: %v. Error: %v", request.URL.Path, err.Error())
	f.Send(response, errorMan.ThrowServer(), 500)
	return
}
