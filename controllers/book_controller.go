package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rivaldiheriyan/managementsystem/dto"
	"github.com/rivaldiheriyan/managementsystem/models"
	"gorm.io/gorm"
)


type BookDb struct {
	DB *gorm.DB
}

func (db *BookDb) ListBook (c *gin.Context){
	var book []models.Book

	if err := db.DB.Find(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve products",
			"error":   err.Error(),
		})
		return
	}

	// Mapping ke DTO
	var bookResponse []dto.ListBookResponse
	for _, p := range book {
		bookResponse = append(bookResponse, dto.ListBookResponse{
			BookID: p.BookID,
			Title: p.Title,
			Description:p.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List Data Book",
		"book": bookResponse,
	})
}

func (db *BookDb) DetailBook (c *gin.Context){
	bookID := c.Param("book_id");

	var book models.Book

	fmt.Println(bookID,"ini datanya product id")

	uuidBookID, err := uuid.Parse(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID format"})
		return
	}

	err = db.DB.First(&book, "book_id = ?", uuidBookID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Kalau tidak ditemukan
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Book not found",
		})
		return
	} else if err != nil {
		// Kalau ada error database lainnya
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.DetailBookResponse{
		BookID:      book.BookID,
		Title:       book.Title,
		Description: book.Description,
		Stock:       book.Stock,
	})
}

func (db *BookDb) CreateBook (c *gin.Context){
	var book models.Book

	err := c.ShouldBindJSON(&book);
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err);
	}

	book.BookID = uuid.New()

	if err := db.DB.Create(&book).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Title": book.Title,
		"Description" :book.Description,
	})

}

func (db *BookDb) EditBook(c *gin.Context) {
	bookID := c.Param("book_id")
	fmt.Println(bookID, "ini datanya book id")

	uuidBookID, err := uuid.Parse(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID format"})
		return
	}

	// Cek apakah bukunya ada
	var book models.Book
	if err := db.DB.First(&book, "book_id = ?", uuidBookID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error", "error": err.Error()})
		return
	}

	// Ambil input update
	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Stock       *int    `json:"stock"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	// Update hanya field yang diisi
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Description != nil {
		book.Description = *input.Description
	}
	if input.Stock != nil {
		book.Stock = *input.Stock
	}

	// Simpan ke database
	if err := db.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update book", "error": err.Error()})
		return
	}

	// Berhasil
	c.JSON(http.StatusOK, gin.H{
		"message":   "Book updated successfully",
		"book_data": book,
	})
}


func (db *BookDb) DeleteBook (c *gin.Context){
	bookID := c.Param("book_id")

	uuidBookID, err := uuid.Parse(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID format"})
		return
	}

	var book models.Book
	if err := db.DB.First(&book, "book_id = ?", uuidBookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if err := db.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":"Failed to Delete Prodict",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
		"deleted": book.Title,
	})
}

