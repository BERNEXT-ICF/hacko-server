package service

import (
	"context"
	"hacko-app/internal/module/modules/entity"
	"hacko-app/internal/module/modules/ports"
)

var _ ports.ModulesService = &modulesService{}

type modulesService struct {
	repo ports.ModulesRepository
}

func NewModulesService(repo ports.ModulesRepository) *modulesService {
	return &modulesService{
		repo: repo,
	}
}

func (s *modulesService) CreateModules(ctx context.Context, req *entity.CreateModulesRequest) (*entity.CreateModulesResponse, error) {

	response, err := s.repo.CreateModules(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *modulesService) UpdateModules(ctx context.Context, req *entity.UpdateModulesRequest) (*entity.UpdateModulesResponse, error) {

	response, err := s.repo.UpdateModules(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
