package ports

import (
	"context"
	"hacko-app/internal/module/class/entity"
)

type ClassRepository interface {

	// utils repo contract
	FindClass(ctx context.Context, id string) error
	CheckEnrollment(ctx context.Context, req *entity.AddUsersToClassRequest) error
	GetAllSyllabus(ctx context.Context, classId string) ([]entity.GetMaterialResponse, error)

	// admin repo contract
	CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error)
	UpdateClass(ctx context.Context, req *entity.UpdateClassRequest) (*entity.UpdateClassResponse, error)
	DeleteClass(ctx context.Context, req *entity.DeleteClassRequest) error
	UpdateVisibilityClass(ctx context.Context, req *entity.UpdateVisibilityClassRequest) (*entity.UpdateVisibilityClassResponse, error)
	GetAllUsersEnrolledClass(ctx context.Context, req *entity.GetAllUsersEnrolledClassRequest) (*entity.GetAllUsersEnrolledClassResponse, error)
	DeleteStudentClass(ctx context.Context, req *entity.DeleteUsersClassRequest) error
	GetAllStudentNotEnrolledClass(ctx context.Context, req *entity.GetAllUserNotEnrolledClassRequest) (*entity.GetAllUserNotEnrolledClassResponse, error)
	AddUserToClass(ctx context.Context, req *entity.AddUsersToClassRequest) (*entity.AddUsersToClassResponse, error)
	GetAllClassAdmin(ctx context.Context, req *entity.GetAllClassAdminRequest) (*[]entity.GetAllClassAdminResponse, error)

	// users repo contract
	GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error)
	GetOverviewClassById(ctx context.Context, req *entity.GetOverviewClassByIdRequest) (*entity.GetOverviewClassByIdResponse, error)
	EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error
	TrackModule(ctx context.Context, req *entity.TrackModuleRequest) (*entity.TrackModuleResponse, error)
	GetProgress(ctx context.Context, req *entity.GetProgressRequest) (*float64, error)
}

type ClassService interface {

	// admin service contract
	CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error)
	UpdateClass(ctx context.Context, req *entity.UpdateClassRequest) (*entity.UpdateClassResponse, error)
	DeleteClass(ctx context.Context, req *entity.DeleteClassRequest) error
	UpdateVisibilityClass(ctx context.Context, req *entity.UpdateVisibilityClassRequest) (*entity.UpdateVisibilityClassResponse, error)
	GetAllUsersEnrolledClass(ctx context.Context, req *entity.GetAllUsersEnrolledClassRequest) (*entity.GetAllUsersEnrolledClassResponse, error)
	DeleteStudentClass(ctx context.Context, req *entity.DeleteUsersClassRequest) error
	GetAllStudentNotEnrolledClass(ctx context.Context, req *entity.GetAllUserNotEnrolledClassRequest) (*entity.GetAllUserNotEnrolledClassResponse, error)
	AddUserToClass(ctx context.Context, req *entity.AddUsersToClassRequest) (*entity.AddUsersToClassResponse, error)
	GetAllClassAdmin(ctx context.Context, req *entity.GetAllClassAdminRequest) (*[]entity.GetAllClassAdminResponse, error)

	// users service contract
	GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error)
	GetOverviewClassById(ctx context.Context, req *entity.GetOverviewClassByIdRequest) (*entity.GetOverviewClassByIdResponse, error)
	EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error
	TrackModule(ctx context.Context, req *entity.TrackModuleRequest) (*entity.TrackModuleResponse, error)
}
