package elven

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)


// ResponseContent template for response.
type ResponseContent struct {
	Meta struct {
		PerPage int    `json:"per_page"`
		Next    string `json:"next"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

type ArticleContent struct {
	Time   int64 `json:"time"`
	Blocks []struct {
		ID   string      `json:"id"`
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	} `json:"blocks"`
	Version string `json:"version"`
}

func (a ArticleContent) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ArticleContent) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("articleContent: failed convert value to []byte")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, &a)
}
