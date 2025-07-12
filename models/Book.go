package models

import (
	"github.com/google/uuid"
)

type Book struct {
	BookID  	uuid.UUID `gorm:"type:uuid;primaryKey" json:"book_id"`
	Title   	string    `json:"title"`
	Description string    `json:"description"`
	Stock   	int       `json:"stock"`

	// Relasi ke Borrow (bukan langsung ke User)
	 Borrows []Borrow `gorm:"-:all"`
}
