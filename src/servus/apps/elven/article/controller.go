package article

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"

	"github.com/gorilla/mux"
)

// getAll - GET url/. Get paginated articles.
func (a *Instance) getAll(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	var err error
	var isAdmin = false
	var pipe = a.pipe.GetByContext(request)
	isAdmin = pipe != nil && pipe.IsAdmin()
	// validate.
	validatedBody := &Paginate{}
	validator := validatedBody.Validate(request.URL.Query(), isAdmin)
	if validator.HasErrors() {
		h.Send(validator.GetJSON(), 400, err)
		return
	}
	// get.
	pag := model.Article{}
	articles, totalPages, err := pag.GetPaginated(validatedBody.By, validatedBody.Start, validatedBody.Show, validatedBody.Page, validatedBody.Preview)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	// fill pagination.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = validatedBody.Page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = model.ArticlePageSize
	responseContent.Data = articles
	// make & send json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// getOne - GET url/id. Get one article.
func (a *Instance) getOne(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	var isAdmin = false
	// get id from params.
	var params = mux.Vars(request)
	var id = params["id"]
	// find.
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}
	// if not published and user not admin, access denied.
	var pAuth = a.pipe.GetByContext(request)
	isAdmin = pAuth.IsAdmin()
	if !article.IsPublished && !isAdmin {
		h.Send(a.throw.Forbidden(), 403, err)
		return
	}
	// send.
	articleJson, err := json.Marshal(article)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// create - POST url/. Creates new article.
func (a *Instance) create(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// validate.
	body := Body{}
	validator := body.Validate(request.Body)
	if validator.HasErrors() {
		h.Send(validator.GetJSON(), 400, nil)
		return
	}
	// fill.
	var pAuth = a.pipe.GetByContext(request)
	var article = model.Article{UserID: pAuth.GetID(), IsPublished: false, Title: *body.Title, Content: *body.Content}
	// create.
	err := article.Create()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	// send created.
	articleJson, err := json.Marshal(&article)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// update - PUT (update all available fields) or PATCH (update specific field) url/id.
func (a *Instance) update(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get id from params.
	var params = mux.Vars(request)
	var id = params["id"]
	// find.
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}
	// validate body.
	var body = Body{}
	validator := body.Validate(request.Body)
	if err != nil {
		h.Send(validator.GetJSON(), 400, err)
		return
	}
	// check fields.
	var isTitle = body.Title != nil
	var isContent = body.Content != nil
	var isPublished = body.IsPublished != nil
	// choose modification by method.
	// if PUT method we need full article to update. If PATCH - we need at least one field.
	if request.Method == http.MethodPut {
		// need full.
		if !isTitle || !isContent || !isPublished {
			validator.Add("body")
			h.Send(validator.GetJSON(), 400, err)
			return
		} else {
			article.Title = *body.Title
			article.Content = *body.Content
			article.IsPublished = *body.IsPublished
		}
	} else if request.Method == http.MethodPatch {
		// need at least one field.
		if isTitle {
			article.Title = *body.Title
		}
		if isContent {
			article.Content = *body.Content
		}
		if isPublished {
			article.IsPublished = *body.IsPublished
		}
	}
	// update.
	err = article.Update()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	// send updated.
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(jsonArticle), 200, err)
}

// delete - DELETE url/id. Deletes one article.
func (a *Instance) delete(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get id from params.
	var params = mux.Vars(request)
	var id = params["id"]
	// find.
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}
	// delete.
	err = article.DeleteByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}
