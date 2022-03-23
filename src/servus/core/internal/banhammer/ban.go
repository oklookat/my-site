package banhammer

type Banner struct {
	db       *SQLite
	maxWarns int
	// hook.
	onBanned   func(ip string)
	onUnbanned func(ip string)
}

func (b *Banner) New(db *SQLite, maxWarns int) {
	b.db = db
	b.maxWarns = maxWarns
}

func (b *Banner) Ban(ip string) error {
	// get entry.
	var entry, err = b.db.GetEntry(ip)
	if err != nil {
		return err
	}
	if entry == nil {
		entry = &IPEntry{}
	} else if entry.IsBanned {
		err = createError(ip + " already banned")
		return err
	}

	// ban.
	entry.IsBanned = true
	err = b.db.AddOrUpdateEntry(ip, *entry)

	// run hook.
	if err == nil && b.onBanned != nil {
		b.onBanned(ip)
	}

	return err
}

func (b *Banner) OnBanned(hook func(ip string)) {
	b.onBanned = hook
}

func (b *Banner) Unban(ip string) error {
	// get entry.
	var entry, err = b.db.GetEntry(ip)
	if err != nil || entry == nil || !entry.IsBanned {
		err = createError(ip + " not banned")
		return err
	}

	// unban.
	entry.IsBanned = false
	entry.WarnsCount = 0
	err = b.db.AddOrUpdateEntry(ip, *entry)

	// run hook.
	if err == nil && b.onUnbanned != nil {
		b.onUnbanned(ip)
	}

	return err
}

func (b *Banner) OnUnbanned(hook func(ip string)) {
	b.onUnbanned = hook
}

func (b *Banner) IsBanned(ip string) (bool, error) {
	var entry, err = b.db.GetEntry(ip)
	if err != nil || entry == nil {
		return false, err
	}
	return entry.IsBanned, err
}
