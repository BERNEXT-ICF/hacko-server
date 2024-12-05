package service

import (
	"context"
	"hacko-app/internal/module/materials/entity"
	"hacko-app/internal/module/materials/ports"
)

var _ ports.MaterialsService = &materialsService{}

type materialsService struct {
	repo ports.MaterialsRepository
}

func NewMaterialsService(repo ports.MaterialsRepository) *materialsService {
	return &materialsService{
		repo: repo,
	}
}

func (s *materialsService) CreateMaterials(ctx context.Context, req *entity.CreateMaterialsRequest) (*entity.CreateMaterialsResponse, error) {

	result, err := s.repo.CreateMaterials(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
