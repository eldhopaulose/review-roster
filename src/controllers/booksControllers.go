package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/eldhopaulose/ReviewRoster/src/config/db"
	"github.com/eldhopaulose/ReviewRoster/src/models"
)

// Get all books
func GetAllBooksController(c *gin.Context) {
    var books []models.Books

    db := db.GetDB()
    if err := db.Find(&books).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Books not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		}
		return
    }

    c.JSON(http.StatusOK, gin.H{"books": books})
}
// GetBookController handles GET requests to fetch a book by its ID.
func GetBookController(c *gin.Context) {
	var book models.Books

	db := db.GetDB()
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// Create a new book
func CreateBookController(c *gin.Context) {
	var book models.Books
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := db.GetDB()
	if err := db.Create(&book).Error; err != nil {
		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "dtails": err.Error })
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book", "dtails": err.Error })
		}
		return
	}

	c.JSON(http.StatusOK, book)
}
// UpdateBookController handles PUT requests to update a book by its ID.
func UpdateBookController(c *gin.Context) {
	var book models.Books
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := db.GetDB()
	if err := db.Model(&book).Where("id = ?", c.Param("id")).Updates(book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// DeleteBookController handles DELETE requests to delete a book by its ID.
func DeleteBookController(c *gin.Context) {
	var book models.Books
	db := db.GetDB()
	if err := db.Where("id = ?", c.Param("id")).Delete(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}