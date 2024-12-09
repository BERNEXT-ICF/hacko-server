package ports

import (
	"context"
	"hacko-app/internal/module/quiz/entity"
)

type QuizRepository interface {
	FindClass(ctx context.Context, req string) error
	CreateQuiz(ctx context.Context, req *entity.CreateQuizRequest) (*entity.CreateQuizResponse, error)
	FindQuiz(ctx context.Context, req int) error
	CreateQuestionQuiz(ctx context.Context, req *entity.CreateQuestionQuizRequest) (*entity.CreateQuestionQuizResponse, error)
}

type QuizService interface {
	CreateQuiz(ctx context.Context, req *entity.CreateQuizRequest) (*entity.CreateQuizResponse, error)
	CreateQuestionQuiz(ctx context.Context, req *entity.CreateQuestionQuizRequest) (*entity.CreateQuestionQuizResponse, error)
}
