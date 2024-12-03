package service

import (
	"context"
	"hacko-app/internal/module/class/entity"
	"hacko-app/internal/module/class/ports"
	"hacko-app/pkg/errmsg"

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
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Creator id not found"))
	}

	return result, nil
}

