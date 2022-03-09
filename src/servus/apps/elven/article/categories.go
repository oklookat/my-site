package article

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
)

// ALL HANDLERS PROTECTED BY SAFE METHODS MIDDLEWARE.

// get all categories (GET)
func (a *Instance) getCategories(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// get all.
	category := model.ArticleCategory{}
	categories, err := category.GetAll()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	// fill response.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = 1
	responseContent.Meta.TotalPages = 1
	responseContent.Meta.PerPage = 1
	responseContent.Data = categories

	// send.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	h.Send(string(jsonResponse), 200, err)
}

// get category by ID (GET)
func (a *Instance) getCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var category = model.ArticleCategory{ID: id}
	found, err := category.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}

	// send.
	categoryJson, err := json.Marshal(category)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(categoryJson), 200, err)
}

// add category (POST)
func (a *Instance) addCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// validate.
	body := &base.CategoryBody{}
	err := ValidateCategoryBody(body, request.Body)
	if err != nil {
		h.Send("bad request", 400, nil)
		return
	}

	// fill.
	var category = model.ArticleCategory{}
	category.Name = body.Name

	// exists?
	found, err := category.FindByName()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if found {
		h.Send(a.throw.Exists(), 409, err)
		return
	}

	// create.
	err = category.Create()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	// send created.
	categoryJson, err := json.Marshal(&category)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(categoryJson), 200, err)
}

// rename category (PUT/PATCH)
func (a *Instance) renameCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// validate.
	body := &base.CategoryBody{}
	err := ValidateCategoryBody(body, request.Body)
	if err != nil {
		h.Send("bad request", 400, nil)
		return
	}

	// get name from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var category = model.ArticleCategory{}
	category.ID = id
	found, err := category.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}

	// update.
	category.Name = body.Name
	err = category.ChangeNameByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	// send.
	categoryJson, err := json.Marshal(&category)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(categoryJson), 200, err)
}

// delete category by ID (DELETE)
func (a *Instance) deleteCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get name from params.
	var id = h.GetRouteArgs()["id"]
	// find.
	var category = model.ArticleCategory{}
	category.ID = id
	found, err := category.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}
	// delete.
	err = category.DeleteByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}
