package banhammer

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// IP's list.
type IPList struct {
	List map[string]*IPEntry `json:"ip"`
}

// one IP file in IP's list.
type IPEntry struct {
	IsBanned   bool `json:"is_banned"`
	WarnsCount int  `json:"warns_count"`
}

// IP list ops.
type List struct {
	dir  string
	path string
	// cached (RAM) ip list.
	cached *IPList
}

func (l *List) SetPath(path string) {
	path = filepath.ToSlash(path)
	path = filepath.Clean(path)
	// dir.
	dirname, _ := filepath.Split(path)
	l.dir = dirname
	// path.
	var fullPath = dirname + "/banhammer.json"
	fullPath = filepath.Clean(fullPath)
	l.path = fullPath
}

func (l *List) GetPath() string {
	return l.path
}

func (l *List) IsExists() (bool, error) {
	_, err := os.Stat(l.GetPath())
	var isExists = errors.Is(err, fs.ErrExist)
	var isNotExists = errors.Is(err, fs.ErrNotExist)
	if isExists || isNotExists {
		return isExists, nil
	}
	return false, err
}

func (l *List) GetList() (*IPList, error) {
	// cache.
	if l.cached != nil {
		return l.cached, nil
	}

	// check.
	var isExists, err = l.IsExists()
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, createError("ip's list not exists")
	}

	// open.
	file, err := l.openFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode.
	var ipList = IPList{}
	err = json.NewDecoder(file).Decode(&ipList)
	if err != nil {
		return nil, err
	}

	// send.
	return &ipList, err
}

func (l *List) GetEntry(ip string) (*IPEntry, error) {
	var ips, err = l.GetList()
	if err != nil {
		return nil, err
	}
	var entry = ips.List[ip]
	return entry, err
}

func (l *List) AddEntry(ip string, entry IPEntry) error {
	var ips, err = l.GetList()
	if err != nil {
		return err
	}
	if ips == nil {
		return createError("AddEntry: nil IP list")
	}
	*ips.List[ip] = entry
	err = l.WriteList(ips)
	return err
}

func (l *List) RemoveEntry(ip string) error {
	var ips, err = l.GetList()
	if err != nil {
		return err
	}
	if ips == nil {
		return createError("RemoveIP: nil IP list")
	}
	delete(ips.List, ip)
	err = l.WriteList(ips)
	return err
}

func (l *List) Recreate() error {
	return l.writeToFile(nil)
}

func (l *List) WriteList(list *IPList) error {
	if list == nil {
		return createError("WriteList: nil list pointer")
	}
	var err = l.writeToFile(list)
	if err == nil {
		l.cached = list
	}
	return err
}

// -- SERVICE FUNC USE GetList() INSTEAD --
//
// open list file.
func (l *List) openFile() (*os.File, error) {
	var path = l.GetPath()
	var file, err = os.Open(path)
	return file, err
}

// -- SERVICE FUNC USE WriteList() INSTEAD --
//
// write to list file.
//
// if content == nil: empties the list.
//
// if file not exists: create and write.
func (l *List) writeToFile(content *IPList) error {
	var err error

	// recreate if not exists.
	isExists, err := l.IsExists()
	if err != nil {
		return err
	}
	if !isExists || content == nil {
		// create empty list.
		var emptyList = IPList{}
		emptyList.List = make(map[string]*IPEntry, 0)
		content = &emptyList
	}

	// create json from struct.
	jsonBytes, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return err
	}

	// write to file.
	err = os.WriteFile(l.GetPath(), jsonBytes, 0644)
	return err
}
