package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

/**
 * Handles all the config routes
 * @param apiGroup *gin.RouterGroup, services *ServivesContainer
 */
func setupConfigRoutes(apiGroup *gin.RouterGroup, services *ServicesContainer) {
	configHandler := handlers.NewConfigHandler(services.ConfigService)
	configGroup := apiGroup.Group("/config")

	configGroup.GET("", configHandler.GetConfiq)
	configGroup.GET("/usage", configHandler.GetUsage)

}
