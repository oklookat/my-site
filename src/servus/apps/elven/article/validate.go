package article

import (
	"encoding/json"
	"io"
	"net/url"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
	"strconv"
	"strings"
)

// Paginate - params to get paginated articles.
type Paginate struct {
	// Page: number of page.
	Page int
	// Show: published; drafts.
	Show string
	// By: created; updated; published.
	By string
	// Start: newest (DESC); oldest (ASC).
	Start string
	// Preview: true (content < 480 symbols); false (gives you full articles).
	Preview bool
}

// Body - represents the body of the request that the user should send. Used in create and update methods.
type Body struct {
	IsPublished *bool                 `json:"is_published"`
	Title       *string               `json:"title"`
	Content     *model.ArticleContent `json:"content"`
}

// Validate - validate params to get paginated articles.
func (a *Paginate) Validate(params url.Values, isAdmin bool) base.Validator {
	var val = validate.Create()
	// "show" param
	var show = params.Get("show")
	if len(show) == 0 {
		show = "published"
	}
	var isShowPublished = show == "published"
	var isShowDrafts = show == "drafts"
	var isShowInvalid = !isShowPublished && !isShowDrafts
	var isShowForbidden = isShowDrafts && !isAdmin
	if isShowInvalid || isShowForbidden {
		val.Add("show")
	} else {
		a.Show = show
	}
	// "by" param.
	var by = params.Get("by")
	if len(by) == 0 {
		by = "published"
	}
	var isByCreated = by == "created"
	var isByPublished = by == "published"
	var isByUpdated = by == "updated"
	var isByInvalid = !isByCreated && !isByPublished && !isByUpdated
	var isByForbidden = (isByCreated || isByUpdated) && !isAdmin
	if isByInvalid || isByForbidden {
		val.Add("by")
	}
	switch by {
	case "created":
		by = "created_at"
	case "updated":
		by = "updated_at"
	case "published":
		by = "published_at"
	}
	a.By = by
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
		a.Start = start
	}
	// "preview" param.
	var preview = params.Get("preview")
	if len(preview) == 0 {
		preview = "true"
	}
	var previewBool bool
	previewBool, err := strconv.ParseBool(preview)
	if err != nil {
		val.Add("preview")
	}
	a.Preview = previewBool
	// "page" param
	var pageStr = params.Get("page")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		val.Add("page")
	} else {
		a.Page = page
	}
	return nil
}

// Validate - validate request body when POST or PUT.
func (a *Body) Validate(body io.ReadCloser) base.Validator {
	var val = validate.Create()
	// body.
	err := json.NewDecoder(body).Decode(a)
	if err != nil {
		val.Add("body")
	}
	// if all fields empty.
	var isContent = a.Content != nil
	var isTitle = a.Title != nil
	var isPublished = a.IsPublished != nil
	if !isContent && !isTitle && !isPublished {
		val.Add("body")
		return val
	}
	// content.
	if isContent {
		var contentInvalid = len(a.Content.Blocks) < 1
		if contentInvalid {
			val.Add("content")
		}
	}
	// title.
	if isTitle {
		if len(*a.Title) < 1 {
			*a.Title = "Untitled"
		} else if len(*a.Title) > 124 {
			val.Add("title")
		}
	}
	// isPublished.
	return val
}
