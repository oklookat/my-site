package article

import (
	"encoding/json"
	"io"
	"net/url"
	"regexp"
	"servus/apps/elven/base"
	"strconv"
	"strings"
)

// params to get paginated articles.
type Paginate struct {
	// number of page.
	Page int
	// published; drafts.
	Show string
	// created; updated; published.
	By string
	// newest (DESC); oldest (ASC).
	Start string
	// true (with content); false (gives you empty content).
	Preview bool
}

// represents article body of the request that the user should send. Used in create and update methods.
type ArticleBody struct {
	CategoryID  *string `json:"category_id"`
	IsPublished *bool   `json:"is_published"`
	Title       *string `json:"title"`
	Content     *string `json:"content"`
}

// user should send this body to create/update category.
type CategoryBody struct {
	Name string `json:"name"`
}

// validate params to get paginated articles.
func (a *Paginate) Validate(params url.Values, isAdmin bool) (val base.Validator) {
	val = validate.Create()
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
	return
}

// validate request body when POST or PUT.
func (a *ArticleBody) Validate(body io.ReadCloser) (val base.Validator) {
	val = validate.Create()
	// body.
	err := json.NewDecoder(body).Decode(a)
	if err != nil {
		val.Add("body")
		return
	}
	// if all fields empty.
	var isContent = a.Content != nil
	var isTitle = a.Title != nil
	var isPublished = a.IsPublished != nil
	if !isContent && !isTitle && !isPublished {
		val.Add("body")
		return
	}
	// content.
	if isContent {
		var contentLen = call.Utils.LenRune(*a.Content)
		var contentInvalid = contentLen < 1 || contentLen > 256000
		if contentInvalid {
			val.Add("content")
		}
	}
	// title.
	if isTitle {
		var titleLen = call.Utils.LenRune(*a.Title)
		if titleLen < 1 {
			*a.Title = "Untitled"
		} else if titleLen > 124 {
			val.Add("title")
		}
	}
	// isPublished.
	return
}

func (c *CategoryBody) Validate(body io.ReadCloser) (val base.Validator) {
	val = validate.Create()
	// body.
	err := json.NewDecoder(body).Decode(c)
	if err != nil {
		val.Add("body")
		return
	}
	// replace new lines with one space
	reg, _ := regexp.Compile(`[\r\n]`)
	c.Name = reg.ReplaceAllString(c.Name, " ")
	// replace 2+ spaces with one space
	reg, _ = regexp.Compile(`[^\S]{2,}`)
	c.Name = reg.ReplaceAllString(c.Name, "")
	// remove spaces at start and end
	c.Name = strings.Trim(c.Name, " ")
	// length?
	var nameLen = call.Utils.LenRune(c.Name)
	var notValid = nameLen < 1 || nameLen > 24
	if notValid {
		val.Add("name")
	}
	return
}
