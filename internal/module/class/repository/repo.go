package repository

import (
	"context"
	"database/sql"
	"hacko-app/internal/module/class/entity"
	"hacko-app/internal/module/class/ports"
	"hacko-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.ClassRepository = &classRepository{}

type classRepository struct {
	db *sqlx.DB
}

func NewClassRepository(db *sqlx.DB) *classRepository {
	return &classRepository{
		db: db,
	}
}

func (r *classRepository) CreateClass(ctx context.Context, req *entity.CreateClassRequest) (*entity.CreateClassResponse, error) {
	var res = new(entity.CreateClassResponse)

	query := `
	INSERT INTO class (
		creator_class_id,
		title,
		description,
		image,
		video,
		status
	)
	VALUES (
		?, ?, ?, ?, ?, ?
	)
	RETURNING id, creator_class_id, title, status, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query), req.UserId, req.Title, req.Description, req.Image, req.Video, req.Status).
		Scan(&res.Id, &res.CreatorId, &res.Title, &res.Status, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Error().Err(err).Any("payload", req).Msg("repo::CreateClass - Failed to insert class")
			return nil, err
		}
		log.Error().Err(pqErr).Any("payload", req).Msg("repo::CreateClass - Database error")
		return nil, err
	}

	return res, nil
}

func (r *classRepository) GetAllClasses(ctx context.Context) (*entity.GetAllClassesResponse, error) {

	query := `
		SELECT 
			id,
			title,
			description,
			image,
			video,
			status,
			creator_class_id,
			created_at,
			updated_at
		FROM 
			class
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetAllClasses - Failed to execute query")
		return nil, err
	}
	defer rows.Close()

	var classes []*entity.GetClassResponse

	for rows.Next() {
		class := new(entity.GetClassResponse)
		err := rows.Scan(
			&class.ID,
			&class.Title,
			&class.Description,
			&class.Image,
			&class.Video,
			&class.Status,
			&class.CreatorClassID,
			&class.CreatedAt,
			&class.UpdatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("repo::GetAllClasses - Failed to scan row")
			return nil, err
		}
		classes = append(classes, class)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("repo::GetAllClasses - Error occurred during rows iteration")
		return nil, err
	}

	response := &entity.GetAllClassesResponse{
		Classes: classes,
		Total:   len(classes),
	}

	return response, nil
}

func (r *classRepository) GetOverviewClassById(ctx context.Context, req *entity.GetOverviewClassByIdRequest) (*entity.GetOverviewClassByIdResponse, error) {
	var res = new(entity.GetOverviewClassByIdResponse)

	query := `
		SELECT 
			c.id,
			c.creator_class_id,
			c.title,
			c.description,
			c.image,
			c.video,
			c.status,
			c.created_at,
			c.updated_at,
			COALESCE(uc.enrollment_status, 'not_enrolled') AS enrollment_status
		FROM class c
		LEFT JOIN users_classes uc ON uc.class_id = c.id AND uc.user_id = ? 
		WHERE c.id = ?
	`

	err := r.db.GetContext(ctx, res, r.db.Rebind(query), req.UserId, req.Id)
	if err != nil {
		log.Error().
			Err(err).
			Str("classId", req.Id).
			Msg("repo::GetOverviewClassById - Failed to retrieve class by ID")

		if err == sql.ErrNoRows {
			log.Warn().
				Str("classId", req.Id).
				Msg("repo::GetOverviewClassById - No class found with the provided ID")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Class with the ID was not found"))
		}

		return nil, err
	}

	return res, nil
}

func (r *classRepository) EnrollClass(ctx context.Context, req *entity.EnrollClassRequest) error {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM users_classes 
		WHERE user_id = $1 AND class_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, req.UserId, req.ClassId).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if user is already enrolled in class")
		return err
	}

	if count > 0 {
		return errmsg.NewCustomErrors(400, errmsg.WithMessage("User is already enrolled in the class"))
	}

	insertQuery := `
		INSERT INTO users_classes (user_id, class_id, enrollment_status, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err = r.db.ExecContext(ctx, insertQuery, req.UserId, req.ClassId, "active")
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Error().Err(err).Any("payload", req).Msg("repo::EnrollClass - Failed to enroll user")
			return err
		}

		switch pqErr.Code.Name() {
		case "foreign_key_violation":
			log.Warn().Msg("repo::EnrollClass - Class with the id not found")
			return errmsg.NewCustomErrors(409, errmsg.WithMessage("Class with that id was not found"))
		default:
			log.Error().Err(err).Any("payload", req).Msg("repo::EnrollClass - Unhandled pq.Error")
			return err
		}
	}

	return nil
}

func (r *classRepository) UpdateClass(ctx context.Context, req *entity.UpdateClassRequest) (*entity.UpdateClassResponse, error) {
	var res = new(entity.UpdateClassResponse)

	query := `
		UPDATE class 
		SET 
			title = ?, 
			description = ?, 
			image = ?, 
			video = ?, 
			status = ?, 
			updated_at = NOW() 
		WHERE id = ? AND creator_class_id = ?
		RETURNING id, title, description, image, video, status, created_at, updated_at, creator_class_id;
	`

	err := r.db.GetContext(ctx, res, r.db.Rebind(query), req.Title, req.Description, req.Image, req.Video, req.Status, req.Id, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::UpdateClass - Failed to update class")
		if err == sql.ErrNoRows {
			log.Warn().
				Msg("repo::UpdateClass - No class found with the provided ID")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Class not found or you do not have update access to the class"))
		}
		return nil, err
	}

	return res, nil
}

func (r *classRepository) DeleteClass(ctx context.Context, req *entity.DeleteClassRequest) error {
	query := `
        DELETE FROM class
        WHERE id = ? AND creator_class_id = ?
    `

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id, req.UserId)
	if err != nil {
		log.Error().
			Err(err).
			Any("payload", req).
			Msg("repo::DeleteClass - Failed to delete class")
		return errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to delete class"))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().
			Err(err).
			Msg("repo::DeleteClass - Failed to get rows affected")
		return errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to process deletion"))
	}

	if rowsAffected == 0 {
		log.Warn().
			Any("payload", req).
			Msg("repo::DeleteClass - No rows affected, invalid classId or userId")
		return errmsg.NewCustomErrors(404, errmsg.WithMessage("Class not found or unauthorized"))
	}

	return nil
}

func (r *classRepository) UpdateVisibilityClass(ctx context.Context, req *entity.UpdateVisibilityClassRequest) (*entity.UpdateVisibilityClassResponse, error) {
	var res = new(entity.UpdateVisibilityClassResponse)

	query := `
		UPDATE class
		SET status = CASE 
			WHEN status = 'public' THEN 'draf'
			WHEN status = 'draf' THEN 'public'
			ELSE status
		END
		WHERE id = $1 AND creator_class_id = $2
		RETURNING id, title, status
	`

	err := r.db.GetContext(ctx, res, query, req.Id, req.UserId)
	if err != nil {
		log.Error().
			Err(err).
			Any("payload", req).
			Msg("repo::UpdateVisibilityClass - Failed to update class visibility")

		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Class not found or unauthorized access"))
		}
		return nil, err
	}

	return res, nil
}
