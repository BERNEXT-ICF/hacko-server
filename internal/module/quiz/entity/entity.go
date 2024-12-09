package entity

import (
	"encoding/json"
	"time"
)

type CreateQuizRequest struct {
	UserId  string `validate:"required"`
	ClassId string `json:"creator_quiz_id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Status  string `json:"status" validate:"required"`
}

type CreateQuizResponse struct {
	Id            int       `json:"id" db:"id"`
	CreatorQuizId string    `json:"creator_quiz_id" db:"creator_quiz_id"`
	ClassId       string    `json:"class_id" db:"class_id"`
	Title         string    `json:"title" db:"title"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateQuestionQuizRequest struct {
	UserId   string          `validate:"required"`
	QuizId   int             `json:"quiz_id" validate:"required"`
	Type     string          `json:"type" validate:"required"`
	Question string          `json:"question" validate:"required"`
	Answers  json.RawMessage `json:"answers" validate:"required"`
}

type CreateQuestionQuizResponse struct {
	Id                    int             `json:"id" db:"id"`
	CreatorQuestionQuizId string          `json:"creator_question_quiz_id" db:"creator_quiz_id"`
	Type                  string          `json:"type" db:"type"`
	Question              string          `json:"question" db:"question"`
	Answers               json.RawMessage `json:"answers" db:"answers"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
}

type GetAllQuizRequest struct {
	UserId  string `validate:"required"`
	ClassId string `json:"class_id" validate:"required"`
}

type GetAllQuizResponse struct {
	Id        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Status    string `json:"status" db:"status"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type GetDetailsQuizRequset struct {
	QuizId int `json:"quiz_id" validate:"required"`
}

type GetQuestionQuizResponse struct {
	Id        string          `json:"id" db:"id"`
	Type      string          `json:"type" db:"type"`
	Question  string          `json:"question" db:"question"`
	Answers   json.RawMessage `json:"answers" db:"answers"`
	CreatedAt string          `json:"created_at" db:"created_at"`
	UpdatedAt string          `json:"updated_at" db:"updated_at"`
}

type GetDetailsQuizResponse struct {
	Id       string                    `json:"id" db:"id"`
	Title    string                    `json:"title" db:"title"`
	Question []GetQuestionQuizResponse `json:"question"`
}
