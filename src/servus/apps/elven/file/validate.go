package file

import (
	"net/url"
	"servus/apps/elven/base"
	"strconv"
	"strings"
)

// Paginate - get paginated files by params.
type Paginate struct {
	Page  int
	Start string
	By    string
}

func (f *Paginate) Validate(params url.Values, isAdmin bool) base.Validator {
	val := validate.Create()
	// "start" param.
	var start = params.Get("start")
	if len(start) == 0 {
		start = "newest"
	}
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isStartInvalid {
		val.Add("start")
	} else {
		if isNewest {
			start = "DESC"
		} else if isOldest {
			start = "ASC"
		}
	}
	f.Start = start
	// "by" param.
	var by = params.Get("by")
	if len(by) == 0 {
		by = "created"
	}
	var isByCreated = strings.EqualFold(by, "created")
	var isByInvalid = !isByCreated
	var isByForbidden = (isByCreated) && !isAdmin
	if isByInvalid || isByForbidden {
		val.Add("created")
	}
	switch by {
	case "created":
		by = "created_at"
	}
	f.By = by
	// "page" param.
	var pageStr = params.Get("page")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		val.Add("page")
	} else {
		f.Page = page
	}
	return nil
}
