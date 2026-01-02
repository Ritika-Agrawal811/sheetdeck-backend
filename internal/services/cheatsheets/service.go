package cheatsheets

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/entities"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/storage"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
)

type CheatsheetsService interface {
	CreateCheatsheet(ctx context.Context, details dtos.Cheatsheet, image multipart.File, imageSize int64) error
	BulkCreateCheatsheets(ctx context.Context, details []dtos.Cheatsheet, files []*multipart.FileHeader) []string
	GetCheatsheetByID(ctx context.Context, id string) (*repository.GetCheatsheetByIDRow, error)
	GetCheatsheetBySlug(ctx context.Context, slug string) (*repository.GetCheatsheetBySlugRow, error)
	GetAllCheatsheets(ctx context.Context, category string, subcategory string, sortBy string, limit string) ([]repository.ListCheatsheetsRow, error)
	UpdateCheatsheet(ctx context.Context, id string, details dtos.UpdateCheatsheetRequest, image multipart.File) error
}

type cheatsheetsService struct {
	repo       repository.Querier
	storageSdk storage.StorageService
}

func NewCheatsheetsService(repo repository.Querier, storageSdk storage.StorageService) CheatsheetsService {
	return &cheatsheetsService{
		repo:       repo,
		storageSdk: storageSdk,
	}
}

/**
 * Create a new cheatsheet
 * @param details dtos.Cheatsheet
 * @return error
 */
func (s *cheatsheetsService) CreateCheatsheet(ctx context.Context, details dtos.Cheatsheet, image multipart.File, imageSize int64) error {
	// Validate required fields
	if details.Slug == "" || details.Title == "" || details.Category == "" || details.SubCategory == "" {
		return fmt.Errorf("missing required fields")
	}

	cheatsheetDetails := repository.CreateCheatsheetParams{
		Slug:        details.Slug,
		Title:       details.Title,
		Category:    repository.Category(details.Category),
		Subcategory: repository.Subcategory(details.SubCategory),
		ImageSize:   utils.PgInt8(imageSize),
	}

	// Upload image to Supabase Storage
	filePath := &entities.FilePaths{
		NewPath: fmt.Sprintf("%s/%s/%s.webp", details.Category, details.SubCategory, details.Slug),
	}

	imageUrl, err := s.uploadCheatsheetImage(image, filePath, "upload")
	if err != nil {
		return fmt.Errorf("failed to upload the cheat sheet in storage: %w", err)
	}

	// Set the image URL
	cheatsheetDetails.ImageUrl = utils.PgText(imageUrl)

	if err := s.repo.CreateCheatsheet(ctx, cheatsheetDetails); err != nil {
		return fmt.Errorf("failed to create cheatsheet: %w", err)
	}

	return nil
}

/**
 * Bulk create cheatsheets
 * @param details []dtos.Cheatsheet
 * @param files []*multipart.FileHeader
 * @return error
 */
func (s *cheatsheetsService) BulkCreateCheatsheets(ctx context.Context, details []dtos.Cheatsheet, files []*multipart.FileHeader) []string {

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
		if err := s.CreateCheatsheet(ctx, details[i], file, fileHeader.Size); err != nil {
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
func (s *cheatsheetsService) GetCheatsheetByID(ctx context.Context, id string) (*repository.GetCheatsheetByIDRow, error) {
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
func (s *cheatsheetsService) GetCheatsheetBySlug(ctx context.Context, slug string) (*repository.GetCheatsheetBySlugRow, error) {
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
func (s *cheatsheetsService) GetAllCheatsheets(ctx context.Context, category string, subcategory string, sortBy string, limit string) ([]repository.ListCheatsheetsRow, error) {
	defaultLimit := 15
	defaultSortBy := "recent"

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil && parsedLimit > 0 {
			defaultLimit = parsedLimit
		}
	}

	if sortBy != "" && entities.Filters[sortBy] {
		defaultSortBy = sortBy
	}

	details := repository.ListCheatsheetsParams{
		Category:    repository.NullCategory{Category: repository.Category(category), Valid: category != ""},
		Subcategory: repository.NullSubcategory{Subcategory: repository.Subcategory(subcategory), Valid: subcategory != ""},
		Limit:       int32(defaultLimit),
		Offset:      0,
		SortBy:      defaultSortBy,
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
 * @param details dtos.Cheatsheet
 * @return error
 */
func (s *cheatsheetsService) UpdateCheatsheet(ctx context.Context, id string, details dtos.UpdateCheatsheetRequest, image multipart.File) error {
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
	}

	// Only upload image is provided
	if image != nil {
		// Upload image to Supabase Storage
		filePaths, err := s.createFilePathForUpdate(ctx, id, details.Category, details.SubCategory, details.Slug)
		if err != nil {
			return err
		}

		uploadType := "replace"
		if strings.EqualFold(filePaths.NewPath, filePaths.OldPath) {
			uploadType = "update"
		}

		imageUrl, err := s.uploadCheatsheetImage(image, filePaths, uploadType)
		if err != nil {
			return err
		}

		// Set the image URL in update params
		updateParams.ImageUrl = imageUrl
	}

	if err := s.repo.UpdateCheatsheet(ctx, updateParams); err != nil {
		return fmt.Errorf("failed to update cheatsheet: %w", err)
	}

	return nil
}

/**
 * Uploads a cheatsheet image to Supabase Storage
 * @param category string
 * @param subCategory string
 * @param slug string
 * @param image multipart.File
 * @return string (image URL), error
 */
func (s *cheatsheetsService) uploadCheatsheetImage(image multipart.File, filePaths *entities.FilePaths, uploadType string) (string, error) {
	type UploadResult struct {
		imageUrl string
		err      error
	}

	uploadChan := make(chan UploadResult, 1)

	// Upload image to Supabase Storage concurrently
	go func() {

		var imageUrl string
		var err error

		switch uploadType {
		case "upload":
			imageUrl, err = s.storageSdk.UploadFile(filePaths.NewPath, image)
		case "update":
			imageUrl, err = s.storageSdk.UpdateFileInPlace(filePaths.OldPath, image)
		case "replace":
			imageUrl, err = s.storageSdk.ReplaceFile(filePaths.OldPath, filePaths.NewPath, image)
		}
		uploadChan <- UploadResult{imageUrl, err}
	}()

	// Wait for upload to complete
	uploadRes := <-uploadChan
	if uploadRes.err != nil {
		return "", fmt.Errorf("failed to upload cheatsheet image: %w", uploadRes.err)
	}

	return uploadRes.imageUrl, nil
}

/**
 * Creates new and previous file paths for cheatsheet image for database storage
 * @param id string - cheatsheet id
 * @param category string
 * @param subCategory string
 * @param slug string
 * @return *entities.FilePaths, error
 */
func (s *cheatsheetsService) createFilePathForUpdate(ctx context.Context, id string, category, subCategory, slug string) (*entities.FilePaths, error) {
	uuid, err := utils.StringToUUID(id)
	if err != nil {
		return &entities.FilePaths{}, fmt.Errorf("failed to convert id to uuid")
	}

	cheatsheetDetails, err := s.repo.GetCheatsheetByID(ctx, uuid)
	if err != nil {
		return &entities.FilePaths{}, fmt.Errorf("failed to fetch previous cheat sheet details")
	}

	// Use provided values or fall back to database values
	if category == "" {
		category = string(cheatsheetDetails.Category)
	}

	if subCategory == "" {
		subCategory = string(cheatsheetDetails.Subcategory)
	}

	if slug == "" {
		slug = cheatsheetDetails.Slug
	}

	fileName := fmt.Sprintf("%s/%s/%s.webp", category, subCategory, slug)
	prevFileName := fmt.Sprintf("%s/%s/%s.webp", cheatsheetDetails.Category, cheatsheetDetails.Subcategory, cheatsheetDetails.Slug)

	results := &entities.FilePaths{
		NewPath: fileName,
		OldPath: prevFileName,
	}

	return results, nil
}
