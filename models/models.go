package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
	Email    string `form:"unique"`
    IsAdmin  bool `gorm:"default:false"`
}
