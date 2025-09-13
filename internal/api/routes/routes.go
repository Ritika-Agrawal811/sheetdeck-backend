package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/analytics"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/cheatsheets"
	"github.com/gin-gonic/gin"
)

type ServicesContainer struct {
	CheatsheetsService cheatsheets.CheatsheetsService
	AnalyticsService   analytics.AnalyticsService
}

func NewServicesContainer(repo *repository.Queries) *ServicesContainer {
	cheatsheetsService := cheatsheets.NewCheatsheetsService(repo)
	analyticsService := analytics.NewAnalyticsService(repo)

	return &ServicesContainer{
		CheatsheetsService: cheatsheetsService,
		AnalyticsService:   analyticsService,
	}
}

func SetupRoutes(apiGroup *gin.RouterGroup, services *ServicesContainer) {
	setupCheatsheetsRoutes(apiGroup, services)
	setupAnalyticsRoutes(apiGroup, services)
}
