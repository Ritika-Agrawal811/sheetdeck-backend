package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/services/analytics"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service analytics.AnalyticsService
}

func NewAnalyticsHandler(service analytics.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

/**
 * Record a page view
 * @param c *gin.Context
 * @success 201 {object} map[string]string{"message": "Pageview recorded successfully"}
 * @failure 400 {object} map[string]string{"error": "Failed to record page view"}
 * @router /api/analytics/pageview [post]
 */
func (h *AnalyticsHandler) RecordPageView(c *gin.Context) {
	var req dtos.PageviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request : %v", err.Error())})
		return
	}

	// Extract IP address and User-Agent from the request context
	req.IpAddress = c.ClientIP()
	req.UserAgent = c.Request.UserAgent()

	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.RecordPageView(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to record page view: %v", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pageview recorded successfully"})
}

/**
 * Record an event - click, download etc
 * @param c *gin.Context
 * @success 201 {object} map[string]string{"message": "Event recorded successfully"}
 * @failure 400 {object} map[string]string{"error": "Failed to record event"}
 * @router /api/analytics/event [post]
 */
func (h *AnalyticsHandler) RecordEvent(c *gin.Context) {
	var req dtos.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request : %v", err.Error())})
		return
	}

	// Extract IP address from the request context
	req.IpAddress = c.ClientIP()

	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.RecordEvent(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to record event: %v", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event recorded successfully"})

}

/**
 * Fetch page views stats
 * @param c *gin.Context
 * @success 200 {object} dtos.PageviewStatsResponse
 * @failure 500 {object} map[string]string{"error": "Failed to fetch page view stats"}
 * @router /api/analytics/pageviews [get]
 */
func (h *AnalyticsHandler) GetPageviewsStats(c *gin.Context) {
	period := c.Query("period")

	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period can not be empty"})
		return
	}

	// Create a context with 45 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	stats, err := h.service.GetPageviewsStats(ctx, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch page view stats: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, stats)
}

/**
 * Fetch page views stats
 * @param c *gin.Context
 * @success 200 {object} dtos.DeviceStatsResponse
 * @failure 500 {object} map[string]string{"error": "Failed to fetch device stats"}
 * @router /api/analytics/summary/devices [get]
 */
func (h *AnalyticsHandler) GetDevicesStats(c *gin.Context) {
	period := c.Query("period")

	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period can not be empty"})
		return
	}

	// Create a context with 45 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	stats, err := h.service.GetDeviceStats(ctx, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch devices stats: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, stats)
}

/**
 * Fetch browser stats
 * @param c *gin.Context
 * @success 200 {object} dtos.BrowserStatsResponse
 * @failure 500 {object} map[string]string{"error": "Failed to fetch browser stats"}
 * @router /api/analytics/summary/browsers [get]
 */
func (h *AnalyticsHandler) GetBrowsersStats(c *gin.Context) {
	period := c.Query("period")

	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period can not be empty"})
		return
	}

	// Create a context with 45 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	stats, err := h.service.GetBrowserStats(ctx, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch browsers stats: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, stats)
}

/**
 * Fetch operating systems stats
 * @param c *gin.Context
 * @success 200 {object} dtos.OSStatsResponse
 * @failure 500 {object} map[string]string{"error": "Failed to fetch operating systems stats"}
 * @router /api/analytics/summary/os [get]
 */
func (h *AnalyticsHandler) GetOperatingSystemsStats(c *gin.Context) {
	period := c.Query("period")

	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period can not be empty"})
		return
	}

	// Create a context with 45 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	stats, err := h.service.GetOperatingSystemsStats(ctx, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch operating systems stats: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, stats)
}
