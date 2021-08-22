package _core

import (
	"servus/_core/database"
	"servus/_core/logger"
)

type Servus struct {
	DB *database.DB
	Logger logger.Logger
	Config ConfigFile
	Utils Utils
}
var servus = Servus{}
