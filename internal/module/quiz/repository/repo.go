package repository

import (
	"context"
	"database/sql"
	"hacko-app/internal/module/quiz/entity"
	"hacko-app/internal/module/quiz/ports"
	"hacko-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.QuizRepository = &quizRepository{}

type quizRepository struct {
	db *sqlx.DB
}

func NewQuizRepository(db *sqlx.DB) *quizRepository {
	return &quizRepository{
		db: db,
	}
}

func (r *quizRepository) FindClass(ctx context.Context, req string) error {
	query := `SELECT id FROM class WHERE id = $1`

	var classId string
	err := r.db.QueryRowContext(ctx, query, req).Scan(&classId)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Str("payload", req).Msg("repo::FindClass - Class not found")
			return errmsg.NewCustomErrors(404, errmsg.WithMessage("Class not found"))
		}
		log.Error().Err(err).Str("payload", req).Msg("repo::FindClass - Failed to query class")
		return err
	}

	return nil
}

func (r *quizRepository) CreateQuiz(ctx context.Context, req *entity.CreateQuizRequest) (*entity.CreateQuizResponse, error) {
	query := `INSERT INTO quiz (class_id, creator_quiz_id, title, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		RETURNING id, class_id, creator_quiz_id, title, status, created_at, updated_at`

	var res entity.CreateQuizResponse
	err := r.db.QueryRowContext(ctx, query, req.ClassId, req.UserId, req.Title, req.Status).Scan(&res.Id, &res.ClassId, &res.CreatorQuizId, &res.Title, &res.Status, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateQuiz - Failed Create Quiz")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}

	return &res, nil
}

func (r *quizRepository) FindQuiz(ctx context.Context, req int) error{
    query := `SELECT id FROM quiz WHERE id = $1`

	var classId string
	err := r.db.QueryRowContext(ctx, query, req).Scan(&classId)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("payload", req).Msg("repo::FindQuiz - Quiz not found")
			return errmsg.NewCustomErrors(404, errmsg.WithMessage("Quiz not found"))
		}
		log.Error().Err(err).Any("payload", req).Msg("repo::FindClass - Failed to query quiz")
		return err
	}

	return nil
}

func (r *quizRepository) CreateQuestionQuiz(ctx context.Context, req *entity.CreateQuestionQuizRequest) (*entity.CreateQuestionQuizResponse, error) {
    query := `
        INSERT INTO questions_quiz (quiz_id, creator_question_quiz_id, type, question, answers, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, creator_question_quiz_id, type, question, answers, created_at, updated_at;
    `

    var resp entity.CreateQuestionQuizResponse
    err := r.db.QueryRowContext(ctx, query,
        req.QuizId,                  // quiz_id
        req.UserId,                  // creator_quiz_id (user ID)
        req.Type,                    // type
        req.Question,                // question
        req.Answers,                 // answers (JSON)
    ).Scan(
        &resp.Id,
        &resp.CreatorQuestionQuizId,
        &resp.Type,
        &resp.Question,
        &resp.Answers,
        &resp.CreatedAt,
        &resp.UpdatedAt,
    )

    if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateQuestionQuiz - Failed create question quiz")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}

    return &resp, nil
}
