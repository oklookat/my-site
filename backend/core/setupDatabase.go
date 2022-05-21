package core

import "servus/core/external/database"

func (i *Instance) setupDatabase() error {
	// connect to DB. Database be available via database.Adapter.
	var conn = database.Connector{}
	return conn.New(i.Config.DB, i.Logger)
}
