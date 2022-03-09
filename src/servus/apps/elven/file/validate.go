package file

import (
	"net/url"
	"servus/apps/elven/base"
	"strconv"
	"strings"
)

func ValidateGetParams(params url.Values, isAdmin bool) (bodyParams *base.FileGetParams, err error) {

	var validationErr = base.ValidationError{}

	// "start" param.
	var start = params.Get("start")
	if len(start) == 0 {
		start = "newest"
	}
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isStartInvalid {
		validationErr.New("start")("invalid value")
		err = &validationErr
		return
	}
	if isNewest {
		start = "DESC"
	} else if isOldest {
		start = "ASC"
	}
	bodyParams = &base.FileGetParams{}
	bodyParams.Start = start

	// "by" param.
	var by = params.Get("by")
	if len(by) == 0 {
		by = "created"
	}
	var isByCreated = strings.EqualFold(by, "created")
	var isByInvalid = !isByCreated
	var isByForbidden = (isByCreated) && !isAdmin
	if isByInvalid || isByForbidden {
		if isByInvalid {
			validationErr.New("by")("invalid value")
			err = &validationErr
		} else {
			validationErr.New("by")("not allowed")
			err = &validationErr
		}
		return
	}
	switch by {
	case "created":
		by = "created_at"
	}
	bodyParams.By = by

	// "page" param.
	var pageStr = params.Get("page")
	if len(pageStr) == 0 {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		validationErr.New("page")("invalid value")
		err = &validationErr
		return
	}
	bodyParams.Page = page

	// "extensions" param.
	var extensions = params.Get("extensions")
	if len(extensions) > 0 {
		var extensionsSlice = strings.Split(extensions, ",")
		bodyParams.Extensions = extensionsSlice
	}

	// "filename" param.
	var filename = params.Get("filename")
	if len(filename) > 0 {
		bodyParams.Filename = &filename
	}

	return
}
