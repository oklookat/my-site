package banhammer

import (
	"net/http"
)

// ban ops.
type Banner interface {
	// ban.
	Set(ip string) error
	// unban.
	Remove(ip string) error
	// is IP banned?
	IsBanned(ip string) (bool, error)
	// when IP banned.
	OnBanned(hook func(ip string))
}

// warn ops.
type Warner interface {
	// add warn. 3 warns = ban.
	Add(ip string) error
	// remove warn.
	Remove(ip string) error
	// when IP warned.
	OnWarned(hook func(ip string))
}

// service funcs.
type Servicer interface {
	// get ban checking middleware.
	GetMiddleware() func(http.Handler) http.Handler
}

// ip's list ops.
type Lister interface {
	// set dir for IP's file.
	SetPath(path string)
	// get path to IP's file.
	GetPath() string
	// is IP's file exists?
	IsExists() (bool, error)
	// get IP list.
	GetList() (*IPList, error)
	// get IP list entry by IP. If entry not exists returns nil.
	GetEntry(ip string) (*IPEntry, error)
	// add IP with entry to list. If entry exists - overwrites it.
	AddEntry(ip string, entry IPEntry) error
	// remove IP from list.
	RemoveEntry(ip string) error
	// write list to IP's file.
	WriteList(list *IPList) error
	// (re)create IP's list file.
	Recreate() error
}

type Hammerer interface {
	// create new instance.
	//
	// listPath: IP's list file.
	New(listPath string, maxWarns int)
	Banner
	Warner
	Lister
	Servicer
}
