package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		// fmt.Println("Failed3")
		log.Fatal("Error loading .env file")
	}
	host = "localhost"
	databaseName = os.Getenv("POSTGRES_DB")
	username = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
}

var (
	host         string
	port         = 5432
	databaseName string
	username     string
	password     string
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
  "password=%s dbname=%s sslmode=disable",
  host, port, username, password, databaseName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println(db)
	fmt.Println("Connect succesfull")
}