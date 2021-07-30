package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DBNAME")
	port := os.Getenv("POSTGRES_PORT")
	sslMode := os.Getenv("POSTGRES_SSLMODE")
	timeZone := os.Getenv("POSTGRES_TIMEZONE")
	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbName, port, sslMode, timeZone)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Printf("error init DB: %v", err)
	}
	conn.AutoMigrate(&User{}, &Content{})
	conn.Model(&Content{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	DB = conn
}
