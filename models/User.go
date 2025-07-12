package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type User struct {
	UserID    	uuid.UUID 	`gorm:"type:uuid;primaryKey" json:"user_id"`
	Username  	string    	`json:"username"`
	Email     	string    	`json:"email"`
	Password  	string    	`json:"password"`
	Role 		string 		`json:"role"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	

	// Relasi ke Borrow
	Borrows 	[]Borrow 	`gorm:"-:all"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error{
	// Generate UUID
	u.UserID = uuid.New()

	// GenerateFromPassword di bcyrpt
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	
	if err != nil {
		fmt.Println("Failed to encrypt password: ", err)
		return err
	}

	u.Password = string(pwd)
	return nil
}