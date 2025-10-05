package config

import (
	"context"
	"fmt"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
)

type ConfigService interface {
	GetConfig(ctx context.Context) (*dtos.ConfigResponse, error)
}

type configService struct {
	repo *repository.Queries
}

func NewConfigService(repo *repository.Queries) ConfigService {
	return &configService{
		repo: repo,
	}
}

func (s *configService) GetConfig(ctx context.Context) (*dtos.ConfigResponse, error) {

	// Get total cheatsheets count
	totalCount, err := s.repo.GetTotalCheasheetsCount(ctx)
	if err != nil {
		return &dtos.ConfigResponse{}, fmt.Errorf("failed to get total cheat sheet count")
	}

	// Get category stats
	categoryStats := s.fetchCategoryStats(ctx)

	// Get categories list
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return &dtos.ConfigResponse{}, fmt.Errorf("failed to get all categories")
	}

	// Get subcategories list
	subcategories, err := s.repo.GetSubcategories(ctx)
	if err != nil {
		return &dtos.ConfigResponse{}, fmt.Errorf("failed to get all subcategories")
	}

	// Get analytics data
	analytics, err := s.repo.GetTotalViewsAndVisitors(ctx)
	if err != nil {
		return &dtos.ConfigResponse{}, fmt.Errorf("failed to fetch total view and visitors")
	}

	stats := dtos.GlobalStats{
		TotalCheatsheets:    totalCount,
		TotalViews:          analytics.TotalViews,
		TotalUniqueVisitors: analytics.TotalVisitors,
	}

	config := &dtos.ConfigResponse{
		Stats:         stats,
		CategoryStats: categoryStats,
		Categories:    categories,
		Subcategories: subcategories,
	}

	return config, nil
}

func (s *configService) fetchCategoryStats(ctx context.Context) []dtos.CategoryStat {
	// Fetch category totals
	categoryDetails, err := s.repo.GetCategoryDetails(ctx)
	if err != nil {
		return []dtos.CategoryStat{}
	}

	// Fetch category + subcategory totals
	subcategoryCounts, err := s.repo.CountCheatsheetsByCategoryAndSubcategory(ctx)
	if err != nil {
		return []dtos.CategoryStat{}
	}

	subcatMap := make(map[string][]dtos.SubcategoryStat)
	for _, row := range subcategoryCounts {
		subcatMap[string(row.Category)] = append(subcatMap[string(row.Category)], dtos.SubcategoryStat{
			Subcategory: string(row.Subcategory),
			Count:       row.CheatsheetCount,
		})
	}

	// Build CategoryStat for each category
	stats := make([]dtos.CategoryStat, 0, len(categoryDetails))
	for _, cat := range categoryDetails {
		subStats := subcatMap[string(cat.Category)]

		stats = append(stats, dtos.CategoryStat{
			Category:           string(cat.Category),
			TotalCount:         cat.CheatsheetCount,
			Subcategories:      cat.Subcategories,
			SubcategoriesStats: subStats,
		})
	}

	return stats
}
