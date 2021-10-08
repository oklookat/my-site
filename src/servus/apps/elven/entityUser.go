package elven

import (
	"time"
)

// user - user entity.
type entityUser struct {
}

// ModelUser - represents user in database.
type ModelUser struct {
	ID        string `json:"id" db:"id"`
	Role      string `json:"role" db:"role"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	RegIP     *string `json:"reg_ip" db:"reg_ip"`
	RegAgent  *string `json:"reg_agent" db:"reg_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
