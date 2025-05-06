package main

import (
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
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint `json:"price"`
}

func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getBook(db *gorm.DB, id int) *Book {
	var book Book
	result := db.First(&book, id) // * first argument is for storing the book we find, second argument is for finding that primary key

	if result.Error != nil {
		log.Fatalf("Error get book: %v", result.Error)
	}

	return &book
}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)

	if result.Error != nil {
		log.Fatalf("Error get books: %v", result.Error)
	}

	return books
}

func updateBook(db *gorm.DB, book *Book) error {	
	result := db.Model(&book).Updates(book) // * update only the selected field (from the book)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func deleteBook(db *gorm.DB, id int) error {
	var book Book
	result := db.Delete(&book, id) // ! soft delete if we have DeletedAt gorm.DeletedAt `gorm:"index"`, but hard delete if we don't
	// result := db.Unscoped().Delete(&book, id) // ! Permanet Delete even if we have DeletedAt gorm.DeletedAt `gorm:"index"`

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func searchBook(db *gorm.DB, bookName string) []Book { // * slice normally is already an address
	var books []Book

	result := db.Where("name = ?", bookName).Order("price desc").Find(&books) // * pass a pointer so that gorm can modify the books and fill it.

	if result.Error != nil {
		log.Fatalf("Search book failed: %v", result.Error)
	}

	return books
}