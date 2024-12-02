package service

import "hacko-app/internal/module/z_v2/ports"

var _ ports.XxxService = &xxxService{}

type xxxService struct {
	repo ports.XxxRepository
}

func NewXxxService(repo ports.XxxRepository) *xxxService {
	return &xxxService{
		repo: repo,
	}
}
