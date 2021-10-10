package elven

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

// ModelUser - represents user in database.
type ModelUser struct {
	ID        string    `json:"id" db:"id"`
	Role      string    `json:"role" db:"role"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	RegIP     *string   `json:"reg_ip" db:"reg_ip"`
	RegAgent  *string   `json:"reg_agent" db:"reg_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ModelToken - represents token in database.
type ModelToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	LastIP    *string   `json:"last_ip" db:"last_ip"`
	LastAgent *string   `json:"last_agent" db:"last_agent"`
	AuthIP    *string   `json:"auth_ip" db:"auth_ip"`
	AuthAgent *string   `json:"auth_agent" db:"auth_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ModelArticle - represents article in database.
type ModelArticle struct {
	ID          string         `json:"id" db:"id"`
	UserID      string         `json:"user_id" db:"user_id"`
	IsPublished bool           `json:"is_published" db:"is_published"`
	Title       string         `json:"title" db:"title"`
	Content     ArticleContent `json:"content" db:"content"`
	Slug        string         `json:"slug" db:"slug"`
	PublishedAt *time.Time     `json:"published_at" db:"published_at"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// ArticleContent - represents ModelArticle content in database.
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
