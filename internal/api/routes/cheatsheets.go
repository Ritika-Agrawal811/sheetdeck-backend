package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func setupCheatsheetsRoutes(router *gin.Engine, services *ServicesContainer) {
	cheatsheetsHandler := handlers.NewCheatsheetsHandler(services.CheatsheetsService)
	cheatsheetsGroup := router.Group("/cheatsheets")

	cheatsheetsGroup.POST("/create", cheatsheetsHandler.CreateCheatsheet)

}
