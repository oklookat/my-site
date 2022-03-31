package banhammer

type Instance struct {
	db *SQLite
	*Banner
	*Warner
	*Service
}

func (i *Instance) Boot(workPath string, maxWarns int) error {
	var err error

	// database.
	var db = &SQLite{}
	if err = db.Boot(workPath, maxWarns); err != nil {
		return err
	}
	i.db = db

	// service.
	i.Service = &Service{}
	i.Service.New(i.db)

	// ban.
	var ban = &Banner{}
	ban.New(i.db, maxWarns)
	i.Banner = ban

	// warn.
	var warn = &Warner{}
	warn.New(i.db, i.Banner, maxWarns)
	i.Warner = warn

	return err
}