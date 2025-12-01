package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/entities"
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
 * @param cheatsheet_image formData file true "Cheatsheet Image (WebP format only)"
 * @success 201 {object} map[string]string{"message": "Cheatsheet created successfully"}
 * @failure 400 {object} map[string]string{"error": "Failed to create cheatsheet"}
 * @router /api/cheatsheets [post]
 */
func (h *CheatsheetsHandler) CreateCheatsheet(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(2 << 20); err != nil { // 2MB limit for single file
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	form := c.Request.MultipartForm

	// Get metadata from form field
	metadataStr := form.Value["metadata"]
	if len(metadataStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing metadata field"})
		return
	}

	// Parse JSON metadata
	var req dtos.Cheatsheet
	if err := json.Unmarshal([]byte(metadataStr[0]), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata JSON"})
		return
	}

	// Get the cheatsheet image from the form data
	cheatsheetImage, header, err := c.Request.FormFile("cheatsheet_image")
	if err != nil {
		c.JSON(400, gin.H{"error": "Cheatsheet image is required"})
		return
	}
	defer cheatsheetImage.Close()

	// Validate image type
	if header.Header.Get("Content-Type") != "image/webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only WebP images are allowed"})
		return
	}

	// Create a context with 25 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()

	if err := h.service.CreateCheatsheet(ctx, req, cheatsheetImage, header.Size); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create cheatsheet: %v", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Cheatsheet created successfully"})
}

/**
 * Bulk create cheatsheets
 * @param cheatsheets body []dtos.CreateCheatsheetRequest
 * @param cheatsheet_images formData file true "Cheatsheet Images (WebP format only, max 5 files)"
 * @success 201 {object} map[string]string{"message": "All Cheatsheets created successfully"}
 * @failure 400 {object} map[string]string{"error": "Failed to create cheatsheets"}
 * @router /api/cheatsheets/bulk [post]
 */
func (h *CheatsheetsHandler) BulkCreateCheatsheets(c *gin.Context) {
	// Parse multipart form (limit total to ~10 MB since max 5 files × 1 MB)
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	form := c.Request.MultipartForm

	// Get all files uploaded under the "cheatsheet_images" key
	files := form.File["cheatsheet_images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least 1 image is required"})
		return
	}
	if len(files) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can upload a maximum of 5 cheatsheets at once"})
		return
	}

	// Get the metadata JSON string from form
	metadataStr := form.Value["metadata"]
	if len(metadataStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing metadata field"})
		return
	}

	// Parse the JSON metadata into array of requests
	var reqs []dtos.Cheatsheet
	if err := json.Unmarshal([]byte(metadataStr[0]), &reqs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata JSON"})
		return
	}

	// Validate that number of metadata entries matches number of files
	if len(reqs) != len(files) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatch between metadata entries and uploaded files"})
		return
	}

	// Create a context with 60 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	results := h.service.BulkCreateCheatsheets(ctx, reqs, files)

	c.JSON(http.StatusCreated, gin.H{"results": results})
}

/**
 * Get a cheatsheet by ID
 * @param id path string true "Cheatsheet ID"
 * @success 200 {object} map[string]interface{}{"cheatsheet": repository.Cheatsheet}
 * @failure 400 {object} map[string]string{"error": "ID parameter is required"}
 * @failure 404 {object} map[string]string{"error": "Cheatsheet not found for id: {id}"}
 * @failure 500 {object} map[string]string{"error": "Failed to fetch cheatsheet"}
 * @router /api/cheatsheets/{id} [get]
 */
func (h *CheatsheetsHandler) GetCheatsheetByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
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
 * @router /api/cheatsheets/{slug} [get]
 */
func (h *CheatsheetsHandler) GetCheatsheetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug parameter is required"})
		return
	}

	// Create a context with 10 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
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
 * @router /api/cheatsheets [get]
 */
func (h *CheatsheetsHandler) GetAllCheatsheets(c *gin.Context) {
	// add query params for category, subcategory etc.
	category := c.Query("category")
	subcategory := c.Query("subcategory")
	sortBy := c.Query("sort")

	if sortBy != "" && !entities.Filters[sortBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort parameter"})
		return
	}

	// Create a context with 45 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	cheatsheets, err := h.service.GetAllCheatsheets(ctx, category, subcategory, sortBy)
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
 * @router /api/cheatsheets/{id} [put]
 */
func (h *CheatsheetsHandler) UpdateCheatsheet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(2 << 20); err != nil { // 2MB limit for single file
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	form := c.Request.MultipartForm

	// Get metadata from form field
	metadataStr := form.Value["metadata"]
	if len(metadataStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing metadata field"})
		return
	}

	// Parse JSON metadata
	var req dtos.UpdateCheatsheetRequest
	if err := json.Unmarshal([]byte(metadataStr[0]), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata JSON"})
		return
	}

	// Get the cheatsheet image from the form data
	var cheatsheetImage multipart.File
	cheatsheetImageFile, header, err := c.Request.FormFile("cheatsheet_image")
	if err == nil {
		defer cheatsheetImageFile.Close()

		// Validate image type
		if header.Header.Get("Content-Type") != "image/webp" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only WebP images are allowed"})
			return
		}

		// Validate file size (optional)
		if header.Size > 1<<20 { // 1MB limit
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image file too large (max 1MB)"})
			return
		}

		cheatsheetImage = cheatsheetImageFile
	}

	// Create a context with 45 seconds timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	if err := h.service.UpdateCheatsheet(ctx, id, req, cheatsheetImage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update cheatsheet: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cheatsheet updated successfully"})
}
