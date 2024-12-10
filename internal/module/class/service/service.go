package service

import (
	"context"
	"hacko-app/internal/module/class/entity"
	"hacko-app/internal/module/class/ports"

	// "hacko-app/pkg/response"

	// "github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.ClassService = &classService{}

type classService struct {
	repo ports.ClassRepository
}

func NewClassService(repo ports.ClassRepository) *classService {
	return &classService{
		repo: repo,
	}
}

func (s *classService) CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error) {

	result, err := s.repo.CreateClass(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateClass - Failed to create class")
		return nil, err
	}

	return result, nil
}

func (s *classService) GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error) {
	classes, err := s.repo.GetAllClasses(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetAllClasses - Failed to retrieve classes from repository")
		return nil, err
	}

	return classes, nil
}

func (s *classService) GetOverviewClassById(ctx context.Context, req *entity.GetOverviewClassByIdRequest) (*entity.GetOverviewClassByIdResponse, error) {
	if err := s.repo.FindClass(ctx, req.Id); err != nil {
		return nil, err
	}

	class, err := s.repo.GetOverviewClassById(ctx, req)
	if err != nil {
		return nil, err
	}

	syllabus, err := s.repo.GetAllSyllabus(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	class.Syllabus = syllabus

	return class, nil
}

func (s *classService) EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error {
	err := s.repo.EnrollClass(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *classService) UpdateClass(ctx context.Context, req *entity.UpdateClassRequest) (*entity.UpdateClassResponse, error) {
	updatedClass, err := s.repo.UpdateClass(ctx, req)
	if err != nil {
		return nil, err
	}

	return updatedClass, nil
}

func (s *classService) DeleteClass(ctx context.Context, req *entity.DeleteClassRequest) error {
	err := s.repo.DeleteClass(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *classService) UpdateVisibilityClass(ctx context.Context, req *entity.UpdateVisibilityClassRequest) (*entity.UpdateVisibilityClassResponse, error) {
	res, err := s.repo.UpdateVisibilityClass(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *classService) GetAllUsersEnrolledClass(ctx context.Context, req *entity.GetAllUsersEnrolledClassRequest) (*entity.GetAllUsersEnrolledClassResponse, error) {
	res, err := s.repo.GetAllUsersEnrolledClass(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *classService) DeleteStudentClass(ctx context.Context, req *entity.DeleteUsersClassRequest) error {
	err := s.repo.DeleteStudentClass(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *classService) GetAllStudentNotEnrolledClass(ctx context.Context, req *entity.GetAllUserNotEnrolledClassRequest) (*entity.GetAllUserNotEnrolledClassResponse, error) {
	if err := s.repo.FindClass(ctx, req.ClassId); err != nil {
		return nil, err
	}

	res, err := s.repo.GetAllStudentNotEnrolledClass(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *classService) AddUserToClass(ctx context.Context, req *entity.AddUsersToClassRequest) (*entity.AddUsersToClassResponse, error) {
	if err := s.repo.FindClass(ctx, req.ClassId); err != nil {
		return nil, err
	}

	if err := s.repo.CheckEnrollment(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.repo.AddUserToClass(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *classService) TrackModule(ctx context.Context, req *entity.TrackModuleRequest) (*entity.TrackModuleResponse, error){
	res, err := s.repo.TrackModule(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil	
}
