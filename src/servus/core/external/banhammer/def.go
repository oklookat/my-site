package banhammer

import (
	"net/http"
)

// ban ops.
type Banner interface {
	// ban.
	Ban(ip string) error
	// unban.
	Unban(ip string) error
	// is IP banned?
	IsBanned(ip string) error
}

// warn ops.
type Warner interface {
	// add warn. 3 warns = ban.
	Warn(ip string) error
	// remove warn.
	Unwarn(ip string) error
	// get warns count.
	GetWarnsCount(ip string) (int, error)
}

// service funcs.
type Servicer interface {
	// remove IP from all lists.
	Amnesty(ip string) error
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
	// write list to IP's file.
	WriteList(list *IPList) error
	// (re)create IP's list file.
	Recreate() error
}

type Hammerer interface {
	// create new instance.
	//
	// listPath: IP's list file.
	New(listPath string)
	Banner
	Warner
	Lister
	Servicer
}
