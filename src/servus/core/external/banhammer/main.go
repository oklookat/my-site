package banhammer

type Instance struct {
	Banner
	Warner
	Lister
	Servicer
}

func (i *Instance) New(listPath string, maxWarns int) {
	// service.
	i.Servicer = &Service{}

	// list.
	var list = &List{}
	list.SetPath(listPath)
	i.Lister = list

	// ban.
	var ban = &Ban{}
	ban.New(i.Lister)

	// warn.
	var warn = &Warn{}
	warn.New(i.Lister, i.Banner, maxWarns)
	i.Warner = warn
}
