package model

import (
	"servus/core"
	"servus/core/external/database"
)

var call *core.Instance
var IntAdapter = database.Adapter[int]{}
var StringAdapter = database.Adapter[string]{}
var articleAdapter = database.Adapter[Article]{}
var articleCatAdapter = database.Adapter[ArticleCategory]{}
var fileAdapter = database.Adapter[File]{}
var tokenAdapter = database.Adapter[Token]{}
var userAdapter = database.Adapter[User]{}

func Boot(c *core.Instance) {
	call = c
}
