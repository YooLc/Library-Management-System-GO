package main

import (
	"library-management-system/database"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	var err error
	DB, err = gorm.Open(mysql.Open(user+":"+password+"@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&database.Book{}, &database.Card{}, &database.Borrow{})

	println("Hello, World!")
}
