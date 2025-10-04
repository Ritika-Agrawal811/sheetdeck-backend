package handlers

import (
	"context"
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

func (h *ConfigHandler) GetConfiq(c *gin.Context) {
	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	results, err := h.service.GetConfig(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
