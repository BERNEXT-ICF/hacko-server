package entity

import "time"

type CreateClassRequest struct {
	UserId      string `validate:"required"`
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

type GetClassResponse struct {
	ID             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description,omitempty" db:"description"`
	Image          string    `json:"image,omitempty" db:"image"`
	Video          string    `json:"video,omitempty" db:"video"`
	Status         string    `json:"status" db:"status"`
	CreatorClassID string    `json:"creator_class_id" db:"creator_class_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
type GetAllClassesResponse struct {
	Classes []*GetClassResponse `json:"classes"`
	Total   int                `json:"total"`
}

type GetClassByIdRequest struct {
	Id string `json:"id"`
}

type EnrollClassRequest struct {
	UserId      string `validate:"required"`
	ClassId		int `json:"class_id"`
}