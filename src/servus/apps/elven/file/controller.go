package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
	"servus/core/external/filer"
	"servus/core/external/utils"
	"strings"
)

// ALL ROUTES PROTECTED BY ADMIN ONLY MIDDLEWARE.

// get paginated files (GET).
func getAll(response http.ResponseWriter, request *http.Request) {
	var err error
	var h = call.Http.Get(request)

	// get pipe.
	var pipe = pipe.GetByContext(request)
	var isAdmin = pipe.IsAdmin()

	// validate/filter.
	validatedBody := &base.FileGetParams{}
	if err = ValidateGetParams(validatedBody, request.URL.Query(), isAdmin); err != nil {
		h.Send("bad request", 400, nil)
		return
	}

	// get paginated.
	pag := model.File{}
	files, totalPages, err := pag.GetPaginated(validatedBody)
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// generate response with pagination.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = validatedBody.Page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = model.FilePageSize
	responseContent.Data = files

	// make json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// upload file (POST).
func upload(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)
	var auth = pipe.GetByContext(request)

	// get from form.
	var tempDir = call.Config.Uploads.Temp
	processed, err := filer.ProcessFromForm(request, "file", tempDir)
	if err != nil {
		if errors.Is(err, filer.ErrBadFileProvided) {
			h.Send("bad file provided", 400, err)
			return
		}
		h.Send("", 500, err)
		return
	}

	// find duplicates in db by hash.
	var hash = processed.Hash
	var fileInDB = model.File{Hash: hash}
	found, err := fileInDB.FindByHash()
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// exists?

	if found {
		// send found file.
		_ = os.Remove(processed.Temp.Name())
		fileJSON, err := json.Marshal(fileInDB)
		if err != nil {
			h.Send("", 500, err)
			return
		}
		h.Send(string(fileJSON), 409, err)
		return
	}

	// create new.

	// generate directories struct by hash.
	dirsHash, err := filer.GenerateDirsByHash(hash)
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// generate ULID filename.
	filenameULID, err := utils.GenerateULID()
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// get & correct extension.
	var extension = processed.Extension
	extension = strings.ToLower(extension)
	switch extension {
	case "jpeg":
		extension = "jpg"
	case "mpeg":
		extension = "mpg"
	}

	// concat random name and extension.
	var newFileName string
	if len(extension) > 0 {
		newFileName = filenameULID + "." + extension
	} else {
		newFileName = filenameULID
	}

	// concat dirs struct and filename. Output like: 1c/2d/1513cd2345e.jpg.
	var newFilePathLocal = dirsHash + "/" + newFileName
	var uploadTo = call.Config.Uploads.To

	// full path where temp file will be moved. Like: D:/images/files/1c/2d/.
	var saveAt = fmt.Sprintf("%v/%v/", uploadTo, dirsHash)

	// full path to file. Like: D:/images/files/1c/2d/1513cd2345e.jpg.
	var newFilePath = saveAt + "/" + newFileName

	// make dirs.
	if err = os.MkdirAll(saveAt, os.ModePerm); err != nil {
		h.Send("", 500, err)
		return
	}

	defer func() {
		// if something goes wrong, delete created dirs.
		if err != nil {
			_ = filer.DeleteEmptyDirsRecursive(saveAt)
		}
	}()

	// move file from temp dir to created dir.
	if err = os.Rename(processed.Temp.Name(), newFilePath); err != nil {
		h.Send("", 500, err)
		return
	}

	defer func() {
		// if something goes wrong, delete file.
		if err != nil {
			_ = os.Remove(newFilePath)
		}
	}()

	// create model.
	var filename = processed.Header.Filename
	var size = processed.Header.Size
	fileInDB = model.File{
		UserID:       auth.GetID(),
		Hash:         hash,
		Path:         newFilePathLocal,
		Name:         newFileName,
		OriginalName: filename,
		Extension:    extension,
		Size:         size,
	}

	// save to db.
	if err = fileInDB.Create(); err != nil {
		h.Send("", 500, err)
		return
	}

	defer func() {
		// if something goes wrong, delete from db.
		if err != nil {
			_ = fileInDB.DeleteByID()
		}
	}()

	// send to user.
	fileJSON, err := json.Marshal(&fileInDB)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(fileJSON), 200, err)
}

// delete by ID (DELETE).
func deleteOne(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var file = model.File{ID: id}
	found, err := file.FindByID()
	if err != nil {
		h.Send("", 500, err)
		return
	}
	if !found {
		h.Send("", 404, err)
		return
	}

	// delete file from dir.
	var fullPath = call.Config.Uploads.To + "/" + file.Path
	err = os.Remove(fullPath)
	if err == nil {
		// delete dirs if empty.
		fullPath, _ = filepath.Split(fullPath)
		_ = filer.DeleteEmptyDirsRecursive(fullPath)
	}

	// delete file from db.
	if err = file.DeleteByID(); err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send("", 200, err)
}
