package handlers

import (
	"context"
	"fmt"
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

/**
 * Create a new cheatsheet
 * @param cheatsheet body dtos.CreateCheatsheetRequest
 * @success 201 {object} map[string]string{"message": "Cheatsheet created successfully"}
 * @failure 400 {object} map[string]string{"error": "Failed to create cheatsheet"}
 * @router /cheatsheets/create [post]
 */
func (h *CheatsheetsHandler) CreateCheatsheet(c *gin.Context) {
	var req dtos.CreateCheatsheetRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request : %v", err.Error())})
		return
	}

	// Create a context with 5 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.CreateCheatsheet(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create cheatsheet: %v", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Cheatsheet created successfully"})
}

/**
 * Bulk create cheatsheets
 * @param cheatsheets body []dtos.CreateCheatsheetRequest
 * @success 201 {object} map[string]string{"message": "All Cheatsheets created successfully"}
 * @failure 400 {object} map[string]string{"error": "Failed to create cheatsheets"}
 * @router /cheatsheets/create/bulk [post]
 */
func (h *CheatsheetsHandler) BulkCreateCheatsheets(c *gin.Context) {
	var req dtos.BulkCreateCheatsheetRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request : %v", err.Error())})
		return
	}

	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.BulkCreateCheatsheets(ctx, req.Cheatsheets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create cheatsheets: %v", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "All Cheatsheets created successfully"})
}
