package repository

import (
	"context"
	"database/sql"
	"errors"
	"hacko-app/internal/module/modules/entity"
	"hacko-app/internal/module/modules/ports"
	"hacko-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.ModulesRepository = &modulesRepository{}

type modulesRepository struct {
	db *sqlx.DB
}

func NewModulesRepository(db *sqlx.DB) *modulesRepository {
	return &modulesRepository{
		db: db,
	}
}

func (r *modulesRepository) CreateModules(ctx context.Context, req *entity.CreateModulesRequest) (*entity.CreateModulesResponse, error) {
	var res = new(entity.CreateModulesResponse)

	query := `
		INSERT INTO modules (
		creator_modules_id,
		materials_id,
		title,
		content,
		attachments,
		videos
	)
	SELECT
		$1, $2, $3, $4, $5, $6
	WHERE EXISTS (
		SELECT 1
		FROM materials
		WHERE id = $2 AND creator_materials_id = $1 
	)
	RETURNING id, title, content, attachments, videos, created_at, updated_at
    `

	err := r.db.QueryRowContext(ctx, query,
		req.UserId,
		req.MaterialsId,
		req.Title,
		req.Content,
		pq.Array(req.Attachments),
		pq.Array(req.Videos),
	).Scan(
		&res.Id,
		&res.Title,
		&res.Content,
		pq.Array(&res.Attachments),
		pq.Array(&res.Videos),
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Any("payload", req).Msg("repo::CreateModules - Module with the given data not found")
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Module not found or unable to create"))
		}

		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				log.Warn().Msg("repo::CreateModules - Foreign key constraint failed")
				return nil, errmsg.NewCustomErrors(409, errmsg.WithMessage("Referenced data not found"))
			default:
				log.Error().Err(pqErr).Any("payload", req).Msg("repo::CreateModules - Unhandled PostgreSQL error")
				return nil, err
			}
		}

		log.Error().Err(err).Any("payload", req).Msg("repo::CreateModules - Failed to insert module")
		return nil, err
	}

	return res, nil
}

func (r *modulesRepository) UpdateModules(ctx context.Context, req *entity.UpdateModulesRequest) (*entity.UpdateModulesResponse, error) {
	var res = new(entity.UpdateModulesResponse)

	query := `
		UPDATE modules
		SET 
			title = $1,
			content = $2,
			attachments = $3,
			videos = $4,
			updated_at = NOW()
		WHERE 
			id = $5 AND
			creator_modules_id = $6
		RETURNING 
			id, title, content, attachments, videos, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		req.Title,                    
		req.Content,                  
		pq.Array(req.Attachments),    
		pq.Array(req.Videos),         
		req.ModulesId,                  
		req.UserId,                   
	).Scan(
		&res.Id,
		&res.Title,
		&res.Content,
		pq.Array(&res.Attachments),
		pq.Array(&res.Videos),
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Any("payload", req).Msg("repo::UpdateModules - Module not found or unauthorized")
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Module not found or you are not authorized to update it"))
		}

		pqErr, ok := err.(*pq.Error)
		if ok {
			log.Error().Err(pqErr).Any("payload", req).Msg("repo::UpdateModules - Unhandled PostgreSQL error")
			return nil, err
		}

		log.Error().Err(err).Any("payload", req).Msg("repo::UpdateModules - Failed to update module")
		return nil, err
	}

	return res, nil
}

