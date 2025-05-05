package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Book{}) // AutoMigrate will not delete col, it can only create col
	fmt.Println("Migrate succesfull")

	// * Update Book
	// currentBook := getBook(db, 1) // getBook return an address

	// currentBook.Name = "BOBA JOHN"
	// currentBook.Price = 440

	// updateBook(db, currentBook)
	// * --------------------------------

	// * Create Book
	// createBook(db, &Book{
	// 	Name: "suzy",
	// 	Author: "john",
	// 	Price: 400,
	// 	Description: "Test",
	// })
	// * --------------------------------

	// * Get Book
	// currentBook := getBook(db,1)
	// fmt.Println(currentBook)
	// * --------------------------------

	// * Delete Book
	deleteBook(db, 1)
	// * --------------------------------

	// currentBook := getBook(db, 1)
	// fmt.Println(currentBook)
}