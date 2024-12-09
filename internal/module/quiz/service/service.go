package service

import (
	"context"
	"hacko-app/internal/module/quiz/entity"
	"hacko-app/internal/module/quiz/ports"
)

var _ ports.QuizService = &quizService{}

type quizService struct {
	repo ports.QuizRepository
}

func NewQuizService(repo ports.QuizRepository) *quizService {
	return &quizService{
		repo: repo,
	}
}

func (s *quizService) CreateQuiz(ctx context.Context, req *entity.CreateQuizRequest) (*entity.CreateQuizResponse, error) {

	err := s.repo.FindClass(ctx, req.ClassId)
	if err != nil {
		return nil, err
	}

	response, err := s.repo.CreateQuiz(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *quizService) CreateQuestionQuiz(ctx context.Context, req *entity.CreateQuestionQuizRequest) (*entity.CreateQuestionQuizResponse, error) {
	err := s.repo.FindQuiz(ctx, req.QuizId)
	if err != nil {
		return nil, err
	}

	response, err := s.repo.CreateQuestionQuiz(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *quizService) GetAllQuiz(ctx context.Context, req *entity.GetAllQuizRequest) ([]entity.GetAllQuizResponse, error) {
	err := s.repo.FindClass(ctx, req.ClassId)
	if err != nil {
		return nil, err
	}

	response, err := s.repo.GetAllQuiz(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
