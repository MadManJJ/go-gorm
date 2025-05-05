package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// * gorm.Model definition
// type Model struct {
//   ID        uint           `gorm:"primaryKey"`
//   CreatedAt time.Time
//   UpdatedAt time.Time
//   DeletedAt gorm.DeletedAt `gorm:"index"`
// }

type Book struct {
	gorm.Model
	Name        string
	Author      string
	Description string
	Price       uint
}

func createBook(db *gorm.DB, book *Book) {
	result := db.Create(book)

	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}

	fmt.Println("Create Book Successful")
}

func getBook(db *gorm.DB, id int) *Book {
	var book Book
	result := db.First(&book, id) // * first argument is for storing the book we find, second argument is for finding that primary key

	if result.Error != nil {
		log.Fatalf("Error get book: %v", result.Error)
	}

	return &book
}

func updateBook(db *gorm.DB, book *Book) {
	result := db.Save(book)

	if result.Error != nil {
		log.Fatalf("Update book failed: %v", result.Error)
	}

	fmt.Println("Update Book Successful")
}

func deleteBook(db *gorm.DB, id uint) {
	var book Book
	result := db.Delete(&book, id) // ! soft delete if we have DeletedAt gorm.DeletedAt `gorm:"index"`, but hard delete if we don't
	// result := db.Unscoped().Delete(&book, id) // ! Permanet Delete

	if result.Error != nil {
		log.Fatalf("Delete book failed: %v", result.Error)
	}

	fmt.Println("Delete Book Successful")
}