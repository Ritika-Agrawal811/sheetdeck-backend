package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func setupAnalyticsRoutes(apiGroup *gin.RouterGroup, services *ServicesContainer) {
	analyticsHandler := handlers.NewAnalyticsHandler(services.AnalyticsService)
	analyticsGroup := apiGroup.Group("/analytics")

	analyticsGroup.GET("/pageviews", analyticsHandler.GetPageviewsStats)
	analyticsGroup.GET("/summary/devices", analyticsHandler.GetDevicesStats)
	analyticsGroup.POST("/pageview", analyticsHandler.RecordPageView)
	analyticsGroup.POST("/event", analyticsHandler.RecordEvent)

}
