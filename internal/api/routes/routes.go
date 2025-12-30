package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/analytics"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/cheatsheets"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/config"
	"github.com/gin-gonic/gin"
)

type ServicesContainer struct {
	CheatsheetsService cheatsheets.CheatsheetsService
	AnalyticsService   analytics.AnalyticsService
	ConfigService      config.ConfigService
}

/**
 * Creates a common services struct for all the services
 * @param repo *repository.Queries
 */
func NewServicesContainer(repo *repository.Queries) *ServicesContainer {
	cheatsheetsService := cheatsheets.NewCheatsheetsService(repo)
	analyticsService := analytics.NewAnalyticsService(repo)
	configService := config.NewConfigService(repo)

	return &ServicesContainer{
		CheatsheetsService: cheatsheetsService,
		AnalyticsService:   analyticsService,
		ConfigService:      configService,
	}
}

/**
 * Sets up routes for all api groups - cheatsheets, config, analytics
 * @param apiGroup *gin.RouterGroup, services *ServicesContainer
 */
func SetupRoutes(apiGroup *gin.RouterGroup, services *ServicesContainer) {
	setupCheatsheetsRoutes(apiGroup, services)
	setupAnalyticsRoutes(apiGroup, services)
	setupConfigRoutes(apiGroup, services)
}
