package models
//
//import (
//	"time"
//)
//
//type User struct {
//	ID        uint   `gorm:"primaryKey"`
//	Role      string `gorm:"default: user; notNull;"`
//	Username  string `gorm:"check: username > 4; size: 24; unique; notNull;"`
//	Password  string `gorm:"check: password > 8; size: 64; notNull;"`
//	RegIP     string `gorm:"default: unknown;"`
//	RegAgent  string `gorm:"default: unknown;"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	Tokens    []Token   `gorm:"constraint:OnDelete:CASCADE; foreignKey:User; references:ID"`
//	Articles  []Article `gorm:"constraint:OnDelete:CASCADE; foreignKey:User; references:ID"`
//	Files     []File    `gorm:"constraint:OnDelete:CASCADE; foreignKey:User; references:ID"`
//}
