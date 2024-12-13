package repository

import (
	"context"
	"database/sql"
	"errors"
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
		SELECT DISTINCT
			c.id,
			c.title,
			c.description,
			c.image,
			c.video,
			c.status,
			c.creator_class_id,
			c.created_at,
			c.updated_at,
			COALESCE(uc.enrollment_status, 'not_enrolled') AS status_enrollment,
			COALESCE(up.progress, '0') AS progress
		FROM 
			class c
		LEFT JOIN 
			users_classes uc ON c.id = uc.class_id
		LEFT JOIN 
			users_progress up ON c.id = up.class_id
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
			&class.StatusEnrollment,
			&class.Progress,
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



func (r *classRepository) FindClass(ctx context.Context, id string) error {
	query := `
        SELECT 
            1 
        FROM 
            class
        WHERE 
            id = $1
    `

	var exists int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Any("class_id", id).Msg("repo::FindClass - Class not found")
			return errmsg.NewCustomErrors(404, errmsg.WithMessage("Class not found"))
		}

		log.Error().Err(err).Any("class_id", id).Msg("repo::FindClass - Failed to query class")
		return err
	}

	return nil
}

func (r *classRepository) GetAllSyllabus(ctx context.Context, classId string) ([]entity.GetMaterialResponse, error) {
	materialsQuery := `
        SELECT 
            id, 
            title 
        FROM 
            materials 
        WHERE 
            class_id = $1
    `

	materialsRows, err := r.db.QueryContext(ctx, materialsQuery, classId)
	if err != nil {
		log.Error().Err(err).Str("class_id", classId).Msg("repo::GetAllSyllabus - Failed to query materials")
		return nil, err
	}
	defer materialsRows.Close()

	var materials []entity.GetMaterialResponse

	// Iterasi untuk setiap material
	for materialsRows.Next() {
		var material entity.GetMaterialResponse
		err := materialsRows.Scan(&material.Id, &material.Title)
		if err != nil {
			log.Error().Err(err).Msg("repo::GetAllSyllabus - Failed to scan material data")
			return nil, err
		}

		modulesQuery := `
            SELECT 
                id, 
                title, 
                content, 
                attachments, 
                videos 
            FROM 
                modules 
            WHERE 
                materials_id = $1
        `
		modulesRows, err := r.db.QueryContext(ctx, modulesQuery, material.Id)
		if err != nil {
			log.Error().Err(err).Int("material_id", material.Id).Msg("repo::GetAllSyllabus - Failed to query modules")
			return nil, err
		}

		var modules []entity.GetModuleResponse

		for modulesRows.Next() {
			var module entity.GetModuleResponse
			var attachments, videos []string

			err := modulesRows.Scan(
				&module.Id,
				&module.Title,
				&module.Content,
				pq.Array(&attachments),
				pq.Array(&videos),
			)
			if err != nil {
				log.Error().Err(err).Msg("repo::GetAllSyllabus - Failed to scan module data")
				return nil, err
			}

			module.Attachments = attachments
			module.Videos = videos
			modules = append(modules, module)
		}

		modulesRows.Close()

		material.Modules = modules
		materials = append(materials, material)
	}

	if err := materialsRows.Err(); err != nil {
		log.Error().Err(err).Msg("repo::GetAllSyllabus - Error occurred during materials iteration")
		return nil, err
	}

	return materials, nil
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
		return errmsg.NewCustomErrors(404, errmsg.WithMessage("Class not found or unauthorized access"))
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

func (r *classRepository) GetAllUsersEnrolledClass(ctx context.Context, req *entity.GetAllUsersEnrolledClassRequest) (*entity.GetAllUsersEnrolledClassResponse, error) {
	var response entity.GetAllUsersEnrolledClassResponse
	var users []entity.GetUsersEnrolledClassResponse
	var total int // Counter for rows

	query := `
		SELECT  
			u.id AS user_id,
			u.name AS name
		FROM 
			users AS u
		INNER JOIN 
			users_classes AS ce ON u.id = ce.user_id
		INNER JOIN 
			class AS c ON ce.class_id = c.id
		WHERE 
			ce.class_id = $1 AND c.creator_class_id = $2;
	`

	rows, err := r.db.QueryContext(ctx, query, req.ClassId, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllUsersEnrolledClass - Failed to fetch enrolled users")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.GetUsersEnrolledClassResponse
		if err := rows.Scan(&user.UserId, &user.Name); err != nil {
			log.Error().Err(err).Msg("repo::GetAllUsersEnrolledClass - Failed to scan user data")
			return nil, err
		}
		users = append(users, user)
		total++
	}

	if rows.Err() != nil {
		log.Error().Err(rows.Err()).Msg("repo::GetAllUsersEnrolledClass - Error while iterating rows")
		return nil, rows.Err()
	}

	response.UsersEnrolled = users
	response.Total = total

	return &response, nil
}

func (r *classRepository) DeleteStudentClass(ctx context.Context, req *entity.DeleteUsersClassRequest) error {
	query := `
		DELETE FROM users_classes
		WHERE class_id = $1 AND user_id = $2 AND EXISTS (
			SELECT 1
			FROM class
			WHERE id = $1 AND creator_class_id = $3
		)
	`

	result, err := r.db.ExecContext(ctx, query, req.ClassId, req.StudentId, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::DeleteStudentClass - Failed to delete student from class")
		return errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to delete student from class"))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::DeleteStudentClass - Failed to get rows affected")
		return errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to delete student from class"))
	}

	if rowsAffected == 0 {
		log.Warn().Any("payload", req).Msg("repo::DeleteStudentClass - No rows affected, unauthorized or invalid class/student")
		return errmsg.NewCustomErrors(404, errmsg.WithMessage("Class or student not found, or you are not authorized"))
	}

	return nil
}

func (r *classRepository) GetAllStudentNotEnrolledClass(ctx context.Context, req *entity.GetAllUserNotEnrolledClassRequest) (*entity.GetAllUserNotEnrolledClassResponse, error) {
	query := `
        SELECT 
            u.name,
            u.email,
            u.image_url
        FROM 
            users u
        LEFT JOIN 
            users_classes uc ON u.id = uc.user_id AND uc.class_id = $1
        WHERE 
            uc.user_id IS NULL;
    `

	var students []entity.GetUserNotEnrolledClassResponse

	rows, err := r.db.QueryContext(ctx, query, req.ClassId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllStudentsNotEnrolledClass - Failed to get all student not enrolled class")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}
	defer rows.Close()

	for rows.Next() {
		var student entity.GetUserNotEnrolledClassResponse
		if err := rows.Scan(&student.Name, &student.Email, &student.ImageUrl); err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::GetAllStudentsNotEnrolledClass - Failed to get all student not enrolled class")
			return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllStudentsNotEnrolledClass - Failed to get all student not enrolled class")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal server error"))
	}

	response := &entity.GetAllUserNotEnrolledClassResponse{
		Total:    len(students),
		Students: students,
	}

	return response, nil
}

func (r *classRepository) CheckEnrollment(ctx context.Context, req *entity.AddUsersToClassRequest) error {
	query := `
        SELECT COUNT(*) 
        FROM users_classes 
        WHERE user_id = $1 AND class_id = $2
    `

	var count int

	err := r.db.QueryRowContext(ctx, query, req.StudentId, req.ClassId).Scan(&count)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllStudentsNotEnrolledClass - Failed to get all student not enrolled class")
		return errmsg.NewCustomErrors(400, errmsg.WithMessage("Id user not valid"))
	}

	if count > 0 {
		return errmsg.NewCustomErrors(400, errmsg.WithMessage("User is already enrolled in this class"))
	}

	return nil
}

func (r *classRepository) AddUserToClass(ctx context.Context, req *entity.AddUsersToClassRequest) (*entity.AddUsersToClassResponse, error) {
	query := `
        INSERT INTO users_classes (user_id, class_id, enrollment_status, created_at, updated_at)
        VALUES ($1, $2, 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, user_id, class_id, 'active' AS enrollment_status, created_at, updated_at;
    `

	var response entity.AddUsersToClassResponse

	err := r.db.QueryRowContext(ctx, query, req.StudentId, req.ClassId).Scan(
		&response.Id,
		&response.StudentId,
		&response.ClassId,
		&response.StatusEnrollment,
		&response.CreatedAt,
		&response.UpdatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllStudentsNotEnrolledClass - Failed to get all student not enrolled class")
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Id user not valid"))
	}

	return &response, nil
}

func (r *classRepository) TrackModule(ctx context.Context, req *entity.TrackModuleRequest) (*entity.TrackModuleResponse, error) {
	query := `
		INSERT INTO users_progress (
			user_id, class_id, material_id, module_id, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW(), NOW()
		)
		RETURNING id, user_id, progress, status, created_at, updated_at;
	`
	req.StatusProgress = "done"

	// Eksekusi query
	row := r.db.QueryRowContext(ctx, query,
		req.UserId,
		req.ClassId,
		// req.UsersClassesId,
		req.MaterialId,
		req.ModuleId,
		// req.QuizId,
		req.StatusProgress,
		// req.Progress,
	)

	// Menyiapkan respons
	var res entity.TrackModuleResponse
	err := row.Scan(
		&res.Id,
		&res.UserId,
		&res.Progress,
		&res.StatusProgress,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Error().Err(err).Any("payload", req).Msg("repo::TrackModule - Failed to track module")
			return nil, err
		}

		switch pqErr.Code.Name() {
		case "foreign_key_violation":
			log.Warn().Msg("repo::TrackModule - Class, Materials, Module, or user with the ID not found")
			return nil, errmsg.NewCustomErrors(409, errmsg.WithMessage("Invalid class ID, material ID, module ID, or user ID"))
		default:
			log.Error().Err(err).Any("payload", req).Msg("repo::TrackModule - Unhandled pq.Error")
			return nil, err
		}
	}

	progressReq := &entity.GetProgressRequest{
		UserId:  req.UserId,
		ClassId: req.ClassId,
	}

	progress, err := r.GetProgress(ctx, progressReq)
	if err != nil {
		log.Error().Err(err).Any("payload", progressReq).Msg("repo::TrackModule - Failed to calculate progress")
		return nil, err
	}

	updateQuery := `
		UPDATE users_progress
		SET progress = $1, updated_at = NOW()
		WHERE user_id = $2 AND class_id = $3;
	`
	_, err = r.db.ExecContext(ctx, updateQuery, progress, req.UserId, req.ClassId)
	if err != nil {
		log.Error().Err(err).Msg("repo::TrackModule - Failed to update user progress")
		return nil, err
	}

	res.Progress = progress
	return &res, nil
}

func (r *classRepository) GetProgress(ctx context.Context, req *entity.GetProgressRequest) (*float64, error) {
	query := `
		SELECT 
			COUNT(DISTINCT m.id) AS total_modules,
			COUNT(DISTINCT up.module_id) FILTER (WHERE up.status = 'done') AS completed_modules
		FROM 
			class c
		JOIN 
			materials mat ON c.id = mat.class_id
		JOIN 
			modules m ON mat.id = m.materials_id
		LEFT JOIN 
			users_progress up ON m.id = up.module_id AND up.user_id = $1
		WHERE 
			c.id = $2;
	`

	var totalModules, completedModules int
	err := r.db.QueryRowContext(ctx, query, req.UserId, req.ClassId).Scan(&totalModules, &completedModules)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetProgress - Failed to calculate progress")
		return nil, err
	}

	progress := 0.0
	if totalModules > 0 {
		progress = (float64(completedModules) / float64(totalModules)) * 100
	}

	return &progress, nil
}

func (r *classRepository) GetAllClassAdmin(ctx context.Context, req *entity.GetAllClassAdminRequest) (*[]entity.GetAllClassAdminResponse, error) {
	query := `
        SELECT
            c.id,
            c.title,
            c.description AS desc,
            c.status,
            c.created_at,
            c.updated_at,
            COALESCE(m.materials_total, 0) AS materials_total,
            COALESCE(mods.modules_total, 0) AS modules_total,
            COALESCE(uc.student_enrolled_total, 0) AS student_enrolled_total
        FROM
            class c
        LEFT JOIN (
            SELECT
                class_id,
                COUNT(*) AS materials_total
            FROM
                materials
            GROUP BY
                class_id
        ) m ON c.id = m.class_id
        LEFT JOIN (
            SELECT
                mt.class_id,
                COUNT(modules.id) AS modules_total
            FROM
                modules
            JOIN materials mt ON modules.materials_id = mt.id
            GROUP BY
                mt.class_id
        ) mods ON c.id = mods.class_id
        LEFT JOIN (
            SELECT
                class_id,
                COUNT(*) AS student_enrolled_total
            FROM
                users_classes
            GROUP BY
                class_id
        ) uc ON c.id = uc.class_id
        WHERE
            c.creator_class_id = $1
    `

	var classes []entity.GetAllClassAdminResponse

	rows, err := r.db.QueryContext(ctx, query, req.UserId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch class data")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var class entity.GetAllClassAdminResponse
		err := rows.Scan(
			&class.Id,
			&class.Title,
			&class.Desc,
			// &class.Tags,
			&class.Status,
			&class.CreatedAt,
			&class.UpdatedAt,
			&class.MaterialsTotal,
			&class.ModulesTotal,
			&class.StudentEnrolledTotal,
		)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::GetAllClassAdmin - Failed to get all class admin")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Internal server error"))
		}
		classes = append(classes, class)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetAllClassAdmin - Error iterating over class rows")
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Internal server error"))
	}

	return &classes, nil
}
