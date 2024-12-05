package repository

import (
	"context"
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
			log.Error().Err(err).Any("payload", req).Msg("repo::CreateClass - Failed to insert class")
			return nil, err
		}

		switch pqErr.Code.Name() {
		case "foreign_key_violation":
			log.Warn().Msg("repo::EnrollClass - Class with the id not found")
			return nil, errmsg.NewCustomErrors(409, errmsg.WithMessage("Class with that id not found"))
		default:
			log.Error().Err(err).Any("payload", req).Msg("repo::EnrollClass - Unhandled pq.Error")
			return nil, err
		}
	}

	return res, nil
}
