package repository

import (
	"context"
	"database/sql"
	"hacko-app/internal/module/assignment/entity"
	"hacko-app/internal/module/assignment/ports"
	"hacko-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.AssignmentRepository = &assignmentRepository{}

type assignmentRepository struct {
	db *sqlx.DB
}

func NewAssignmentRepository(db *sqlx.DB) *assignmentRepository {
	return &assignmentRepository{
		db: db,
	}
}

func (r *assignmentRepository) CreateAssignment(ctx context.Context, req *entity.CreateAssignmentRequest) (*entity.CreateAssignmentResponse, error) {
	query := `
		INSERT INTO assignments (
			creator_assignment_id, 
			class_id, 
			title, 
			description, 
			due_date, 
			created_at, 
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
		RETURNING id, creator_assignment_id, class_id, title, description, due_date, created_at, updated_at
	`

	var response entity.CreateAssignmentResponse
	err := r.db.QueryRowContext(
		ctx,
		query,
		req.UserId,
		req.ClassId,
		req.Title,
		req.Description,
		req.DueDate,
	).Scan(
		&response.Id,
		&response.UserId,
		&response.ClassId,
		&response.Title,
		&response.Description,
		&response.DueDate,
		&response.CreatedAt,
		&response.UpdatedAt,
	)

	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateAssignment - Failed to insert assignment")
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Module not found or unable to create"))
	}

	return &response, nil
}

func (r *assignmentRepository) FindClass(ctx context.Context, req string) error {
	query := `SELECT id FROM class WHERE id = $1`

	var classId string
	err := r.db.QueryRowContext(ctx, query, req).Scan(&classId)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Str("class_id", req).Msg("repo::FindClass - Class not found")
			return errmsg.NewCustomErrors(404, errmsg.WithMessage("Class not found"))
		}
		log.Error().Err(err).Str("class_id", req).Msg("repo::FindClass - Failed to query class")
		return err
	}

	return nil
}

func (r *assignmentRepository) GetAllAssignmentByClassId(ctx context.Context, req *entity.GetAllAssignmentByClassIdRequest) ([]entity.GetAssignmentByClassIdResponse, error) {
	query := `
        SELECT id, creator_assignment_id, class_id, title, description, due_date, created_at, updated_at
        FROM assignments
        WHERE class_id = $1;
    `

	rows, err := r.db.QueryContext(ctx, query, req.ClassId)
	if err != nil {
		log.Error().Err(err).Str("class_id", req.ClassId).Msg("repo::GetAllAssignmentByClassId - Failed to query assignments")
		return nil, err
	}
	defer rows.Close()

	var assignments []entity.GetAssignmentByClassIdResponse
	for rows.Next() {
		var assignment entity.GetAssignmentByClassIdResponse
		err := rows.Scan(
			&assignment.Id,
			&assignment.CreatorAssignmentId,
			&assignment.ClassId,
			&assignment.Title,
			&assignment.Description,
			&assignment.DueDate,
			&assignment.CreatedAt,
			&assignment.UpdatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("repo::GetAllAssignmentByClassId - Failed to scan assignment")
			return nil, err
		}

		// Memanggil GetAssignmentStatus dengan parameter yang tepat
		status := r.GetAssignmentStatus(ctx, &entity.GetAssignmentStatusRequest{
			ClassId:      req.ClassId,
			AssignmentId: assignment.Id, // Menggunakan ID assignment untuk submission_id
			UserId:       req.UserId,    // Misalkan Anda ingin menggunakan UserId dari request
		})

		if status == "" {
			assignment.Status = "not_submit_yet"
		} else {
			assignment.Status = status
		}

		assignments = append(assignments, assignment)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("repo::GetAllAssignmentByClassId - Error during rows iteration")
		return nil, err
	}

	return assignments, nil
}

func (r *assignmentRepository) GetAssignmentStatus(ctx context.Context, req *entity.GetAssignmentStatusRequest) string {
	query := `
        SELECT status
        FROM submission
        WHERE class_id = $1 AND submission_id = $2 AND user_id = $3
        LIMIT 1; 
    `
	var status string
	err := r.db.QueryRowContext(ctx, query, req.ClassId, req.AssignmentId, req.UserId).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Str("class_id", req.ClassId).Msg("repo::GetAssignmentStatus - No assignments found")
			return "" // Mengembalikan string kosong jika tidak ada hasil
		}
		log.Error().Err(err).Str("class_id", req.ClassId).Msg("repo::GetAssignmentStatus - Failed to query status")
		return "" // Mengembalikan string kosong jika ada error lain
	}

	return status
}
