package dto

import "github.com/google/uuid"

type ListBookResponse struct{
	BookID   	uuid.UUID `json:"book_id"`
	Title   	string    `json:"title"`
	Description string    `json:"description"`
}

type DetailBookResponse struct {
	BookID     uuid.UUID `json:"book_id"`
	Title      string    `json:"title"`
	Description string   `json:"description"`
	Stock      int       `json:"stock"`
}

type BorrowBookRequest struct {
	BookID string `json:"book_id" binding:"required,uuid"`
}