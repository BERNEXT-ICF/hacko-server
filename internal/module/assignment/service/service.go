package service

import (
	"context"
	"hacko-app/internal/module/assignment/entity"
	"hacko-app/internal/module/assignment/ports"
)

var _ ports.AssignmentService = &assignmentService{}

type assignmentService struct {
	repo ports.AssignmentRepository
}

func NewAssignmentService(repo ports.AssignmentRepository) *assignmentService {
	return &assignmentService{
		repo: repo,
	}
}

func (s *assignmentService) CreateAssignment(ctx context.Context, req *entity.CreateAssignmentRequest) (*entity.CreateAssignmentResponse, error) {
	response, err := s.repo.CreateAssignment(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *assignmentService) GetAllAssignmentByClassId(ctx context.Context, req *entity.GetAllAssignmentByClassIdRequest) ([]entity.GetAssignmentByClassIdResponse, error) {

	err := s.repo.FindClass(ctx, req.ClassId)
	if err != nil {
		return nil, err
	}

	response, err := s.repo.GetAllAssignmentByClassId(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *assignmentService) GetAssignmentDetails(ctx context.Context, req *entity.GetAssignmentDetailsRequest) (*entity.GetAssignmentDetailsResponse, error) {
	response, err := s.repo.GetAssignmentDetails(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
