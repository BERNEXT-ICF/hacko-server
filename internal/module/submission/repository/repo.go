package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *submissionRepository) GetSubmissionDetails(ctx context.Context, req *entity.GetSubmissionDetailsRequest) (*entity.GetSubmissionDetailsResponse, error) {
    query := `
        SELECT 
            s.id AS id,
            s.id AS submission_id,
            u.name AS name,
            u.image_url AS image_url,
            s.link AS link,
            s.status AS status,
            s.grade AS grade,
            s.feedback AS feedback,
            s.submitted_at AS submitted_at,
            s.graded_at AS graded_at
        FROM 
            submissions s
        INNER JOIN 
            users u ON s.student_id = u.id
        INNER JOIN 
            assignments a ON s.assignment_id = a.id
        WHERE 
            s.id = $1 AND u.id = $2
    `

    var response entity.GetSubmissionDetailsResponse
    err := r.db.QueryRowContext(ctx, query, req.SubmissionId, req.UserId).Scan(
        &response.Id,
        &response.SubmissionId,
        &response.Name,
        &response.Image,
        &response.Link,
        &response.Status,
        &response.Grade,
        &response.Feedback,
        &response.SubmittedAt,
        &response.GradedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            log.Error().Any("payload", req).Msg("No submission details found")
            return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Submission with that ID not found"))
        }
        log.Error().Err(err).Msg("Failed to get submission details")
        return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
    }

    return &response, nil
}

func (r *submissionRepository) GradingSubmission(ctx context.Context, req *entity.GradingSubmissionRequest) (*entity.GradingSubmissionResponse, error) {
    query := `
        UPDATE 
            submissions
        SET 
            grade = $1,
            feedback = $2,
            status = $3,
            graded_at = NOW()
        WHERE 
            id = $4
        RETURNING 
            id, grade, feedback, status, graded_at
    `

    var response entity.GradingSubmissionResponse

    err := r.db.QueryRowContext(
        ctx,
        query,
        req.Grade,
        req.Feedback,
        req.Status,
        req.SubmissionId,
    ).Scan(
        &response.Id,
        &response.Grade,
        &response.Feedback,
        &response.Status,
        &response.GradedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            log.Error().Any("payload", req).Msg("Submission not found for grading")
            return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Submission not found"))
        }
        log.Error().Err(err).Msg("Failed to grade submission")
        return nil, err
    }

    return &response, nil
}

