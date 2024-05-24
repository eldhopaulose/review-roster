package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/eldhopaulose/ReviewRoster/config/db"
	"github.com/eldhopaulose/ReviewRoster/initialize"
	"github.com/eldhopaulose/ReviewRoster/models"
)

func main() {
	// Load environment variables
	initialize.LoadEnv()

	// Connect to the database
	db.ConnectDB()
	db.Migrate()

	r := gin.Default()
	r.POST("/", createBook)
	r.GET("/", getAllBooks) // Endpoint to get all books
	r.Run(":8080")                // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func createBook(c *gin.Context) {
	var book models.Books
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := db.GetDB()
	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func getAllBooks(c *gin.Context) {
    var books []models.Books

    db := db.GetDB()
    if err := db.Find(&books).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"books": books})
}
