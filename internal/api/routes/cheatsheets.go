package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func setupCheatsheetsRoutes(router *gin.Engine, services *ServicesContainer) {
	cheatsheetsHandler := handlers.NewCheatsheetsHandler(services.CheatsheetsService)
	cheatsheetsGroup := router.Group("/cheatsheets")

	cheatsheetsGroup.GET("/", cheatsheetsHandler.GetAllCheatsheets)
	cheatsheetsGroup.GET("/:id", cheatsheetsHandler.GetCheatsheetByID)
	cheatsheetsGroup.GET("/slug/:slug", cheatsheetsHandler.GetCheatsheetBySlug)
	cheatsheetsGroup.POST("/create", cheatsheetsHandler.CreateCheatsheet)
	cheatsheetsGroup.POST("/create/bulk", cheatsheetsHandler.BulkCreateCheatsheets)
	cheatsheetsGroup.PUT("/update/:id", cheatsheetsHandler.UpdateCheatsheet)

}
