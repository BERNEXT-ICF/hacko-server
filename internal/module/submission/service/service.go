package service

import (
	"context"
	"hacko-app/internal/module/submission/entity"
	"hacko-app/internal/module/submission/ports"
)

var _ ports.SubmissionService = &submissionService{}

type submissionService struct {
	repo ports.SubmissionRepository
}

func NewSubmissionService(repo ports.SubmissionRepository) *submissionService {
	return &submissionService{
		repo: repo,
	}
}

func (s *submissionService) SubmitAssignment(ctx context.Context, req *entity.SubmitRequest) (*entity.SubmitResponse, error) {

	err := s.repo.FindAssignment(ctx, req.AssignmentId)
	if err != nil {
		return nil, err
	}

	response, err := s.repo.SubmitAssignment(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *submissionService) GetSubmissionDetails(ctx context.Context, req *entity.GetSubmissionDetailsRequest) (*entity.GetSubmissionDetailsResponse, error){
	response, err := s.repo.GetSubmissionDetails(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
