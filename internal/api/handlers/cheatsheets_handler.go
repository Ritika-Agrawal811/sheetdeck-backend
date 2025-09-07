package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/cheatsheets"
	"github.com/gin-gonic/gin"
)

type CheatsheetsHandler struct {
	service cheatsheets.CheatsheetsService
}

func NewCheatsheetsHandler(service cheatsheets.CheatsheetsService) *CheatsheetsHandler {
	return &CheatsheetsHandler{
		service: service,
	}
}

func (h *CheatsheetsHandler) CreateCheatsheet(c *gin.Context) {
	var req dtos.CreateCheatsheetRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a context with 5 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.CreateCheatsheet(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cheatsheet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Cheatsheet created successfully"})

}
