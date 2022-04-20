package core

import "servus/core/internal/banhammer"

func (i *Instance) setupBanhammer() error {
	var err error

	var cfg = i.Config.Security.Banhammer
	if !cfg.Active {
		i.Banhammer, _ = banhammer.New(false, "", 0)
		return nil
	}

	var dbPath = cfg.Database
	if dbPath == nil {
		dbPath = new(string)

		// get dir.
		*dbPath, err = i.Dirs.GetData()
		if err != nil {
			return err
		}
	}

	var maxWarns = cfg.MaxWarns

	// boot banhammer.
	i.Banhammer, err = banhammer.New(true, *dbPath, maxWarns)
	if err != nil {
		return err
	}

	return err
}
