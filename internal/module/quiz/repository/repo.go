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

func (r *quizRepository) FindQuiz(ctx context.Context, req int) error {
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
		req.QuizId,   // quiz_id
		req.UserId,   // creator_quiz_id (user ID)
		req.Type,     // type
		req.Question, // question
		req.Answers,  // answers (JSON)
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

func (r *quizRepository) GetAllQuiz(ctx context.Context, req *entity.GetAllQuizRequest) ([]entity.GetAllQuizResponse, error) {
	query := `
		SELECT id, title, status, created_at, updated_at
		FROM quiz
		WHERE class_id = $1
	`

	var quizzes []entity.GetAllQuizResponse
	err := r.db.SelectContext(ctx, &quizzes, query, req.ClassId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllQuiz - Failed Get All Quiz")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}

	return quizzes, nil
}

func (r *quizRepository) GetDetailsQuiz(ctx context.Context, req *entity.GetDetailsQuizRequset) (*entity.GetDetailsQuizResponse, error) {
	quizQuery := `
		SELECT id, title
		FROM quiz
		WHERE id = $1
	`
	var quiz entity.GetDetailsQuizResponse
	err := r.db.GetContext(ctx, &quiz, quizQuery, req.QuizId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("payload", req).Msg("repo::GetAllQuiz - Quiz with that id not found")
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Quiz with that id not found"))
		}
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllQuiz - Quiz with that id not found")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}

	questionsQuery := `
		SELECT id, type, question, answers, created_at, updated_at
		FROM questions_quiz
		WHERE quiz_id = $1
	`
	var questions []entity.GetQuestionQuizResponse
	err = r.db.SelectContext(ctx, &questions, questionsQuery, req.QuizId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("payload", req).Msg("repo::GetAllQuiz - Quiz with that id not found")
			quiz.Question = []entity.GetQuestionQuizResponse{}
		} else {
			log.Error().Err(err).Any("payload", req).Msg("repo::GetAllQuiz - Quiz with that id not found")
			return nil, err
		}
	} else {
		quiz.Question = questions
	}

	return &quiz, nil
}

func (r *quizRepository) FindUsersCompletedQuiz(ctx context.Context, req *entity.SubmitQuizRequest) error {
	query := `SELECT quiz_id, user_id FROM users_completed_quiz WHERE user_id = $1 AND quiz_id = $2`

	var quizID, userID string

	err := r.db.QueryRowContext(ctx, query, req.UserId, req.QuizId).Scan(&quizID, &userID)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::FindUsersCompletedQuiz - Failed to query quiz")
		return errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}

	return errmsg.NewCustomErrors(200, errmsg.WithMessage("Quiz already completed"))
}

func (r *quizRepository) SubmitQuiz(ctx context.Context, req *entity.SubmitQuizRequest) (*entity.SubmitQuizResponse, error) {

	query := `
        INSERT INTO users_completed_quiz (quiz_id, user_id)
        VALUES ($1, $2)
        RETURNING id, quiz_id, user_id, created_at
    `

	var response entity.SubmitQuizResponse
	err := r.db.QueryRowContext(ctx, query, req.QuizId, req.UserId).Scan(
		&response.Id,
		&response.QuizId,
		&response.UserId,
		&response.CreatedAt,
	)

	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::SubmitQuiz - Quiz with that id not found")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}

	response.Status = "completed"

	return &response, nil
}
