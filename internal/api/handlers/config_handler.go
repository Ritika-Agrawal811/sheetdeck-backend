package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/config"
	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	service config.ConfigService
}

func NewConfigHandler(service config.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		service: service,
	}
}

/**
 * Get Config - categories, subcategories, etc.
 * @success 200 {object} *dtos.ConfigResponse
 * @failure 500 {object} map[string]string
 * @router /config [get]
 */
func (h *ConfigHandler) GetConfiq(c *gin.Context) {
	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	results, err := h.service.GetConfig(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch config : %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, results)
}

/**
 * Get Usage Statistics
 * @success 200 {object} *dtos.UsageResponse
 * @failure 500 {object} map[string]string
 * @router /config/usage [get]
 */
func (h *ConfigHandler) GetUsage(c *gin.Context) {
	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	results, err := h.service.GetUsage(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch usage info : %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, results)
}
