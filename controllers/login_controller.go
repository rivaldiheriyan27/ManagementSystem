package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rivaldiheriyan/managementsystem/helpers"
	"github.com/rivaldiheriyan/managementsystem/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginDB struct {
	DB *gorm.DB
}

func (db *LoginDB) Register(c *gin.Context) {
	var req models.User

	// Bind dari JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Validasi sederhana
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password is required"})
		return
	}

	// Cek apakah email sudah digunakan
	var existingUser models.User
	err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		// Data ditemukan, email sudah dipakai
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email already used",
		})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Error lain saat query
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
			"error":   err.Error(),
		})
		return
	}

	// Generate UUID jika belum diisi
	if req.UserID == uuid.Nil {
		req.UserID = uuid.New()
	}

	// Buat user baru
	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	// Sukses
	c.JSON(http.StatusCreated, gin.H{
		"message":  "User registered successfully",
		"user_id":  req.UserID,
		"username": req.Username,
		"email":    req.Email,
	})
}


func (db *LoginDB) Login(c *gin.Context) {
	var req models.User

	// Bind data JSON login (email dan password)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Ambil data user berdasarkan email
	var dbResult models.User
	if err := db.DB.Where("email = ?", req.Email).First(&dbResult).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found"})
		return
	}

	// Debug password
	fmt.Println("Password input:", req.Password)
	fmt.Println("Password DB   :", dbResult.Password)

	// Cek password yang dimasukkan cocok dengan hash
	if err := bcrypt.CompareHashAndPassword([]byte(dbResult.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Buat JWT token
	token := helpers.GenerateToken(dbResult.Username)

	// Response berhasil
	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"token":   token,
		"user": gin.H{
			"user_id":  dbResult.UserID,
			"username": dbResult.Username,
			"email":    dbResult.Email,
		},
		"role": dbResult.Role,
	})
}
