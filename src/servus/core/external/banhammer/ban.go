package banhammer

type Ban struct {
	list Lister
	// hook.
	onBanned func(ip string)
}

func (b *Ban) New(l Lister) {
	b.list = l
}

func (b *Ban) Set(ip string) error {
	// get entry.
	var entry, err = b.list.GetEntry(ip)
	if err != nil {
		return err
	}
	if entry == nil {
		entry = &IPEntry{}
	}
	if entry.IsBanned {
		return err
	}

	// ban.
	entry.IsBanned = true
	err = b.list.AddEntry(ip, *entry)

	// run hook.
	if err == nil && b.onBanned != nil {
		b.onBanned(ip)
	}
	return err
}

func (b *Ban) Remove(ip string) error {
	// get entry.
	var entry, err = b.list.GetEntry(ip)
	if err != nil || entry == nil || !entry.IsBanned {
		return err
	}

	// unban.
	entry.IsBanned = false
	err = b.list.AddEntry(ip, *entry)
	return err
}

func (b *Ban) IsBanned(ip string) (bool, error) {
	var entry, err = b.list.GetEntry(ip)
	if err != nil || entry == nil {
		return false, err
	}
	return entry.IsBanned, err
}

func (b *Ban) OnBanned(hook func(ip string)) {
	b.onBanned = hook
}
