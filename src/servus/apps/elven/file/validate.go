package file

import (
	"net/url"
	"servus/apps/elven/base"
	"strconv"
	"strings"
)

func ValidateGetParams(f *base.FileGetParams, params url.Values, isAdmin bool) (err error) {
	// "start" param.
	var start = params.Get("start")
	if start != "newest" && start != "oldest" {
		start = "newest"
	}
	if start == "newest" {
		start = "DESC"
	} else {
		start = "ASC"
	}
	f.Start = start

	// "by" param.
	var by = params.Get("by")
	if by != "created" {
		by = "created"
	}
	by += "_at"
	f.By = by

	// "page" param
	var pageStr = params.Get("page")
	if len(pageStr) < 1 {
		f.Page = 1
	} else {
		var page = 0
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
		f.Page = page
	}

	// "extensions" param.
	var extensions = params.Get("extensions")
	if len(extensions) > 0 {
		var extensionsSlice = strings.Split(extensions, ",")
		f.Extensions = extensionsSlice
	}

	// "filename" param.
	var filename = params.Get("filename")
	if len(filename) > 0 {
		f.Filename = &filename
	}

	return
}
