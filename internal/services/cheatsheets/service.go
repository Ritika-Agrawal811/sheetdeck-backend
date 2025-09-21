package cheatsheets

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/storage"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
)

type CheatsheetsService interface {
	CreateCheatsheet(ctx context.Context, details dtos.CreateCheatsheetRequest, image multipart.File) error
	BulkCreateCheatsheets(ctx context.Context, details []dtos.CreateCheatsheetRequest, files []*multipart.FileHeader) []string
	GetCheatsheetByID(ctx context.Context, id string) (*repository.Cheatsheet, error)
	GetCheatsheetBySlug(ctx context.Context, slug string) (*repository.Cheatsheet, error)
	GetAllCheatsheets(ctx context.Context, category string, subcategory string) ([]repository.Cheatsheet, error)
	UpdateCheatsheet(ctx context.Context, id string, details dtos.UpdateCheatsheetRequest) error
}

type cheatsheetsService struct {
	repo       *repository.Queries
	storageSdk *storage.StorageSdk
}

func NewCheatsheetsService(repo *repository.Queries) CheatsheetsService {
	storageSdk, err := storage.NewStorageSdk()
	if err != nil {
		log.Fatal("Warning: Storage SDK not configured:", err)
	}

	return &cheatsheetsService{
		repo:       repo,
		storageSdk: storageSdk,
	}
}

/**
 * Create a new cheatsheet
 * @param details dtos.CreateCheatsheetRequest
 * @return error
 */
func (s *cheatsheetsService) CreateCheatsheet(ctx context.Context, details dtos.CreateCheatsheetRequest, image multipart.File) error {
	// Validate required fields
	if details.Slug == "" || details.Title == "" || details.Category == "" || details.SubCategory == "" {
		return fmt.Errorf("missing required fields")
	}

	// Create channels for concurrent operations
	type UploadResult struct {
		imageUrl string
		err      error
	}

	uploadChan := make(chan UploadResult, 1)

	// Upload image to Supabase Storage
	go func() {
		fileName := fmt.Sprintf("%s/%s/%s.webp", details.Category, details.SubCategory, details.Slug)
		imageUrl, err := s.storageSdk.UploadFile(fileName, image)
		uploadChan <- UploadResult{imageUrl, err}
	}()

	cheatsheetDetails := repository.CreateCheatsheetParams{
		Slug:        details.Slug,
		Title:       details.Title,
		Category:    repository.Category(details.Category),
		Subcategory: repository.Subcategory(details.SubCategory),
	}

	// Wait for upload to complete
	uploadRes := <-uploadChan
	if uploadRes.err != nil {
		return fmt.Errorf("failed to upload cheatsheet image: %w", uploadRes.err)
	}

	// Set the image URL
	cheatsheetDetails.ImageUrl = utils.PgText(uploadRes.imageUrl)

	if err := s.repo.CreateCheatsheet(ctx, cheatsheetDetails); err != nil {
		return fmt.Errorf("failed to create cheatsheet: %w", err)
	}

	return nil
}

/**
 * Bulk create cheatsheets
 * @param details []dtos.CreateCheatsheetRequest
 * @param files []*multipart.FileHeader
 * @return error
 */
func (s *cheatsheetsService) BulkCreateCheatsheets(ctx context.Context, details []dtos.CreateCheatsheetRequest, files []*multipart.FileHeader) []string {

	// Process each file with its corresponding metadata
	results := make([]string, 0, len(files))

	for i, fileHeader := range files {
		// Enforce max 1 MB per file
		if fileHeader.Size > 1<<20 {
			results = append(results, fmt.Sprintf("File too large: %s (max 1 MB)", fileHeader.Filename))
			continue
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			results = append(results, fmt.Sprintf("Failed to open file %s", fileHeader.Filename))
			continue
		}
		defer file.Close()

		// Validate file type - only allow WebP
		if fileHeader.Header.Get("Content-Type") != "image/webp" {
			results = append(results, fmt.Sprintf("Invalid file type for %s (only WebP allowed)", fileHeader.Filename))
			continue
		}

		// Call existing single cheatsheet service
		if err := s.CreateCheatsheet(ctx, details[i], file); err != nil {
			results = append(results, fmt.Sprintf("Failed: %s (%v)", details[i].Title, err))
		} else {
			results = append(results, fmt.Sprintf("Success: %s", details[i].Title))
		}
	}

	return results
}

/**
 * Get a cheatsheet by its ID
 * @param id string
 * @return *repository.Cheatsheet, error
 */
func (s *cheatsheetsService) GetCheatsheetByID(ctx context.Context, id string) (*repository.Cheatsheet, error) {
	cheatsheetID, err := utils.StringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	cheatsheet, err := s.repo.GetCheatsheetByID(ctx, cheatsheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cheatsheet: %w", err)
	}

	return &cheatsheet, nil
}

/**
 * Get a cheatsheet by its slug
 * @param slug string
 * @return *repository.Cheatsheet, error
 */
func (s *cheatsheetsService) GetCheatsheetBySlug(ctx context.Context, slug string) (*repository.Cheatsheet, error) {
	cheatsheet, err := s.repo.GetCheatsheetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cheatsheet: %w", err)
	}

	return &cheatsheet, nil
}

/**
 * Get all cheatsheets, optionally filtered by category and subcategory
 * @param category string
 * @param subcategory string
 * @return []repository.Cheatsheet, error
 */
func (s *cheatsheetsService) GetAllCheatsheets(ctx context.Context, category string, subcategory string) ([]repository.Cheatsheet, error) {

	details := repository.ListCheatsheetsParams{
		Category:    repository.NullCategory{Category: repository.Category(category), Valid: category != ""},
		Subcategory: repository.NullSubcategory{Subcategory: repository.Subcategory(subcategory), Valid: subcategory != ""},
		Limit:       100,
		Offset:      0,
	}

	cheatsheets, err := s.repo.ListCheatsheets(ctx, details)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cheatsheets: %w", err)
	}

	return cheatsheets, nil
}

/**
 * Update an existing cheatsheet
 * @param id string
 * @param details dtos.UpdateCheatsheetRequest
 * @return error
 */
func (s *cheatsheetsService) UpdateCheatsheet(ctx context.Context, id string, details dtos.UpdateCheatsheetRequest) error {
	cheatsheetID, err := utils.StringToUUID(id)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}

	updateParams := repository.UpdateCheatsheetParams{
		ID:          cheatsheetID,
		Slug:        details.Slug,
		Title:       details.Title,
		Category:    repository.NullCategory{Category: repository.Category(details.Category), Valid: details.Category != ""},
		Subcategory: repository.NullSubcategory{Subcategory: repository.Subcategory(details.SubCategory), Valid: details.SubCategory != ""},
		ImageUrl:    details.ImageURL,
	}

	if err := s.repo.UpdateCheatsheet(ctx, updateParams); err != nil {
		return fmt.Errorf("failed to update cheatsheet: %w", err)
	}

	return nil
}
