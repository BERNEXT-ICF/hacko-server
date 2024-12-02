package seeds

import (
	"context"
	"hacko-app/internal/adapter"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	// "github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// Run seeds.
func (s *Seed) run(table string, total int) {

	switch table {
	// case "roles":
	// 	s.rolesSeed()
	case "users":
		s.usersSeed(total)
	case "all":
		// s.rolesSeed()
		s.usersSeed(total)
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}

	if table != "" {
		log.Info().Msg("Seed ran successfully")
		log.Info().Msg("Exiting ...")
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
		os.Exit(0)
	}
}

func (s *Seed) deleteAll() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Msg("Error committing transaction")
			}
		}
	}()

	_, err = tx.Exec(`DELETE FROM users`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting users")
		return
	}
	log.Info().Msg("users table deleted successfully")

	_, err = tx.Exec(`DELETE FROM roles`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting roles")
		return
	}
	log.Info().Msg("roles table deleted successfully")

	log.Info().Msg("=== All tables deleted successfully ===")
}

// users
func (s *Seed) usersSeed(total int) {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	var userMaps = make([]map[string]any, 0)

	validRoles := []string{"user", "admin", "teacher"}

	for i := 0; i < total; i++ {
		dataUserToInsert := make(map[string]any)
		dataUserToInsert["id"] = uuid.New().String() // Menghasilkan UUID yang valid

		role := validRoles[gofakeit.Number(0, len(validRoles)-1)]

		dataUserToInsert["role"] = role
		dataUserToInsert["name"] = gofakeit.Name()
		dataUserToInsert["email"] = gofakeit.Email()
		dataUserToInsert["password"] = "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G" // password hash

		userMaps = append(userMaps, dataUserToInsert)
	}

	EndUser := map[string]any{
		"id":       uuid.New().String(),
		"role":     "user",
		"name":     "user",
		"email":    "user@gmail.com",
		"password": "123456789", 
	}

	AdminUser := map[string]any{
		"id":       uuid.New().String(),
		"role":     "admin", // Role yang valid
		"name":     "admin",
		"email":    "admin@gmail.com",
		"password": "123456789", 
	}

	TeacherUser := map[string]any{
		"id":       uuid.New().String(),
		"role":     "teacher", 
		"name":     "teacher",
		"email":    "teacher@gmail.com",
		"password": "123456789",
	}

	userMaps = append(userMaps, EndUser)
	userMaps = append(userMaps, AdminUser)
	userMaps = append(userMaps, TeacherUser)

	_, err = tx.NamedExec(`
		INSERT INTO users (id, role, name, email, password)
		VALUES (:id, :role, :name, :email, :password)
	`, userMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating users")
		return
	}

	log.Info().Msg("users table seeded successfully")
}
