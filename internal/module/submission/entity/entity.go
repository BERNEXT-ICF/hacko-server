package entity

import "time"

type SubmitRequest struct {
	UserId       string `validate:"required"`
	AssignmentId string `json:"assignment_id" validate:"required"`
	Link         string `json:"link" validate:"required"`
	DueDate      string `json:"due_date"`
}

type SubmitResponse struct {
	Id          int       `json:"id" db:"id"`
	UserId      string    `json:"creator_assignment_id" db:"creator_assignment_id"`
	Link        string    `json:"link" db:"link"`
	Status      string    `json:"status" db:"status"`
	SubmittedAt time.Time `json:"submitted_at" db:"submitted_at"`
}
