package elven

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/h2non/filetype"
	"io"
	"net/http"
	"os"
	"servus/core"
	"servus/core/modules/errorMan"
)

const filesPageSize = 2
const filesInMemorySize = 250 << 20 // 250 MB
var filesSaveTo = core.Utils.GetExecuteDir() + "/uploads"
var filesSaveToTemp = core.Utils.GetExecuteDir() + "/uploads/temp"

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
	files, err := f.databaseGetAll(&val)
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
	auth := oUtils.getPipeAuth(request)
	em := errorMan.NewValidation()
	if !auth.UserAndTokenExists || !auth.IsAdmin {
		f.Send(response, errorMan.ThrowForbidden(), 403)
	}
	// get file from form and validate.
	err := request.ParseMultipartForm(filesInMemorySize)
	fileFromForm, header, err := request.FormFile("file")
	if err != nil {
		em.Add("file", "bad file provided.")
		f.Send(response, em.GetJSON(), 400)
		return
	}
	defer fileFromForm.Close()
	if header.Size < 270 {
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
		_, err := fileFromForm.Read(buf)
		if err == io.EOF {
			break
		}
		tempFile.Write(buf)
		md5hr.Write(buf)
	}
	hash := hex.EncodeToString(md5hr.Sum(nil))
	// check is file exists.
	var fileInDB, _ = f.databaseFindBy("hash", hash)
	if fileInDB != nil {
		fileJSON, err := json.Marshal(fileInDB)
		if err != nil {
			f.err500(response, request, err)
			return
		}
		f.Send(response, string(fileJSON), 200)
		return
	}
	// get type (extension).
	// https://github.com/h2non/filetype#file-header
	_, _ = fileFromForm.Seek(0, io.SeekStart)
	fileFromFormHeader := make([]byte, 261)
	_, err = fileFromForm.Read(fileFromFormHeader)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	fType, err := filetype.Match(fileFromFormHeader)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	// create file and dirs in system.
	var hashDirs, _ = oUtils.generateDirsByHash(hash)
	// newFileDir - full path to file dir like D:/images/files/8c/9d
	var newFileDir = fmt.Sprintf("%v/%v/", filesSaveTo, hashDirs)
	err = os.MkdirAll(newFileDir, os.ModePerm)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	// newFileName - 1513cd2345e.jpg
	var newFileName = hash + "." + fType.Extension
	// newFilePathLocal - 1c/2d/1513cd2345e.jpg
	var newFilePathLocal = hashDirs + "/" + newFileName
	// newFilePath - full path to file like D:/images/files/1c/2d/1513cd2345e.jpg
	var newFilePath = fmt.Sprintf("%v/%v", newFileDir, newFileName)
	// move and rename temp file to folder
	tempFile.Close()
	err = os.Rename(tempFile.Name(), newFilePath)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	// create file model
	var newFileModel = ModelFile{}
	newFileModel.UserID = auth.User.ID
	newFileModel.Hash = hash
	newFileModel.Path = newFilePathLocal
	newFileModel.Name = newFileName
	newFileModel.OriginalName = header.Filename
	newFileModel.Extension = fType.Extension
	newFileModel.Size = header.Size
	err = f.databaseCreate(&newFileModel)
	if err != nil {
		f.err500(response, request, err)
		return
	}
	fileJSON, err := json.Marshal(&newFileModel)
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
	found, err := f.databaseFind(id)
	if err != nil {
		f.Send(response, errorMan.ThrowNotFound(), 404)
		return
	}
	// delete dirs if empty
	var fullPath = filesSaveTo + "/" + found.Path
	err = os.Remove(fullPath)
	if err == nil {
		err = oUtils.deleteEmptyDirsRecursive(filesSaveTo, found.Path)
	}
	err = f.databaseDelete(found.ID)
	f.Send(response, "", 200)
	return
}
