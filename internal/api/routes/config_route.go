package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func setupConfigRoutes(apiGroup *gin.RouterGroup, services *ServicesContainer) {
	configHandler := handlers.NewConfigHandler(services.ConfigService)
	configGroup := apiGroup.Group("/config")

	configGroup.GET("", configHandler.GetConfiq)

}
