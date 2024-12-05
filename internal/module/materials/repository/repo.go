package repository

import (
	"context"
	"database/sql"
	"errors"
	"hacko-app/internal/module/materials/entity"
	"hacko-app/internal/module/materials/ports"
	"hacko-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.MaterialsRepository = &materialsRepository{}

type materialsRepository struct {
	db *sqlx.DB
}

func NewMaterialsRepository(db *sqlx.DB) *materialsRepository {
	return &materialsRepository{
		db: db,
	}
}

func (r *materialsRepository) CreateMaterials(ctx context.Context, req *entity.CreateMaterialsRequest) (*entity.CreateMaterialsResponse, error) {
	var res = new(entity.CreateMaterialsResponse)

	query := `
	INSERT INTO materials (
		creator_materials_id,
		class_id,
		title
	)
	VALUES (
		?, ?, ?
	)
	RETURNING id, creator_materials_id, class_id, title, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query), req.UserId, req.ClassId, req.Title).
		Scan(&res.Id, &res.CreatorId, &res.ClassId, &res.Title, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Error().Err(err).Any("payload", req).Msg("repo::CreateMaterials - Failed to insert class")
			return nil, err
		}

		switch pqErr.Code.Name() {
		case "foreign_key_violation":
			log.Warn().Msg("repo::CreateMaterials - Class with the id not found")
			return nil, errmsg.NewCustomErrors(409, errmsg.WithMessage("Class with that id not found"))
		default:
			log.Error().Err(err).Any("payload", req).Msg("repo::EnrollClass - Unhandled pq.Error")
			return nil, err
		}
	}

	return res, nil
}

func (r *materialsRepository) UpdateMaterials(ctx context.Context, req *entity.UpdateMaterialsRequest) (*entity.UpdateMaterialsResponse, error) {
	var res = new(entity.UpdateMaterialsResponse)

	query := `
		UPDATE materials
		SET title = ?, updated_at = NOW()
		WHERE id = ? AND creator_materials_id = ?
		RETURNING id, title, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query), req.Title, req.MaterialId, req.UserId).
		Scan(&res.MaterialId, &res.Title, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Any("payload", req).Msg("repo::UpdateMaterials - Material with the given ID not found or unauthorized")
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Material with the given ID not found or unauthorized"))
		}

		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				log.Warn().Msg("repo::UpdateMaterials - Foreign key violation detected")
				return nil, errmsg.NewCustomErrors(409, errmsg.WithMessage("Foreign key violation detected"))
			default:
				log.Error().Err(err).Any("payload", req).Msg("repo::UpdateMaterials - Unhandled pq.Error")
				return nil, err
			}
		}

		log.Error().Err(err).Any("payload", req).Msg("repo::UpdateMaterials - Unexpected error")
		return nil, err
	}

	return res, nil
}

func (r *materialsRepository) DeleteMaterials(ctx context.Context, req *entity.DeleteMaterialsRequest) error {
	query := `
		DELETE FROM materials
		WHERE id = $1 AND creator_materials_id = $2
		RETURNING id
	`

	var deletedMaterialId int
	err := r.db.QueryRowContext(ctx, query, req.MaterialId, req.UserId).Scan(&deletedMaterialId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Any("payload", req).Msg("repo::DeleteMaterials - Material with the given ID not found or user is not authorized")
			return errmsg.NewCustomErrors(404, errmsg.WithMessage("Material not found or you are not authorized to delete it"))
		}
		log.Error().Err(err).Any("payload", req).Msg("repo::DeleteMaterials - Failed to delete material")
		return err
	}

	return nil
}
