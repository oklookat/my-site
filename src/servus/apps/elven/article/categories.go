package article

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
)

// ALL HANDLERS PROTECTED BY SAFE METHODS MIDDLEWARE.

// get all categories (GET)
func getCategories(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// get all.
	category := model.ArticleCategory{}
	categories, err := category.GetAll()
	if err != nil {
		h.Send(throw.Server(), 500, err)
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
		h.Send(throw.Server(), 500, err)
		return
	}

	h.Send(string(jsonResponse), 200, err)
}

// get category by name (GET)
func getCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// get id from params.
	var name = h.GetRouteArgs()["name"]

	// find.
	var category = model.ArticleCategory{Name: name}
	found, err := category.FindByName()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotFound(), 404, err)
		return
	}

	// send.
	categoryJson, err := json.Marshal(category)
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	h.Send(string(categoryJson), 200, err)
}

// add category (POST)
func addCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// validate.
	body := &base.CategoryBody{}
	if err := ValidateCategoryBody(body, request.Body); err != nil {
		h.Send("bad request", 400, nil)
		return
	}

	// fill.
	var category = model.ArticleCategory{}
	category.Name = body.Name

	// exists?
	found, err := category.FindByName()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if found {
		h.Send(throw.Exists(), 409, err)
		return
	}

	// create.
	if err = category.Create(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}

	// send created.
	categoryJson, err := json.Marshal(&category)
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	h.Send(string(categoryJson), 200, err)
}

// rename category (PUT/PATCH)
func renameCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// validate.
	body := &base.CategoryBody{}
	if err := ValidateCategoryBody(body, request.Body); err != nil {
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
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotFound(), 404, err)
		return
	}

	// update.
	category.Name = body.Name
	if err = category.ChangeNameByID(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}

	// send.
	categoryJson, err := json.Marshal(&category)
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	h.Send(string(categoryJson), 200, err)
}

// delete category by ID (DELETE)
func deleteCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)
	// get name from params.
	var id = h.GetRouteArgs()["id"]
	// find.
	var category = model.ArticleCategory{}
	category.ID = id
	found, err := category.FindByID()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotFound(), 404, err)
		return
	}
	// delete.
	if err = category.DeleteByID(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}
