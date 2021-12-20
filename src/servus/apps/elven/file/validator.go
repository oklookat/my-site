package file

import (
	"net/http"
	"servus/core/external/errorMan"
	"strconv"
	"strings"
)

// getAll - validate query params when getting files list.
func (f *validator) getAll(request *http.Request, isAdmin bool) (val queryGetAll, em *errorMan.EValidation, err error) {
	em = errorMan.NewValidation()
	val = queryGetAll{}
	var queryParams = request.URL.Query()
	// validate start param.
	var start = queryParams.Get("start")
	if len(start) == 0 {
		start = "newest"
	}
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isStartInvalid {
		em.Add("start", "wrong value provided.")
	} else {
		if isNewest {
			start = "DESC"
		} else if isOldest {
			start = "ASC"
		}
	}
	val.start = start
	// validate by param.
	var by = queryParams.Get("by")
	if len(by) == 0 {
		by = "created"
	}
	var isByCreated = strings.EqualFold(by, "created")
	var isByInvalid = !isByCreated
	var isByForbidden = (isByCreated) && !isAdmin
	if isByInvalid || isByForbidden {
		em.Add("by", "wrong value provided.")
	}
	switch by {
	case "created":
		by = "created_at"
	}
	val.by = by
	// validate page param.
	var pageStr = queryParams.Get("page")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		em.Add("page", "wrong value provided.")
	} else {
		val.page = page
	}
	return
}
