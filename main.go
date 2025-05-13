package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/MadManJJ/go-gorm/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v4"
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
	gormdb *gorm.DB
)

func authRequired(c *fiber.Ctx) error {
    // First check for JWT in Authorization header
    tokenStr := c.Get("Authorization")
    if tokenStr == "" {
        // Fallback to cookie if no Authorization header is provided
        tokenStr = c.Cookies("jwt")
    }
    if tokenStr == "" {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    // Strip "Bearer " if it's in the Authorization header
    if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
        tokenStr = tokenStr[7:]
    }

    jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
    token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecretKey), nil
    })

    if err != nil || !token.Valid {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    claim := token.Claims.(jwt.MapClaims)
    fmt.Println(claim["user_id"])

    return c.Next()
}


// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	gormdb = db
	gormdb.AutoMigrate(&Book{}, &User{}) // * AutoMigrate won't delete col, it can only create col

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	
	app.Use("/books", authRequired) // * Middleware

	app.Get("/books", GetBooks)

	app.Get("/books/:id", GetBook)

	app.Post("/books", CreateBook)

	app.Put("/books/:id", UpdateBook)

	app.Delete("/books/:id", DeleteBook)

	app.Post("/register", Register)

	app.Post("/login", LoginUser)

	app.Listen(":8080")

	// * Create Book
	// createBook(db, &Book{
	// 	Name: "suzy",
	// 	Author: "kkk",
	// 	Price: 200,
	// 	Description: "Test",
	// })
	// * --------------------------------

	// * Get Book
	// currentBook := getBook(db,1)
	// fmt.Println(currentBook)
	// * --------------------------------

	// * Update Book
	// currentBook := getBook(db, 1) // getBook return an address

	// currentBook.Name = "BOBA JOHN"
	// currentBook.Price = 440

	// updateBook(db, currentBook)
	// * --------------------------------

	// * Delete Book
	// deleteBook(db, 1)
	// * --------------------------------
	
	// * Search Book
	// currentBook := searchBook(db, "suzy")
	// fmt.Println(currentBook)
	// * --------------------------------
}

// GetBooks godoc
// @Summary Get all books
// @Description Get details of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} BookDTO
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books [get]
func GetBooks(c *fiber.Ctx) error {
	return c.JSON(getBooks(gormdb))
}

// GetBook godoc
// @Summary Get book
// @Description Get book by ID
// @Tags books
// @Produce  json
// @Security ApiKeyAuth
// @Param bookID path int true "Book ID"
// @Success 200 {object} BookDTO
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books/{bookID} [get]
func GetBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	book := getBook(gormdb, id)
	return c.Status(fiber.StatusOK).JSON(book)
}

// CreateBook godoc
// @Summary Create book
// @Description Create book
// @Tags books
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Book body BookDTO true "Book DTO"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books [post]
func CreateBook(c *fiber.Ctx) error {
	book := new(Book) // * book is a pointer
	// var book Book // * book is a regular value

	if err := c.BodyParser(book); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := createBook(gormdb, book)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{
		"message" : "Create Book Successful",
	})
}

// UpdateBook godoc
// @Summary Update book
// @Description Update book
// @Tags books
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param bookID path int true "Book ID"
// @Param Book body BookDTO true "Book DTO"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books/{bookID} [put]
func UpdateBook(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book := new(Book)

		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book.ID = uint(id)

		err = updateBook(gormdb, book)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message" : "Update Book Successful",
		})
	}

// DeleteBook godoc
// @Summary Delete book
// @Description Delete book
// @Tags books
// @Produce  json
// @Security ApiKeyAuth
// @Param bookID path int true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books/{bookID} [delete]
func DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = deleteBook(gormdb, id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{
		"message" : "Delete Book Successful",
	})
}

// RegisterUser godoc
// @Summary User register
// @Description User register
// @Tags auth
// @Accept  json
// @Produce  json
// @Param User body UserDTO true "User DTO"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func Register(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := createUser(gormdb, user)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	
	return c.JSON(fiber.Map{
		"message" : "Register Successful",
	})
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param User body UserDTO true "User DTO"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func LoginUser(c *fiber.Ctx) error {
	var user User

	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := loginUser(gormdb, &user)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// c.Cookie(&fiber.Cookie{
	// 	Name:     "jwt",
	// 	Value:    token,
	// 	Expires:  time.Now().Add(time.Hour * 72),
	// 	HTTPOnly: true,
	// })

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "Login successful",
		"Token" : token,
	})	
}