package banhammer

type Warn struct {
	list     Lister
	ban      Banner
	maxWarns int
	// hook.
	onWarned func(ip string)
}

func (w *Warn) New(l Lister, b Banner, maxWarns int) {
	w.list = l
	w.ban = b
	if maxWarns < 1 {
		maxWarns = 1
	}
	w.maxWarns = maxWarns
}

func (w *Warn) Add(ip string) error {
	// get entry.
	var entry, err = w.list.GetEntry(ip)
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
	if entry.WarnsCount >= w.maxWarns {
		// ban.
		entry.WarnsCount = w.maxWarns
		w.ban.Set(ip)
	}
	err = w.list.AddEntry(ip, *entry)

	// run hook.
	if err == nil && w.onWarned != nil {
		w.onWarned(ip)
	}
	return err
}

func (w *Warn) Remove(ip string) error {
	// get entry.
	var entry, err = w.list.GetEntry(ip)
	if err != nil || entry == nil || entry.WarnsCount < 1 {
		return err
	}

	// remove warn.
	entry.WarnsCount--
	if entry.WarnsCount < w.maxWarns {
		// unban.
		w.ban.Remove(ip)
	}
	err = w.list.AddEntry(ip, *entry)
	return err
}

func (w *Warn) OnWarned(hook func(ip string)) {
	w.onWarned = hook
}
