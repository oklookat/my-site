package banhammer

// active: if false bypass all future calls.
//
// workPath: dir where banhammer.db will be.
//
// maxWarns: max warns to ban.
func New(active bool, workPath string, maxWarns int) (*Instance, error) {
	var hammer = &Instance{}
	var err = hammer.boot(active, workPath, maxWarns)
	return hammer, err
}

type Instance struct {
	active bool
	db     *SQLite
	*Banner
	*Warner
	*Service
}

func (i *Instance) boot(active bool, workPath string, maxWarns int) error {
	i.active = active

	var err error

	if i.active {
		// database.
		var db = &SQLite{}
		if err = db.Boot(workPath, maxWarns); err != nil {
			return err
		}
		i.db = db
	}

	// service.
	i.Service = &Service{}
	i.Service.New(i)

	// ban.
	var ban = &Banner{}
	ban.New(i)
	i.Banner = ban

	// warn.
	var warn = &Warner{}
	warn.New(i)
	i.Warner = warn

	return err
}
