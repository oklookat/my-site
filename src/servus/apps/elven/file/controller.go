package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
	"servus/core/external/filer"
)

// ALL ROUTES PROTECTED BY ADMIN ONLY MIDDLEWARE.

// get paginated files (GET).
func (f *Instance) getAll(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// get pipe.
	pipe := f.pipe.GetByContext(request)
	isAdmin := pipe != nil && pipe.IsAdmin()

	// validate/filter.
	body := ValidateGetParams(request.URL.Query(), isAdmin)
	if body == nil {
		h.Send("invalid request", 400, nil)
		return
	}

	// get paginated.
	pag := model.File{}
	files, totalPages, err := pag.GetPaginated(body)
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}

	// generate response with pagination.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = body.Page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = model.FilePageSize
	responseContent.Data = files

	// make json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// upload file (POST).
func (f *Instance) upload(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	auth := f.pipe.GetByContext(request)

	// get from form.
	var tempDir = call.Config.Uploads.Temp
	processed, err := filer.ProcessFromForm(request, "file", tempDir)
	if err != nil {
		if err == filer.ErrBadFileProvided {
			h.Send("bad file provided", 400, err)
			return
		}
		h.Send(f.throw.Server(), 500, err)
		return
	}

	// find duplicates in db by hash.
	var hash = processed.Hash
	var fileInDB = model.File{Hash: hash}
	found, err := fileInDB.FindByHash()
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}

	// exists?

	if found {
		// send found file.
		_ = os.Remove(processed.Temp.Name())
		fileJSON, err := json.Marshal(fileInDB)
		if err != nil {
			h.Send(f.throw.Server(), 500, err)
			return
		}
		h.Send(string(fileJSON), 200, err)
		return
	}

	// create new.

	// generate directories struct by hash.
	dirsHash, err := filer.GenerateDirsByHash(hash)
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}

	// generate ULID filename.
	filenameULID, err := call.Utils.GenerateULID()
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}

	// get extension.
	var extension = processed.Extension

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
	err = os.MkdirAll(saveAt, os.ModePerm)
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}

	defer func() {
		// if something goes wrong, delete created dirs.
		if err != nil {
			_ = filer.DeleteEmptyDirsRecursive(saveAt)
		}
	}()

	// move file from temp dir to created dir.
	err = os.Rename(processed.Temp.Name(), newFilePath)
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
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
	err = fileInDB.Create()
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
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
		h.Send(f.throw.Server(), 500, err)
		return
	}
	h.Send(string(fileJSON), 200, err)
}

// delete by ID (DELETE).
func (f *Instance) deleteOne(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var file = model.File{ID: id}
	found, err := file.FindByID()
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(f.throw.NotFound(), 404, err)
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
	err = file.DeleteByID()
	if err != nil {
		h.Send(f.throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}
