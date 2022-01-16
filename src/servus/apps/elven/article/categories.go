package article

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"

	"github.com/gorilla/mux"
)

// get all categories (GET)
func (a *Instance) getCategories(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get
	category := model.ArticleCategory{}
	categories, err := category.GetAll()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	// fill info
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = 1
	responseContent.Meta.TotalPages = 1
	responseContent.Meta.PerPage = 1
	responseContent.Data = categories
	// make & send json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// get one category (GET)
func (a *Instance) getCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get name from params.
	var params = mux.Vars(request)
	var name = params["name"]
	// find.
	var category = model.ArticleCategory{Name: name}
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
	body := CategoryBody{}
	validator := body.Validate(request.Body)
	if validator.HasErrors() {
		h.Send(validator.GetJSON(), 400, nil)
		return
	}
	// fill.
	var category = model.ArticleCategory{}
	category.Name = body.Name
	// create.
	err := category.Create()
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
	body := CategoryBody{}
	validator := body.Validate(request.Body)
	if validator.HasErrors() {
		h.Send(validator.GetJSON(), 400, nil)
		return
	}
	// get name from params.
	var params = mux.Vars(request)
	var name = params["name"]
	// find.
	var category = model.ArticleCategory{}
	category.Name = name
	found, err := category.FindByName()
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

// delete category (DELETE)
func (a *Instance) deleteCategory(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get name from params.
	var params = mux.Vars(request)
	var name = params["name"]
	// find.
	var category = model.ArticleCategory{}
	category.Name = name
	found, err := category.FindByName()
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
