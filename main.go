package main

import (
	"log"
	"net/http"

	"github.com/SouksavathPMS/go-fiber-postgres/models"
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

	// Storage connection
	if db, err := storage.NewConnection(config); err != nil {
		log.Fatal(err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
