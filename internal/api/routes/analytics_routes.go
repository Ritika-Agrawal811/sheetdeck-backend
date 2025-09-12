package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func setupAnalyticsRoutes(router *gin.Engine, services *ServicesContainer) {
	analyticsHandler := handlers.NewAnalyticsHandler(services.AnalyticsService)
	analyticsGroup := router.Group("/api/analytics")

	analyticsGroup.POST("/pageview", analyticsHandler.RecordPageView)

}
