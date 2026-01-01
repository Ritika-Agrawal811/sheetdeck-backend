package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/middlewares"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/routes"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/infra/db"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port           string
	AllowedOrigins []string
	Environment    string
	SSLRedirect    bool
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

/**
 * Prepares the configuration details for the server
 * @return *Config the configuration struct
 */
func loadConfig() (*Config, error) {

	/* Loads the .env file */
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found (this may be fine in PROD)")
	}

	env := os.Getenv("ENV")

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	sslRedirect := true
	if env == "TEST" {
		sslRedirect = false
	}

	return &Config{
		Port:           os.Getenv("PORT"),
		Environment:    os.Getenv("ENV"),
		SSLRedirect:    sslRedirect,
		AllowedOrigins: strings.Split(allowedOrigins, ","),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}, nil
}

/**
 * Initialize Gin HTTP server with middleware and routes
 * @param *Config Configuration settings
 * @return *gin.Engine Configured Gin engine, *gin.RouterGroup the gin router group for APIs
 */
func initGin(cfg *Config) (*gin.Engine, *gin.RouterGroup) {

	/**
	 * Sets Gin to release mode in production,
	 * which disables debug output and improves performance.
	 */
	if cfg.Environment == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	/* Sets up Logger middleware for the router */
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		/* Skip noisy endpoints */
		if param.Path == "/healthz" {
			return ""
		}

		log.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Str("latency", param.Latency.String()).
			Msg("request")

		return ""
	}))

	/* Sets up Recovery middleware for the router */
	r.Use(gin.Recovery())

	/* Adds security headers for the router */
	r.Use(secure.New(secure.Config{
		SSLRedirect:           cfg.SSLRedirect,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            31536000,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	/* Handles CORS for the requests */
	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie", "Set-Cookie", "X-Requested-With"},
		MaxAge:       12 * time.Hour,
	}))

	/* Creates /heatlhz endpoint to check health status of the server */
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	/* Create API group with middlewares for checking valid requests and rate limiting */
	apiGroup := r.Group("/api")
	apiGroup.Use(middlewares.ValidateRequestMiddleware(cfg.AllowedOrigins))

	rateLimiter := middlewares.NewRateLimiterStore(150, time.Minute)
	apiGroup.Use(rateLimiter.RateLimitMiddleware())

	return r, apiGroup
}

func main() {
	log.Info().Msg("Starting the backend server...")

	/* Loads configuration */
	cfg, err := loadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load configuration")
		os.Exit(1)
	}

	/* Channel to listen for interrupt or terminate signals from the OS. */
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	/* Initializes Gin engine */
	r, apiGroup := initGin(cfg)

	/* Creates HTTP server with timeouts and max header bytes */
	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        r,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	/* Initializes Postgres */
	pgClient := db.InitPostgresClient()
	if pgClient == nil {
		log.Info().Msg("Failed to initialize Postgres client")
		os.Exit(1)
	}

	/* Setup routes, repo and services */
	repo := repository.New(pgClient.Client)
	db := db.GetPostgresClient()
	services := routes.NewServicesContainer(repo, db)

	routes.SetupRoutes(apiGroup, services)

	/* Starts the server in a goroutine and stores the errors in a channel */
	serverErr := make(chan error, 1)

	go func() {
		log.Info().Msgf("Server is running on port %s in %s mode", cfg.Port, cfg.Environment)
		log.Info().Msg("Press Ctrl+C to stop")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	/* Handles interrupt signal or server error */
	select {
	case sig := <-sigChan:
		log.Info().Msgf("Received signal: %s. Shutting down...", sig)
	case err := <-serverErr:
		log.Error().Err(err).Msg("Server encountered an error")
	}

	/* Performs graceful shutdown */
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db.Close()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
		os.Exit(1)
	}

	log.Info().Msg("Server exited gracefully")
	os.Exit(0)

}
