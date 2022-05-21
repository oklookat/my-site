package banhammer

import (
	"database/sql"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// one IP file in IP's list.
type IPEntry struct {
	ID         int
	IP         string
	IsBanned   bool
	WarnsCount int
	CreatedAt  int64
}

type SQLite struct {
	path       string
	maxWarns   int
	connection *sql.DB
}

func (s *SQLite) Boot(databasePath string, maxWarns int) error {
	s.maxWarns = maxWarns
	s.SetPath(databasePath)

	// check is database exists.
	var isDbExists, err = s.isDatabaseExists()
	if err != nil {
		return err
	}

	// create database if not exists.
	if !isDbExists {
		if err = s.createDatabase(); err != nil {
			return err
		}
	}

	// connect.
	err = s.connectToDatabase()
	if err == nil && isDbExists {
		err = s.someAmnesty()
	}
	return err
}

func (s *SQLite) AddOrUpdateEntry(ip string, entry IPEntry) error {
	if !IsIpValid(ip) {
		return createError("invalid IP")
	}

	// convert / correct.
	var isBannedInt = 0
	if entry.IsBanned || entry.WarnsCount >= s.maxWarns {
		entry.WarnsCount = s.maxWarns
		entry.IsBanned = true
		isBannedInt = 1
	} else if entry.WarnsCount < 0 {
		entry.WarnsCount = 0
	}

	// check is exists.
	var tempEntry, err = s.GetEntry(ip)
	if err != nil {
		return err
	}
	var query string
	var queryArgs = make([]any, 0)
	if tempEntry != nil {
		// exists - overwrite.
		query = `UPDATE ip_list 
		SET is_banned = $1,
		warns_count = $2
		WHERE ip = $3
		`
		queryArgs = append(queryArgs, isBannedInt, entry.WarnsCount, ip)
	} else {
		// not exists - create.
		query = `INSERT INTO ip_list 
		(ip, is_banned, warns_count, created_at) values ($1, $2, $3, $4)
		`
		var currentTime = time.Now().Unix()
		queryArgs = append(queryArgs, ip, isBannedInt, entry.WarnsCount, currentTime)
	}

	// exec.
	_, err = s.connection.Exec(query, queryArgs...)
	return err
}

func (s *SQLite) GetEntry(ip string) (*IPEntry, error) {
	var err error
	// get.
	var row = s.connection.QueryRow("SELECT * FROM ip_list WHERE ip = $1 LIMIT 1", ip)
	if row == nil {
		return nil, nil
	}
	var entry = IPEntry{}
	var isBannedInt = 0
	if err = row.Scan(&entry.ID, &entry.IP, &isBannedInt, &entry.WarnsCount, &entry.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// convert.
	var isBannedBool = false
	if isBannedInt > 0 {
		isBannedBool = true
	}
	entry.IsBanned = isBannedBool
	return &entry, err
}

func (s *SQLite) RemoveEntry(ip string) error {
	var _, err = s.connection.Exec("DELETE FROM ip_list WHERE ip = $1", ip)
	return err
}

// remove some items depend on some conditions.
func (s *SQLite) someAmnesty() error {
	// remove 5 days-old entries.
	var fiveDaysAgo = time.Now().Add(120 - time.Hour).Unix()
	if _, err := s.connection.Exec("DELETE FROM ip_list WHERE created_at <= $1", fiveDaysAgo); err != nil {
		return err
	}

	// remove not banned & not warned
	if _, err := s.connection.Exec("DELETE FROM ip_list WHERE is_banned < 1 AND warns_count < 1"); err != nil {
		return err
	}

	return nil
}

func (s *SQLite) SetPath(path string) {
	path = filepath.ToSlash(path)
	// path.
	var fullPath = path + "/banhammer.db"
	fullPath = filepath.Clean(fullPath)
	s.path = fullPath
}

func (s *SQLite) GetDatabasePath() string {
	return s.path
}

func (s *SQLite) connectToDatabase() error {
	var err error
	// close previous connection if exists.
	if s.connection != nil {
		_ = s.connection.Close()
	}

	// connect.
	s.connection, err = sql.Open("sqlite3", s.GetDatabasePath())
	if err == nil {
		err = s.connection.Ping()
	}
	return err
}

func (s *SQLite) isDatabaseExists() (bool, error) {
	_, err := os.Stat(s.GetDatabasePath())
	if err == nil {
		return true, err
	}
	var isExists = errors.Is(err, fs.ErrExist)
	var isNotExists = errors.Is(err, fs.ErrNotExist)
	if isExists || isNotExists {
		return isExists, nil
	}
	return false, err
}

func (s *SQLite) createDatabase() error {
	// create database file.
	var err error
	if err = os.WriteFile(s.GetDatabasePath(), nil, 0644); err != nil {
		return err
	}

	// open connection.
	if err = s.connectToDatabase(); err != nil {
		return err
	}
	defer s.connection.Close()

	// exec query.
	var query = `
	CREATE TABLE ip_list (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		ip TEXT UNIQUE,
		is_banned INTEGER,
		warns_count INTEGER,
		created_at INTEGER
	  );
	`
	_, err = s.connection.Exec(query)
	return err
}
