package elven

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"servus/core"
	"servus/core/modules/errorMan"
)

const filesPageSize = 2
var filesSaveTo = core.Utils.GetExecuteDir() + "/uploads"
var filesSaveToTemp = core.Utils.GetExecuteDir() + "/uploads/temp"

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
	if len(files) >= filesPageSize {
		var lastElement = len(files) - 1
		responseContent.Meta.Next = files[lastElement].ID
		files = files[:lastElement]
	}
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
	// get user permissions.
	auth := oUtils.getPipeAuth(request)
	em := errorMan.NewValidation()
	if !auth.UserAndTokenExists || !auth.IsAdmin {
		f.Send(response, errorMan.ThrowForbidden(), 403)
	}
	// get file from form and validate.
	fileFromForm, header, err := request.FormFile("file")
	if err != nil {
		em.Add("file", "bad file provided.")
		f.Send(response, em.GetJSON(), 400)
		return
	}
	defer func() {
		_ = fileFromForm.Close()
	}()
	if header == nil {
		em.Add("file", "bad file provided.")
		f.Send(response, em.GetJSON(), 400)
		return
	}
	// create temp file.
	err = os.MkdirAll(filesSaveToTemp, os.ModePerm)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	tempFile, err := os.CreateTemp(filesSaveToTemp, "tmp-*")
	if err != nil {
		f.err500(response, request, err)
		return
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()
	// generate hash and write data to temp file.
	md5hr := md5.New()
	buf := make([]byte, 8192)
	for {
		var n int
		n, err = fileFromForm.Read(buf)
		if err == io.EOF {
			break
		}
		_, err = tempFile.Write(buf[:n])
		if err != nil {
			return
		}
		_, err = md5hr.Write(buf[:n])
	}
	_ = fileFromForm.Close()
	hash := hex.EncodeToString(md5hr.Sum(nil))
	// check is file exists.
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
	// get file extension.
	var extension = header.Filename
	extension = filepath.Ext(extension)
	// create file and dirs in system.
	var hashDirs, _ = oUtils.generateDirsByHash(hash)
	// newFileDir - full path to file dir like D:/images/files/8c/9d
	var newFileDir = fmt.Sprintf("%v/%v/", filesSaveTo, hashDirs)
	err = os.MkdirAll(newFileDir, os.ModePerm)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	// newFileName - ULIDSTRING.jpg
	filenameULID, err := oUtils.generateULID()
	if err != nil {
		f.err500(response, request, err)
		return
	}
	var newFileName = filenameULID + extension
	// newFilePathLocal - 1c/2d/1513cd2345e.jpg
	var newFilePathLocal = hashDirs + "/" + newFileName
	// newFilePath - full path to file like D:/images/files/1c/2d/1513cd2345e.jpg
	var newFilePath = fmt.Sprintf("%v/%v", newFileDir, newFileName)
	// move and rename temp file to folder
	_ = tempFile.Close()
	err = os.Rename(tempFile.Name(), newFilePath)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	// create file model
	fileInDB = ModelFile{
		UserID: auth.User.ID,
		Hash: hash,
		Path: newFilePathLocal,
		Name: newFileName,
		OriginalName: header.Filename,
		Extension: extension,
		Size: header.Size,
	}
	err = fileInDB.create()
	if err != nil {
		f.err500(response, request, err)
		return
	}
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
	var fullPath = filesSaveTo + "/" + file.Path
	err = os.Remove(fullPath)
	if err == nil {
		err = oUtils.deleteEmptyDirsRecursive(filesSaveTo, file.Path)
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
