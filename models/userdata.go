package models

import (
	_ "github.com/go-sql-driver/mysql"
)
// Entity
type UserData struct {
	ID       int    `gorm:"PRIMARY_KEY;" json:"id"`
	Username string `gorm:"NOT_NULL" json:"username"`
	Email    string `gorm:"NOT_NULL" json:"email"`
}
