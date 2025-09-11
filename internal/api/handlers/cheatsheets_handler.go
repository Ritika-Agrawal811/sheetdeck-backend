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

/**
 * Get a cheatsheet by ID
 * @param id path string true "Cheatsheet ID"
 * @success 200 {object} map[string]interface{}{"cheatsheet": repository.Cheatsheet}
 * @failure 400 {object} map[string]string{"error": "ID parameter is required"}
 * @failure 404 {object} map[string]string{"error": "Cheatsheet not found for id: {id}"}
 * @failure 500 {object} map[string]string{"error": "Failed to fetch cheatsheet"}
 * @router /cheatsheets/{id} [get]
 */
func (h *CheatsheetsHandler) GetCheatsheetByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Create a context with 5 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cheatsheet, err := h.service.GetCheatsheetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch cheatsheet: %v", err.Error())})
		return
	}

	if cheatsheet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cheatsheet not found for id: %s", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cheatsheet": cheatsheet})
}

/**
 * Get a cheatsheet by slug
 * @param slug path string true "Cheatsheet Slug"
 * @success 200 {object} map[string]interface{}{"cheatsheet": repository.Cheatsheet}
 * @failure 400 {object} map[string]string{"error": "Slug parameter is required"}
 * @failure 404 {object} map[string]string{"error": "Cheatsheet not found for slug: {slug}"}
 * @failure 500 {object} map[string]string{"error": "Failed to fetch cheatsheet"}
 * @router /cheatsheets/{slug} [get]
 */
func (h *CheatsheetsHandler) GetCheatsheetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug parameter is required"})
		return
	}

	// Create a context with 5 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cheatsheet, err := h.service.GetCheatsheetBySlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch cheatsheet: %v", err.Error())})
		return
	}

	if cheatsheet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cheatsheet not found for slug: %s", slug)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cheatsheet": cheatsheet})
}

/**
 * Get all cheatsheets with optional filters
 * @param category query string false "Category filter"
 * @param subcategory query string false "Subcategory filter"
 * @success 200 {object} map[string]interface{}{"cheatsheets": []repository.Cheatsheet}
 * @failure 500 {object} map[string]string{"error": "Failed to fetch cheatsheets"}
 * @router /cheatsheets [get]
 */
func (h *CheatsheetsHandler) GetAllCheatsheets(c *gin.Context) {
	// add query params for category, subcategory etc.
	category := c.Query("category")
	subcategory := c.Query("subcategory")

	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	cheatsheets, err := h.service.GetAllCheatsheets(ctx, category, subcategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch cheatsheets: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cheatsheets": cheatsheets})
}

/**
 * Update a cheatsheet by ID
 * @param id path string true "Cheatsheet ID"
 * @param cheatsheet body dtos.UpdateCheatsheetRequest
 * @success 200 {object} map[string]string{"message": "Cheatsheet updated successfully"}
 * @failure 400 {object} map[string]string{"error": "ID parameter is required"}
 * @failure 500 {object} map[string]string{"error": "Failed to update cheatsheet"}
 * @router /cheatsheets/update/{id} [put]
 */
func (h *CheatsheetsHandler) UpdateCheatsheet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var req dtos.UpdateCheatsheetRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request : %v", err.Error())})
		return
	}

	// Create a context with 5 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.UpdateCheatsheet(ctx, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update cheatsheet: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cheatsheet updated successfully"})
}
