package model

import (
	"errors"
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

func Boot(core *core.Instance) error {
	if core == nil {
		return errors.New("core nil pointer")
	}
	call = core
	return nil
}
