package entity

import "time"

type CreateAssignmentRequest struct {
	UserId      string `validate:"required"`
	ClassId     int    `json:"class_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	DueDate     string `json:"due_date"`
}

type CreateAssignmentResponse struct {
	Id          int       `json:"id" db:"id"`
	UserId      string    `json:"creator_assignment_id" db:"creator_assignment_id"`
	ClassId     int       `json:"class_id" db:"class_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type GetAllAssignmentByClassIdRequest struct {
	UserId  string `validate:"required"`
	ClassId string `json:"class_id"`
}

type GetAssignmentByClassIdResponse struct {
	Id                  int       `json:"id" db:"id"`
	CreatorAssignmentId string    `json:"creator_assignment_id" db:"creator_assignment_id"`
	ClassId             int       `json:"class_id" db:"class_id"`
	Title               string    `json:"title" db:"title"`
	Description         string    `json:"description" db:"description"`
	Status              string    `json:"status"`
	DueDate             time.Time `json:"due_date" db:"due_date"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type GetAssignmentStatusRequest struct {
	UserId       string
	ClassId      string
	AssignmentId int
}

type GetAssignmentDetailsRequest struct {
	UserId       string `validate:"required"`
	AssignmentId int `json:"assignment_id"`
}

type GetAssignmentDetailsResponse struct {
	Id             int    `json:"id" db:"title"`
	Title          string `json:"title" db:"title"`
	Description    string `json:"description" db:"description"`
	DueDate        string `json:"due_date" db:"due_date"`
	LinkSubmission *string `json:"link_submission" db:"link"`
	Grade          *string `json:"grade_subission" db:"grade"`
	Feedback       *string `json:"feedback_submission" db:"feedback"`
	Status         *string `json:"status_submission" db:"status"`
	SubmittedAt    *string `json:"submitted_at" db:"submitted_at"`
}
