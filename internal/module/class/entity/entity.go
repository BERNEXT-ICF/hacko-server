package entity

import "time"

type CreateClassRequest struct {
	UserId      string `json:"user_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Video       string `json:"video"`
	Status      string `json:"status" validate:"required,oneof=public draf"`
}

type CreateClassResponse struct {
	Id        int    `json:"id"`
	CreatorId string `json:"creator_id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetAllClassResponse struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description,omitempty"`
	Image          string    `json:"image,omitempty"`
	Video          string    `json:"video,omitempty"`
	Status         string    `json:"status"`
	CreatorClassID string    `json:"creator_class_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetAllClassesResponse represents the response for fetching multiple classes.
type GetAllClassesResponse struct {
	Classes []GetAllClassResponse `json:"classes"`
	Total   int                   `json:"total"`
}
