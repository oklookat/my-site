package banhammer

type Hammer interface {
	// (re)create banned ip's file.
	RecreateTXT() error
	// add ip to banned list.
	Add(ip string) error
	// remove ip from banned list.
	Remove(ip string) error
}
