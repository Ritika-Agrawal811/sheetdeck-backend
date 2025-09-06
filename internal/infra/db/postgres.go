package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// PostgresClient is the struct that holds the PostgreSQL client instance.
type PostgresClient struct {
	Client *pgxpool.Pool
}

var (
	pgInstance *PostgresClient
	pgOnce     sync.Once
)

// InitPostgresClient initializes the PostgreSQL client as a singleton.
// This function should be called only once, ideally in the main package.
func InitPostgresClient() *PostgresClient {
	pgOnce.Do(func() {
		pgInstance = &PostgresClient{Client: connectPostgres()}
		log.Info().Msg("Postgres client successfully connected and initialized.")
	})

	return pgInstance
}

// connectPostgres sets up a new connection to the PostgreSQL database.
func connectPostgres() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		utils.GetEnv("POSTGRES_USER", "user"),
		utils.GetEnv("POSTGRES_PASSWORD", "password"),
		utils.GetEnv("POSTGRES_HOST", "localhost"),
		utils.GetEnv("POSTGRES_PORT", "5432"),
		utils.GetEnv("POSTGRES_DB", "dbname"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal().Msgf("Unable to parse DSN: %v", err)
	}

	// Set the connection pool configuration
	config.MaxConns = 10
	config.MinConns = 1
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal().Msgf("Unable to create PostgreSQL pool: %v", err)
	}

	// Check the connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal().Msgf("Failed to connect to PostgreSQL: %v", err)
	}

	return pool
}

// GetPostgresClient returns the PostgreSQL client instance.
// This can be used by other packages to access the database.
func GetPostgresClient() *pgxpool.Pool {
	if pgInstance == nil {
		log.Fatal().Msgf("Postgres client is not initialized. Please call InitPostgresClient first.")
	}
	return pgInstance.Client
}

// Close closes the PostgreSQL client connection.
func Close() {
	if pgInstance != nil {
		pgInstance.Client.Close()
		log.Info().Msg("Postgres client connection closed.")
	}
}
