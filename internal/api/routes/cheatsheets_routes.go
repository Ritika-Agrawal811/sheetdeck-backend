package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

/**
 * Handles all the cheatsheets routes
 * @param apiGroup *gin.RouterGroup, services *ServivesContainer
 */
func setupCheatsheetsRoutes(apiGroup *gin.RouterGroup, services *ServicesContainer) {
	cheatsheetsHandler := handlers.NewCheatsheetsHandler(services.CheatsheetsService)
	cheatsheetsGroup := apiGroup.Group("/cheatsheets")

	cheatsheetsGroup.GET("", cheatsheetsHandler.GetAllCheatsheets)
	cheatsheetsGroup.GET("/:id", cheatsheetsHandler.GetCheatsheetByID)
	cheatsheetsGroup.GET("/slug/:slug", cheatsheetsHandler.GetCheatsheetBySlug)
	cheatsheetsGroup.POST("", cheatsheetsHandler.CreateCheatsheet)
	cheatsheetsGroup.POST("/bulk", cheatsheetsHandler.BulkCreateCheatsheets)
	cheatsheetsGroup.PUT("/:id", cheatsheetsHandler.UpdateCheatsheet)

}
