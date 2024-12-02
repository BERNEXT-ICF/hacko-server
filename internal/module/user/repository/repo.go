package repository

import (
	"context"
	"database/sql"
	"hacko-app/internal/module/user/entity"
	"hacko-app/internal/module/user/ports"
	"hacko-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.UserRepository = &userRepository{}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	var res = new(entity.RegisterResponse)

	query := `
	INSERT INTO users (
    email,
    name,
    password
	)
	VALUES (
		?, ?, ?
	)
	RETURNING id
	`
	err := r.db.QueryRowContext(ctx, r.db.Rebind(query), req.Email, req.Name, req.HassedPassword).Scan(&res.Id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Error().Err(err).Any("payload", req).Msg("repo::Register - Failed to insert user")
			return nil, err
		}

		switch pqErr.Code.Name() {
		case "unique_violation":
			log.Warn().Msg("Email already registered")
			return nil, errmsg.NewCustomErrors(409, errmsg.WithMessage("Email sudah terdaftar"))
		case "not_null_violation":
			log.Error().Err(err).Any("payload", req).Msg("Missing required fields")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Data tidak lengkap"))
		case "syntax_error":
			log.Error().Err(err).Any("payload", req).Msg("Query syntax error")
			return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Kesalahan sintaks"))
		default:
			log.Error().Err(err).Any("payload", req).Msg("Unhandled pq.Error")
			return nil, err
		}

	}

	return res, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.UserResult, error) {
	var res = new(entity.UserResult)

	query := `
	SELECT
		u.id,
		u.role,
		u.name,
		u.email,
		u.password
	FROM
		users u
	WHERE
		u.email = ?
`

	err := r.db.GetContext(ctx, res, r.db.Rebind(query), email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Str("email", email).Msg("repo::FindByEmail - User not found")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Email atau password salah"))
		}
		log.Error().Err(err).Str("email", email).Msg("repo::FindByEmail - Failed to get user")
		return nil, err
	}

	return res, nil
}

func (r *userRepository) FindById(ctx context.Context, id string) (*entity.ProfileResponse, error) {
	var res = new(entity.ProfileResponse)

	query := `
	SELECT
		u.id,
		u.role,
		u.name,
		u.email
	FROM
		users u
	WHERE
		u.id = ?
`

	err := r.db.GetContext(ctx, res, r.db.Rebind(query), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Str("id", id).Msg("repo::FindById - User not found")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("User tidak ditemukan"))
		}

		log.Error().Err(err).Str("id", id).Msg("repo::FindById - Failed to get user")
		return nil, err
	}

	return res, nil
}
