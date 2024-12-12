package ports

import (
	"context"
	"hacko-app/internal/module/assignment/entity"
)

type AssignmentRepository interface {
	// utils validation
	FindClass(ctx context.Context, req string) error
	GetAssignmentStatus(ctx context.Context, req *entity.GetAssignmentStatusRequest) string

	// admin contract
	CreateAssignment(ctx context.Context, req *entity.CreateAssignmentRequest) (*entity.CreateAssignmentResponse, error)
	GetAllAssignmentByClassIdAdmin(ctx context.Context, req *entity.GetAllAssignmentByClassIdAdminRequest) (*[]entity.GetAllAssignmentByClassIdAdminResponse, error)
	
	// user contract
	GetAssignmentDetails(ctx context.Context, req *entity.GetAssignmentDetailsRequest)(*entity.GetAssignmentDetailsResponse, error)
	GetAllAssignmentByClassId(ctx context.Context, req *entity.GetAllAssignmentByClassIdRequest) ([]entity.GetAssignmentByClassIdResponse, error)
}

type AssignmentService interface {
	// admin contract
	CreateAssignment(ctx context.Context, req *entity.CreateAssignmentRequest) (*entity.CreateAssignmentResponse, error)
	GetAllAssignmentByClassIdAdmin(ctx context.Context, req *entity.GetAllAssignmentByClassIdAdminRequest) (*[]entity.GetAllAssignmentByClassIdAdminResponse, error)

	// user contract
	GetAllAssignmentByClassId(ctx context.Context, req *entity.GetAllAssignmentByClassIdRequest) ([]entity.GetAssignmentByClassIdResponse, error)
	GetAssignmentDetails(ctx context.Context, req *entity.GetAssignmentDetailsRequest)(*entity.GetAssignmentDetailsResponse, error)
}
