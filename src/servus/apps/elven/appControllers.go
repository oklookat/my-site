package elven

// ResponseContent - template for response.
type ResponseContent struct {
	Meta struct {
		PerPage     int `json:"per_page"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

// utils - useful utilities.
var oUtils *objectUtils

// oCmd - commandline methods. Used when app starts.
var oCmd *objectCmd

// eBase - parent for all other entities.
var eBase *entityBase

// eAuth - manage the authorization.
var eAuth *entityAuth

// eArticle - manage the articles.
var eArticle *entityArticle

// eFile - manage the files.
var eFile *entityFile

// eUser - manage the users.
var eUser *entityUser

// eToken - manage the tokens.
var eToken *entityToken

func bootEntities() {
	oUtils = &objectUtils{}
	oCmd = &objectCmd{}
	eBase = &entityBase{instance.HTTP}
	eAuth = &entityAuth{entityBase: eBase}
	eArticle = &entityArticle{entityBase: eBase}
	eFile = &entityFile{entityBase: eBase}
	eUser = &entityUser{entityBase: eBase}
	eToken = &entityToken{}
}
