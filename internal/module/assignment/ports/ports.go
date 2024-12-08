package ports

import (
	"context"
	"hacko-app/internal/module/assignment/entity"
)

type AssignmentRepository interface {
	CreateAssignment(ctx context.Context, req *entity.CreateAssignmentRequest) (*entity.CreateAssignmentResponse, error)
	GetAllAssignmentByClassId(ctx context.Context, req *entity.GetAllAssignmentByClassIdRequest) ([]entity.GetAssignmentByClassIdResponse, error)
	GetAssignmentStatus(ctx context.Context, req *entity.GetAssignmentStatusRequest) string
	FindClass(ctx context.Context, req string) error
	GetAssignmentDetails(ctx context.Context, req *entity.GetAssignmentDetailsRequest)(*entity.GetAssignmentDetailsResponse, error)
}

type AssignmentService interface {
	CreateAssignment(ctx context.Context, req *entity.CreateAssignmentRequest) (*entity.CreateAssignmentResponse, error)
	GetAllAssignmentByClassId(ctx context.Context, req *entity.GetAllAssignmentByClassIdRequest) ([]entity.GetAssignmentByClassIdResponse, error)
	GetAssignmentDetails(ctx context.Context, req *entity.GetAssignmentDetailsRequest)(*entity.GetAssignmentDetailsResponse, error)
}
