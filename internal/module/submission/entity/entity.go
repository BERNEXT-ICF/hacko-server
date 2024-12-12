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

type GetSubmissionDetailsRequest struct {
	UserId       string `json:"user_id" validate:"required"`
	SubmissionId string `json:"submission_id" validate:"required"`
}

type GetSubmissionDetailsResponse struct {
	Id           int       `json:"id" db:"id"`
	SubmissionId string    `json:"submission_id" db:"submission_id"`
	Name         string    `json:"name" db:"name"`
	Image        *string    `json:"image_url" db:"image_url"`
	Link         string    `json:"link" db:"link"`
	Status       string    `json:"status" db:"status"`
	Grade        *string    `json:"grade" db:"grade"`
	Feedback     *string    `json:"feedback" db:"feedback"`
	SubmittedAt  time.Time `json:"submitted_at" db:"submitted_at"`
	GradedAt     *time.Time `json:"graded_at" db:"graded_at"`
}

// type GradingSubmission struct {
// 	UserId       string `validate:"required"`
// 	AssignmentId string `json:"assignment_id" validate:"required"`
// 	Link         string `json:"link" validate:"required"`
// 	DueDate      string `json:"due_date"`
// }
