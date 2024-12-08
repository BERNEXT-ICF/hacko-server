package ports

import (
	"context"
	"hacko-app/internal/module/submission/entity"
)

type SubmissionRepository interface {
	FindAssignment(ctx context.Context, assignmentId string) error
	SubmitAssignment(ctx context.Context, req *entity.SubmitRequest) (*entity.SubmitResponse, error)
}

type SubmissionService interface {
	SubmitAssignment(ctx context.Context, req *entity.SubmitRequest) (*entity.SubmitResponse, error)
}
