package article

import (
	"net/url"
	"servus/apps/elven/file"
	"servus/core/external/utils"
	"strconv"
	"strings"
)

// correct/validate fields.
func FilterAndValidate(what *Model) (isValid bool, err error) {
	// title.
	var isTitleEmpty = utils.LenRune(strings.TrimSpace(what.Title)) < 1
	if isTitleEmpty {
		what.Title = "Untitled"
	} else if utils.LenRune(what.Title) > 124 {
		return
	}

	// content.
	if utils.LenRune(what.Content) > 816000 {
		return
	}

	// cover.
	if what.CoverID != nil {
		var file = file.Model{}
		file.ID = *what.CoverID
		var isFound bool
		isFound, err = file.FindByID()
		if err != nil {
			return
		} else if !isFound {
			what.CoverID = nil
		} else {
			// can this file be a cover?
			var extension = strings.ToUpper(file.Extension)
			var isCoverable = extension == "JPG" || extension == "JPEG" || extension == "PNG" ||
				extension == "GIF" || extension == "WEBP" || extension == "MP4"
			if !isCoverable {
				what.CoverID = nil
			}
		}
	}

	isValid = true
	return
}

// validate params to get paginated articles.
func ValidateGetParams(a *GetParams, params url.Values, isAdmin bool) {
	var err error

	// "published" param
	var published = params.Get("drafts")
	if len(published) < 1 {
		a.Drafts = false
	} else {
		a.Drafts, err = strconv.ParseBool(published)
		if err != nil {
			a.Drafts = false
		}
	}
	if a.Drafts && !isAdmin {
		a.Drafts = false
	}

	// "by" param.
	var by = params.Get("by")
	if by == "created" || by == "updated" {
		if !isAdmin {
			by = "published"
		}
	} else {
		by = "published"
	}
	by += "_at"
	a.By = by

	// "newest" param.
	var newest = params.Get("newest")
	if len(newest) < 1 {
		a.Newest = true
	} else {
		a.Newest, err = strconv.ParseBool(newest)
		if err != nil {
			a.Newest = true
		}
	}

	// "page" param
	var pageStr = params.Get("page")
	if len(pageStr) < 1 {
		a.Page = 1
	} else {
		var page = 0
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
		a.Page = page
	}

	// "title" param
	var title = params.Get("title")
	if len(title) > 0 {
		a.Title = &title
	}
}
