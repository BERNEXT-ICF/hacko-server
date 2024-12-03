package service

import (
	"context"
	"hacko-app/internal/module/class/entity"
	"hacko-app/internal/module/class/ports"

	"github.com/rs/zerolog/log"
)

var _ ports.ClassService = &classService{}

type classService struct {
	repo ports.ClassRepository
}

func NewClassService(repo ports.ClassRepository) *classService {
	return &classService{
		repo: repo,
	}
}

func (s *classService) CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error) {

	result, err := s.repo.CreateClass(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateClass - Failed to create class")
		return nil, err
	}

	return result, nil
}

func (s *classService) GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error) {
	classes, err := s.repo.GetAllClasses(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetAllClasses - Failed to retrieve classes from repository")
		return nil, err
	}

	return classes, nil
}

func (s *classService) GetClassById(ctx context.Context, req *entity.GetClassByIdRequest) (*entity.GetClassResponse, error) {
	class, err := s.repo.GetClassById(ctx, req)
	if err != nil {
		return nil, err
	}

	return class, nil
}

