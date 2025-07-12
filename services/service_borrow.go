package services

import (
	"github.com/rivaldiheriyan/managementsystem/models"
	"gorm.io/gorm"
)

func GetAllBorrows(db *gorm.DB) ([]models.Borrow, error) {
	var borrows []models.Borrow
	err := db.Preload("User").Preload("Book").Find(&borrows).Error
	return borrows, err
}

func GetMyBorrows(db *gorm.DB, userID string) ([]models.Borrow, error) {
	var borrows []models.Borrow
	err := db.Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&borrows).Error
	return borrows, err
}
