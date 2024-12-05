package ports

import (
	"context"
	"hacko-app/internal/module/materials/entity"
)

type MaterialsRepository interface {
	CreateMaterials(ctx context.Context, req *entity.CreateMaterialsRequest) (*entity.CreateMaterialsResponse, error)
}

type MaterialsService interface {
	CreateMaterials(ctx context.Context, req *entity.CreateMaterialsRequest) (*entity.CreateMaterialsResponse, error)
}
