package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func setupAnalyticsRoutes(apiGroup *gin.RouterGroup, services *ServicesContainer) {
	analyticsHandler := handlers.NewAnalyticsHandler(services.AnalyticsService)
	analyticsGroup := apiGroup.Group("/analytics")

	analyticsGroup.GET("/overview", analyticsHandler.GetMetricsOverview)
	analyticsGroup.GET("/summary/devices", analyticsHandler.GetDevicesStats)
	analyticsGroup.GET("/summary/browsers", analyticsHandler.GetBrowsersStats)
	analyticsGroup.GET("/summary/os", analyticsHandler.GetOperatingSystemsStats)
	analyticsGroup.GET("/summary/referrers", analyticsHandler.GetReferrerStats)
	analyticsGroup.POST("/pageview", analyticsHandler.RecordPageView)
	analyticsGroup.POST("/event", analyticsHandler.RecordEvent)

}
