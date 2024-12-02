package adapter

import (
	// "log"

	"hacko-app/internal/infrastructure/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func WithHackoPostgres() Option {
	return func(a *Adapter) {
		dbUser := config.Envs.HackoPostgres.Username
		dbPassword := config.Envs.HackoPostgres.Password
		dbName := config.Envs.HackoPostgres.Database
		dbHost := config.Envs.HackoPostgres.Host
		dbSSLMode := config.Envs.HackoPostgres.SslMode
		dbPort := config.Envs.HackoPostgres.Port

		dbMaxPoolSize := config.Envs.DB.MaxOpenCons
		dbMaxIdleConns := config.Envs.DB.MaxIdleCons
		dbConnMaxLifetime := config.Envs.DB.ConnMaxLifetime

		connectionString := "user=" + dbUser + " password=" + dbPassword + " host=" + dbHost + " port=" + dbPort + " dbname=" + dbName + " sslmode=" + dbSSLMode + " TimeZone=UTC"
		db, err := sqlx.Connect("postgres", connectionString)
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Postgres")
		}

		db.SetMaxOpenConns(dbMaxPoolSize)
		db.SetMaxIdleConns(dbMaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConnMaxLifetime) * time.Second)

		// check connection
		err = db.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Hacko Postgres")
		}

		a.HackoPostgres = db
		log.Info().Msg("Hacko Postgres connected")
	}
}
