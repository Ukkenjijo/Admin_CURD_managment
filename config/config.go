package config

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    var err error
    dsn := "host=localhost user=postgres password=root dbname=app port=5432 sslmode=disable TimeZone=UTC"
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }
}
