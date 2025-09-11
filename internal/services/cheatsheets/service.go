package cheatsheets

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
)

type CheatsheetsService interface {
	CreateCheatsheet(ctx context.Context, details dtos.CreateCheatsheetRequest) error
	BulkCreateCheatsheets(ctx context.Context, details []dtos.CreateCheatsheetRequest) error
	GetCheatsheetByID(ctx context.Context, id string) (*repository.Cheatsheet, error)
	GetCheatsheetBySlug(ctx context.Context, slug string) (*repository.Cheatsheet, error)
	GetAllCheatsheets(ctx context.Context, category string, subcategory string) ([]repository.Cheatsheet, error)
	UpdateCheatsheet(ctx context.Context, id string, details dtos.UpdateCheatsheetRequest) error
}

type cheatsheetsService struct {
	repo *repository.Queries
}

func NewCheatsheetsService(repo *repository.Queries) CheatsheetsService {
	return &cheatsheetsService{
		repo: repo,
	}
}

/**
 * Create a new cheatsheet
 * @param details dtos.CreateCheatsheetRequest
 * @return error
 */
func (s *cheatsheetsService) CreateCheatsheet(ctx context.Context, details dtos.CreateCheatsheetRequest) error {

	if details.Slug == "" || details.Title == "" {
		return errors.New("slug and title are required")
	}

	if details.Category == "" || details.SubCategory == "" {
		return errors.New("category and subcategory are required")
	}

	cheatsheetDetails := repository.CreateCheatsheetParams{
		Slug:        details.Slug,
		Title:       details.Title,
		Category:    repository.Category(details.Category),
		Subcategory: repository.Subcategory(details.SubCategory),
		ImageUrl:    utils.PgText(details.ImageURL),
	}

	if err := s.repo.CreateCheatsheet(ctx, cheatsheetDetails); err != nil {
		return fmt.Errorf("failed to create cheatsheet: %w", err)
	}

	return nil
}

/**
 * Bulk create cheatsheets
 * @param details []dtos.CreateCheatsheetRequest
 * @return error
 */
func (s *cheatsheetsService) BulkCreateCheatsheets(ctx context.Context, details []dtos.CreateCheatsheetRequest) error {
	for _, cheatsheet := range details {
		if err := s.CreateCheatsheet(ctx, cheatsheet); err != nil {
			return err
		}
	}
	return nil
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
