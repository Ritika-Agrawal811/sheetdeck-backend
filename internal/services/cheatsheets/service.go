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

	// Call the repository method to create a cheatsheet
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
