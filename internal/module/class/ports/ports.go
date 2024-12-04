package ports

import (
	"context"
	"hacko-app/internal/module/class/entity"
)

type ClassRepository interface {
	CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error)
	GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error)
	GetClassById(ctx context.Context, req *entity.GetClassByIdRequest) (*entity.GetClassResponse, error)
	EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error
}

type ClassService interface {
	CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error)
	GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error)
	GetClassById(ctx context.Context, req *entity.GetClassByIdRequest) (*entity.GetClassResponse, error)
	EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error
}
