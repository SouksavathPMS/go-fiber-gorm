package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SouksavathPMS/go-fiber-postgres/models"
	"github.com/SouksavathPMS/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

// CreateBook handles /create-books to create the book
func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := &models.Books{}
	if err := context.BodyParser(&book); err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"},
		)
		return err
	}
	if err := r.DB.Create(&book).Error; err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"},
		)
		return err
	}
	context.Status(http.StatusCreated).JSON(
		&fiber.Map{"message": "Book has been added"},
	)
	return nil
}

// Getbooks handles /books to get all the available books
func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	if err := r.DB.Find(bookModels).Error; err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Could not get the books"},
		)
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "books fetched successfully", "data": bookModels},
	)
	return nil
}

// DeleteBook handles /delete-book/:{id} for delete the book
func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Book not found"},
		)
		return nil
	}
	err := r.DB.Delete(bookModel, id)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not delete book"},
		)
		return err.Error
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Book delete successfully"},
	)
	return nil
}

// GetBookByID handles /get-books/:{id} which will get the book info by its ID
func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Book not found"},
		)
		return nil
	}
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get the book with the ID"},
		)
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Fected book with the ID successfully", "data": bookModel},
	)
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create-books", r.CreateBook)
	api.Delete("/delete-books/:id", r.DeleteBook)
	api.Get("/get-books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	// Load our .env file
	err := godotenv.Load(".evn")
	if err != nil {
		log.Fatalf("Error while loading .evn with: %s", err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	// Storage connection
	db, err := storage.NewConnection(config)
	if db != nil {
		log.Fatal(err)
	}

	if err = models.MigrateBooks(db); err != nil {
		log.Fatal("Could not migrate DB")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
