package routes

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/cheatsheets"
	"github.com/gin-gonic/gin"
)

type ServicesContainer struct {
	CheatsheetsService cheatsheets.CheatsheetsService
}

func NewServicesContainer(repo *repository.Queries) *ServicesContainer {
	cheatsheetsService := cheatsheets.NewCheatsheetsService(repo)

	return &ServicesContainer{
		CheatsheetsService: cheatsheetsService,
	}
}

func SetupRoutes(router *gin.Engine, services *ServicesContainer) {
	setupCheatsheetsRoutes(router, services)
}
