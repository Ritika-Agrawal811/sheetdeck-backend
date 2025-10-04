package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
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
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&pgbouncer=true",
		utils.GetEnv("POSTGRES_USER", "user"),
		utils.GetEnv("POSTGRES_PASSWORD", "password"),
		utils.GetEnv("POSTGRES_HOST", "localhost"),
		utils.GetEnv("POSTGRES_PORT", "6543"),
		utils.GetEnv("POSTGRES_DB", "postgres"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal().Msgf("Unable to parse DSN: %v", err)
	}

	// 🚨 Important: Disable statement cache to avoid "already exists" errors
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Set the connection pool configuration
	config.MaxConns = 1
	config.MinConns = 0
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 5
	config.HealthCheckPeriod = time.Minute

	// 🔥 Add connection timeout
	config.ConnConfig.ConnectTimeout = 15 * time.Second

	// Retry logic with exponential backoff
	var pool *pgxpool.Pool
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Warn().Msgf("Attempt %d: Unable to create PostgreSQL pool: %v", i+1, err)
			if i == maxRetries-1 {
				log.Fatal().Msgf("Failed to create PostgreSQL pool after %d attempts: %v", maxRetries, err)
			}
			time.Sleep(time.Second * time.Duration(2*(i+1)))
			continue
		}

		// Check the connection with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = pool.Ping(ctx)
		cancel()

		if err != nil {
			log.Warn().Msgf("Attempt %d: Failed to ping PostgreSQL: %v", i+1, err)
			pool.Close()
			if i == maxRetries-1 {
				log.Fatal().Msgf("Failed to connect to PostgreSQL after %d attempts: %v", maxRetries, err)
			}
			time.Sleep(time.Second * time.Duration(2*(i+1)))
			continue
		}
		// Connection successful
		log.Info().Msgf("Successfully connected to PostgreSQL on attempt %d", i+1)
		break
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
