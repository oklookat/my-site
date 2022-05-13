package article

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"servus/apps/elven/base"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

// ALL HANDLERS PROTECTED BY SAFE METHODS MIDDLEWARE.

// get paginated articles (GET url/).
func getArticles(response http.ResponseWriter, request *http.Request) {
	var err error
	var h = call.Http.Get(request)
	var pipe = pipe.GetByContext(request)
	var isAdmin = pipe.IsAdmin()

	// validate.
	validatedBody := &GetParams{}
	ValidateGetParams(validatedBody, request.URL.Query(), isAdmin)

	// get.
	article := Model{}
	articles, totalPages, err := article.GetPaginated(validatedBody, isAdmin)
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// fill response.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = validatedBody.Page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = pageSize
	responseContent.Data = articles

	// send.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// get one article (GET url/id).
func getArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)
	var isAdmin = false

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	var pAuth = pipe.GetByContext(request)
	isAdmin = pAuth.IsAdmin()

	// find.
	var article = Model{ID: id}
	found, err := article.FindByID(isAdmin)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	if !found {
		h.Send("", 404, err)
		return
	}

	// if not published and user not admin, access denied.
	if !article.IsPublished && !isAdmin {
		h.Send("", 403, err)
		return
	}

	// send.
	articleJson, err := json.Marshal(article)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// сreate new article (POST url/).
func createArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// decode body.
	var bodyStruct = &Body{}
	var err = json.NewDecoder(request.Body).Decode(bodyStruct)
	if err != nil {
		h.Send("", 400, err)
		return
	}
	var article = &Model{}
	article.Title = bodyStruct.Title
	article.Content = bodyStruct.Content
	article.IsPublished = bodyStruct.IsPublished
	if len(bodyStruct.CoverID) > 0 {
		article.CoverID = new(string)
		*article.CoverID = bodyStruct.CoverID
	}

	// validate.
	var isValid bool
	isValid, err = FilterAndValidate(article)
	if err != nil {
		h.Send("", 500, err)
		return
	} else if !isValid {
		h.Send("", 400, err)
		return
	}

	var pAuth = pipe.GetByContext(request)
	article.UserID = pAuth.GetID()

	// create.
	if err = article.Create(); err != nil {
		h.Send("", 500, err)
		return
	}

	// send created.
	articleJson, err := json.Marshal(&article)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// update specific fields by json-patch (PATCH url/id).
func updateArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)
	var id = h.GetRouteArgs()["id"]
	var article = Model{ID: id}
	found, err := article.FindByID(true)
	if err != nil {
		h.Send("", 500, err)
		return
	} else if !found {
		h.Send("", 404, err)
		return
	}

	// body to bytes -> bytes to JSON patch.
	var buf = new(bytes.Buffer)
	_, err = io.Copy(buf, request.Body)
	if err != nil {
		h.Send("", 400, err)
		return
	}
	var bytesPatch = buf.Bytes()
	var patch jsonpatch.Patch
	patch, err = jsonpatch.DecodePatch(bytesPatch)
	if err != nil {
		h.Send("", 400, err)
		return
	}

	// found article -> JSON.
	var articleJSON []byte
	articleJSON, err = json.Marshal(article)
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// apply patch to found article (JSON).
	var patchedArticleJSON []byte
	patchedArticleJSON, err = patch.Apply(articleJSON)
	if err != nil {
		h.Send("", 400, err)
		return
	}

	// patched JSON -> to Model.
	var patched = &Model{}
	err = json.Unmarshal(patchedArticleJSON, patched)
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// validate.
	var isValid bool
	isValid, err = FilterAndValidate(patched)
	if err != nil {
		h.Send("", 500, err)
		return
	} else if !isValid {
		h.Send("", 400, err)
		return
	}
	article = *patched

	// update.
	if err = article.Update(); err != nil {
		h.Send("", 500, err)
		return
	}

	// send patched.
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(jsonArticle), 200, err)
}

// DELETE url/id.
func deleteArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)
	var id = h.GetRouteArgs()["id"]
	var article = Model{ID: id}
	found, err := article.FindByID(true)
	if err != nil {
		h.Send("", 500, err)
		return
	} else if !found {
		h.Send("", 404, err)
		return
	}

	if err = article.DeleteByID(); err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send("", 200, err)
}
