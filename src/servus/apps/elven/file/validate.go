package file

import (
	"net/url"
	"servus/apps/elven/base"
	"strconv"
	"strings"
)

func ValidateGetParams(params url.Values, isAdmin bool) (bodyParams *base.FileGetParams) {

	// "start" param.
	var start = params.Get("start")
	if len(start) == 0 {
		start = "newest"
	}
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isStartInvalid {
		return nil
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
		return nil
	}
	switch by {
	case "created":
		by = "created_at"
	}
	bodyParams.By = by

	// "page" param.
	var pageStr = params.Get("page")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return nil
	}
	bodyParams.Page = page

	return
}
