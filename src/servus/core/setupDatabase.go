package core

import "servus/core/external/database"

func (i *Instance) setupDatabase() {
	// connect to DB. Database be available via database.Adapter.
	var conn = database.Connector{}
	conn.New(i.Config.DB, i.Logger)
}
