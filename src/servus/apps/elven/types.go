package elven

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type ResponseContent struct {
	Meta struct {
		PerPage int    `json:"per_page"`
		Next    string `json:"next"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

type JSON json.RawMessage

func (j *JSON) Scan(value interface{}) error {
	var str, ok = value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}
	bytes := []byte(str)
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	result, err := json.RawMessage(j).MarshalJSON()
	return result, err
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSON: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
