package elven

import (
	"net/http"
	"servus/core/modules/errorMan"
	"strings"
	"time"
)

// entityFile - manage files.
type entityFile struct {
	*entityBase
}

type ModelFile struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Hash         string    `json:"hash" db:"hash"`
	Path         string    `json:"path" db:"path"`
	Name         string    `json:"name" db:"name"`
	OriginalName string    `json:"original_name" db:"original_name"`
	Extension    string    `json:"extension" db:"extension"`
	Size         int64     `json:"size" db:"size"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type queryFileControllerGetAll struct {
	cursor string
	start  string
}

// GET url/
// params:
// cursor = ULID
// start = newest (DESC), oldest (ASC)
func (f *entityFile) controllerGetAll(response http.ResponseWriter, request *http.Request) {

}

// POST url/
func (f *entityFile) controllerCreateOne(response http.ResponseWriter, request *http.Request) {

}

// DELETE url/id
func (f *entityFile) controllerDeleteOne(response http.ResponseWriter, request *http.Request) {

}

func (f *entityFile) validatorControllerGetAll(request *http.Request, ec *errorMan.ErrorCollector, isAdmin bool) (queryControllerGetAll queryFileControllerGetAll, err error) {
	queryControllerGetAll = queryFileControllerGetAll{}
	var queryParams = request.URL.Query()
	// validate start param.
	var start = queryParams.Get("start")
	if len(start) == 0 {
		start = "newest"
	}
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isStartInvalid {
		ec.AddEValidationAllowed([]string{"start"}, []string{"newest", "oldest"})
	} else {
		if isNewest {
			start = "DESC"
		} else if isOldest {
			start = "ASC"
		}
	}
	queryControllerGetAll.start = start
	// validate cursor param.
	var cursor = queryParams.Get("cursor")
	if len(cursor) == 0 {
		cursor = "0"
	}
	queryControllerGetAll.cursor = cursor
	return
}
