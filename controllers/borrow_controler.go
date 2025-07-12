package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rivaldiheriyan/managementsystem/dto"
	"github.com/rivaldiheriyan/managementsystem/models"
	"github.com/rivaldiheriyan/managementsystem/services"
	"gorm.io/gorm"
)

type BorrowDb struct {
	DB *gorm.DB
}

func (bc *BorrowDb) BorrowBook(c *gin.Context) {
	var input dto.BorrowBookRequest

	// Ambil user ID dari context (dari JWT Middleware)
	userIDStr, exists := c.Get("userUUID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr.(string))

	// Parse Book UUID
	bookUUID, err := uuid.Parse(input.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid book_id"})
		return
	}

	// Cek apakah user sudah pinjam buku ini dan belum dikembalikan
	var existingBorrow models.Borrow
	err = bc.DB.Where("user_id = ? AND book_id = ? AND returned_at IS NULL", userID, bookUUID).First(&existingBorrow).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You already borrowed this book and haven't returned it"})
		return
	}

	// Cek apakah buku tersedia
	var book models.Book
	if err := bc.DB.First(&book, "book_id = ?", bookUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	if book.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book is out of stock"})
		return
	}

	// Kurangi stok buku
	book.Stock -= 1
	bc.DB.Save(&book)

	// Buat catatan peminjaman
	newBorrow := models.Borrow{
		BorrowID:   uuid.New(),
		UserID:     userID,
		BookID:     bookUUID,
		BorrowedAt: time.Now(),
	}

	if err := bc.DB.Create(&newBorrow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to borrow book", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book borrowed successfully",
		"data":    newBorrow,
	})
}


func (bc *BorrowDb) ReturnBook(c *gin.Context) {
	var input dto.BorrowBookRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	userIDStr, exists := c.Get("userUUID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr.(string))

	bookUUID, err := uuid.Parse(input.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid book_id"})
		return
	}

	// Cari pinjaman yang aktif
	var borrow models.Borrow
	err = bc.DB.Where("user_id = ? AND book_id = ? AND returned_at IS NULL", userID, bookUUID).First(&borrow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "No active borrow found for this book"})
		return
	}

	// Tandai sebagai dikembalikan
	now := time.Now()
	borrow.ReturnedAt = &now
	if err := bc.DB.Save(&borrow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to return book", "error": err.Error()})
		return
	}

	// Tambah stok buku
	bc.DB.Model(&models.Book{}).Where("book_id = ?", bookUUID).Update("stock", gorm.Expr("stock + ?", 1))

	c.JSON(http.StatusOK, gin.H{
		"message": "Book returned successfully",
	})
}


func (db *BorrowDb) ListBorrow (c *gin.Context){
	role := c.GetString("role")
	userUUID := c.GetString("userUUID")

	var borrows []models.Borrow
	var err error

	if role == "admin" {
		borrows, err = services.GetAllBorrows(db.DB)
	} else {
		borrows, err = services.GetMyBorrows(db.DB, userUUID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve borrow history",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"borrows": borrows,
	})
}

