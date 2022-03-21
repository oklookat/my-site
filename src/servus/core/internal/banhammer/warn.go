package banhammer

type Warner struct {
	db       *SQLite
	banner   *Banner
	maxWarns int
	// hook.
	onWarned func(ip string)
}

func (w *Warner) New(db *SQLite, b *Banner, maxWarns int) {
	w.db = db
	w.banner = b
	if maxWarns < 1 {
		maxWarns = 1
	}
	w.maxWarns = maxWarns
}

func (w *Warner) Warn(ip string) error {
	// get entry.
	var entry, err = w.db.GetEntry(ip)
	if err != nil {
		return err
	}
	if entry == nil {
		entry = &IPEntry{}
	}

	// add warn.
	if entry.WarnsCount < w.maxWarns {
		entry.WarnsCount++
	}

	err = w.db.AddOrUpdateEntry(ip, *entry)

	// run hook.
	if err == nil && w.onWarned != nil {
		w.onWarned(ip)
	}
	return err
}

func (w *Warner) Unwarn(ip string) error {
	// get entry.
	var entry, err = w.db.GetEntry(ip)
	if err != nil || entry == nil || entry.WarnsCount < 1 {
		return err
	}

	// remove warn.
	entry.WarnsCount--
	if entry.WarnsCount < w.maxWarns {
		// unban.
		w.banner.Unban(ip)
	}
	err = w.db.AddOrUpdateEntry(ip, *entry)
	return err
}

func (w *Warner) OnWarned(hook func(ip string)) {
	w.onWarned = hook
}
