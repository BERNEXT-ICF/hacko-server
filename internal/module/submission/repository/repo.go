package repository

import (
	"context"
	"database/sql"
	"hacko-app/internal/module/submission/entity"
	"hacko-app/internal/module/submission/ports"
	"hacko-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.SubmissionRepository = &submissionRepository{}

type submissionRepository struct {
	db *sqlx.DB
}

func NewSubmissionRepository(db *sqlx.DB) *submissionRepository {
	return &submissionRepository{
		db: db,
	}
}

func (r *submissionRepository) FindAssignment(ctx context.Context, assignmentId string) error {
    query := `
        SELECT id 
        FROM assignments 
        WHERE id = $1
    `

    var id int

	err := r.db.QueryRowContext(ctx, query, assignmentId).Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Error().Str("assignment_id", assignmentId).Msg("repo::FindAssignment - Assignment not found")
            return errmsg.NewCustomErrors(400, errmsg.WithMessage("Assignment not found"))
        }
        log.Error().Err(err).Msg("repo::FindAssignment - Failed to check assignment")
        return err
    }

    return nil
}


func (r *submissionRepository) SubmitAssignment(ctx context.Context, req *entity.SubmitRequest) (*entity.SubmitResponse, error) {
    // Query SQL untuk memasukkan data baru ke tabel submissions
    query := `
        INSERT INTO submissions (assignment_id, student_id, link, status, submitted_at)
        VALUES ($1, $2, $3, DEFAULT, DEFAULT)
        RETURNING id, student_id, link, status, submitted_at
    `

    var response entity.SubmitResponse

    // Eksekusi query
    err := r.db.QueryRowContext(ctx, query, req.AssignmentId, req.UserId, req.Link).
        Scan(&response.Id, &response.UserId, &response.Link, &response.Status, &response.SubmittedAt)
    if err != nil {
        log.Error().Err(err).Msg("repo::SubmitAssignment - Failed to submit assignment")
        return nil, err
    }

    return &response, nil
}
