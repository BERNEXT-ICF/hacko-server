package ports

import (
	"context"
	"hacko-app/internal/module/class/entity"
)

type ClassRepository interface {
	CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error)
	GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error)
	GetOverviewClassById(ctx context.Context, req *entity.GetOverviewClassByIdRequest) (*entity.GetOverviewClassByIdResponse, error)
	EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error
	UpdateClass(ctx context.Context, req *entity.UpdateClassRequest) (*entity.UpdateClassResponse, error)
	DeleteClass(ctx context.Context, req *entity.DeleteClassRequest) error
	UpdateVisibilityClass(ctx context.Context, req *entity.UpdateVisibilityClassRequest) (*entity.UpdateVisibilityClassResponse, error)
	GetAllUsersEnrolledClass(ctx context.Context, req *entity.GetAllUsersEnrolledClassRequest) (*entity.GetAllUsersEnrolledClassResponse, error)
}

type ClassService interface {
	CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error)
	GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error)
	GetOverviewClassById(ctx context.Context, req *entity.GetOverviewClassByIdRequest) (*entity.GetOverviewClassByIdResponse, error)
	EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error
	UpdateClass(ctx context.Context, req *entity.UpdateClassRequest) (*entity.UpdateClassResponse, error)
	DeleteClass(ctx context.Context, req *entity.DeleteClassRequest) error
	UpdateVisibilityClass(ctx context.Context, req *entity.UpdateVisibilityClassRequest) (*entity.UpdateVisibilityClassResponse, error)
	GetAllUsersEnrolledClass(ctx context.Context, req *entity.GetAllUsersEnrolledClassRequest) (*entity.GetAllUsersEnrolledClassResponse, error)
}
