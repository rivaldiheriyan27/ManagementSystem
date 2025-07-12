package models

import (
	"time"

	"github.com/google/uuid"
)


type Borrow struct {
	BorrowID   uuid.UUID  `gorm:"type:uuid;primaryKey" json:"borrow_id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	BookID     uuid.UUID  `gorm:"type:uuid;not null" json:"book_id"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`

	User User `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Book Book `gorm:"foreignKey:BookID;references:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
