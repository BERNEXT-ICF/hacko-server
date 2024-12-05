package ports

import (
	"context"
	"hacko-app/internal/module/modules/entity"
)

type ModulesRepository interface {
	CreateModules(ctx context.Context, req *entity.CreateModulesRequest) (*entity.CreateModulesResponse, error)
}

type ModulesService interface {
	CreateModules(ctx context.Context, req *entity.CreateModulesRequest) (*entity.CreateModulesResponse, error)
}
