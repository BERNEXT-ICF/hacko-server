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

	var classes []entity.GetClassResponse

	for rows.Next() {
		var class entity.GetClassResponse
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

func (r *classRepository) GetClassById(ctx context.Context, req *entity.GetClassByIdRequest) (*entity.GetClassResponse, error) {
	var res = new(entity.GetClassResponse)

	query := `
		SELECT 
			id,
			creator_class_id,
			title,
			description,
			image,
			video,
			status,
			created_at,
			updated_at
		FROM class
		WHERE id = ?
	`

	err := r.db.GetContext(ctx, res, r.db.Rebind(query), req.Id)

	if err != nil {
		log.Error().
			Err(err).
			Str("classId", req.Id).
			Msg("repo::GetClassById - Failed to retrieve class by ID")

		if err == sql.ErrNoRows {
			log.Warn().
				Str("classId", req.Id).
				Msg("repo::GetClassById - No class found with the provided ID")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Class not found"))
		}

		return nil, err
	}

	return res, nil
}

