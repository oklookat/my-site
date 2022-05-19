package banhammer

type Warner struct {
	hammer *Instance

	// hook.
	onWarned func(ip string)
}

func (w *Warner) New(i *Instance) {
	w.hammer = i
	if !w.hammer.active {
		return
	}

	if w.hammer.db.maxWarns < 1 {
		w.hammer.db.maxWarns = 1
	}
}

func (w *Warner) Warn(ip string) error {
	if !w.hammer.active {
		return nil
	}

	// get entry.
	var entry, err = w.hammer.db.GetEntry(ip)
	if err != nil {
		return err
	}
	if entry == nil {
		entry = &IPEntry{}
	}

	// add warn.
	if entry.WarnsCount < w.hammer.db.maxWarns {
		entry.WarnsCount++
	}

	// notify banner if max warns count reached.
	if entry.WarnsCount >= w.hammer.db.maxWarns {
		entry.WarnsCount = w.hammer.db.maxWarns
		err = w.hammer.Ban(ip)
		if err != nil {
			return err
		}
	}

	err = w.hammer.db.AddOrUpdateEntry(ip, *entry)

	// run hook.
	if err == nil && w.onWarned != nil {
		w.onWarned(ip)
	}
	return err
}

func (w *Warner) Unwarn(ip string) error {
	if !w.hammer.active {
		return nil
	}

	// get entry.
	var entry, err = w.hammer.db.GetEntry(ip)
	if err != nil || entry == nil || entry.WarnsCount < 1 {
		return err
	}

	// remove warn.
	entry.WarnsCount--
	if entry.WarnsCount < w.hammer.db.maxWarns {
		// unban.
		w.hammer.Unban(ip)
	}
	err = w.hammer.db.AddOrUpdateEntry(ip, *entry)
	return err
}

func (w *Warner) OnWarned(hook func(ip string)) {
	w.onWarned = hook
}
