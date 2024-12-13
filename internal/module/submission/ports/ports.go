package ports

import (
	"context"
	"hacko-app/internal/module/submission/entity"
)

type SubmissionRepository interface {
	// utils contract
	FindAssignment(ctx context.Context, assignmentId string) error

	// users contract
	SubmitAssignment(ctx context.Context, req *entity.SubmitRequest) (*entity.SubmitResponse, error)

	// admin contract
	GetSubmissionDetails(ctx context.Context, req *entity.GetSubmissionDetailsRequest) (*entity.GetSubmissionDetailsResponse, error)
	GradingSubmission(ctx context.Context, req *entity.GradingSubmissionRequest) (*entity.GradingSubmissionResponse, error)
}

type SubmissionService interface {
	// user contract
	SubmitAssignment(ctx context.Context, req *entity.SubmitRequest) (*entity.SubmitResponse, error)

	// admin contract
	GetSubmissionDetails(ctx context.Context, req *entity.GetSubmissionDetailsRequest) (*entity.GetSubmissionDetailsResponse, error)
	GradingSubmission(ctx context.Context, req *entity.GradingSubmissionRequest) (*entity.GradingSubmissionResponse, error)
}
