package ports

import (
	"context"
	"hacko-app/internal/module/modules/entity"
)

type ModulesRepository interface {
	CreateModules(ctx context.Context, req *entity.CreateModulesRequest) (*entity.CreateModulesResponse, error)
	UpdateModules(ctx context.Context, req *entity.UpdateModulesRequest) (*entity.UpdateModulesResponse, error)
	DeleteModules(ctx context.Context, req *entity.DeleteModulesRequest) error
}

type ModulesService interface {
	CreateModules(ctx context.Context, req *entity.CreateModulesRequest) (*entity.CreateModulesResponse, error)
	UpdateModules(ctx context.Context, req *entity.UpdateModulesRequest) (*entity.UpdateModulesResponse, error)
	DeleteModules(ctx context.Context, req *entity.DeleteModulesRequest) error
}
