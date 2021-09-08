package core

import (
	"servus/core/database"
	"servus/core/logger"
)

type Servus struct {
	DB *database.DB
	Logger logger.Logger
	Config ConfigFile
	Utils Utils
}
var servus = Servus{}
