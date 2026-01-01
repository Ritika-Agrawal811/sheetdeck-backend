package config

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/entities"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
)

type ConfigService interface {
	GetConfig(ctx context.Context) (*dtos.ConfigResponse, error)
	GetUsage(ctx context.Context) (*dtos.UsageResponse, error)
}

type configService struct {
	repo repository.Querier
	db   repository.DBTX
}

func NewConfigService(repo repository.Querier, db repository.DBTX) ConfigService {

	return &configService{
		repo: repo,
		db:   db,
	}
}

/**
 * Fetches Config - categories, subcategories, stats, etc.
 * @param ctx context.Context
 * @return *dtos.ConfigResponse, error
 */
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

	// Get traffic data
	traffic, err := s.repo.GetTotalViewsAndVisitors(ctx)
	if err != nil {
		return &dtos.ConfigResponse{}, fmt.Errorf("failed to fetch total view and visitors")
	}

	analytics, err := s.repo.GetTotalClicksAndDownloads(ctx)
	if err != nil {
		return &dtos.ConfigResponse{}, fmt.Errorf("failed to fetch total clicks and downloads")
	}

	stats := dtos.GlobalStats{
		TotalCheatsheets:    totalCount,
		TotalViews:          traffic.TotalViews,
		TotalUniqueVisitors: traffic.TotalVisitors,
		TotalClicks:         analytics.TotalClicks,
		TotalDownloads:      analytics.TotalDownloads,
	}

	config := &dtos.ConfigResponse{
		Stats:         stats,
		CategoryStats: categoryStats,
		Categories:    categories,
		Subcategories: subcategories,
	}

	return config, nil
}

/**
 * Fetch category statistics including subcategory breakdowns
 * @param ctx context.Context
 * @return []dtos.CategoryStat
 */
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

/**
 * Get Usage - database size, storage size, etc.
 * @param ctx context.Context
 * @return *dtos.UsageResponse, error
 */
func (s *configService) GetUsage(ctx context.Context) (*dtos.UsageResponse, error) {

	/*  Get database usage info */
	databaseUsage, err := s.getDatabaseSize(ctx)
	if err != nil {
		return &dtos.UsageResponse{}, fmt.Errorf("failed to get database usage info")
	}

	databaseDetails, ok := entities.ResourceLimits["database"]
	if !ok {
		return &dtos.UsageResponse{}, fmt.Errorf("database resource limit config not found")
	}

	databaseUsagePercent := (float64(databaseUsage.SizeBytes) / float64(databaseDetails.LimitBytes)) * 100
	if databaseUsagePercent > 100 {
		databaseUsagePercent = 100
	}

	largestTables, err := s.getLargestTables(ctx)
	if err != nil {
		return &dtos.UsageResponse{}, fmt.Errorf("failed to get largest tables usage info")
	}

	/*  Get storage usage info */
	storageUsage, err := s.repo.GetTotalImageSize(ctx)
	if err != nil {
		return &dtos.UsageResponse{}, fmt.Errorf("failed to get storage usage info")
	}

	storageDetails, ok := entities.ResourceLimits["storage"]
	if !ok {
		return &dtos.UsageResponse{}, fmt.Errorf("storage resource limit config not found")
	}

	storageUsagePercent := (float64(storageUsage.TotalSize) / float64(storageDetails.LimitBytes)) * 100
	if storageUsagePercent > 100 {
		storageUsagePercent = 100
	}

	largestFiles, err := s.getLargestFiles(ctx)
	if err != nil {
		return &dtos.UsageResponse{}, fmt.Errorf("failed to get largest tables usage info")
	}

	response := &dtos.UsageResponse{
		Database: dtos.DatabaseUsage{
			SizeBytes:     databaseUsage.SizeBytes,
			SizePretty:    databaseUsage.SizePretty,
			LimitBytes:    databaseDetails.LimitBytes,
			LimitPretty:   databaseDetails.LimitPretty,
			UsagePercent:  math.Round(databaseUsagePercent),
			LargestTables: largestTables,
		},
		Storage: dtos.StorageUsage{
			SizeBytes:    storageUsage.TotalSize,
			SizePretty:   storageUsage.TotalSizePretty,
			LimitBytes:   storageDetails.LimitBytes,
			LimitPretty:  storageDetails.LimitPretty,
			UsagePercent: math.Round(storageUsagePercent),
			Files:        largestFiles,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	return response, nil
}

/**
 * Retrieve the total size of the PostgreSQL database.
 * @param ctx context.Context
 * @return *entities.DatabaseSize, error
 */
func (s *configService) getDatabaseSize(ctx context.Context) (*entities.DatabaseSize, error) {
	query := `
		SELECT 
			sum(pg_database_size(pg_database.datname))::bigint as size_bytes,
			pg_size_pretty(sum(pg_database_size(pg_database.datname))) as size_pretty
		FROM pg_database;
	`

	var result entities.DatabaseSize
	err := s.db.QueryRow(ctx, query).Scan(&result.SizeBytes, &result.SizePretty)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/**
 * Retrieve 4 largest tables by size
 * @param ctx context.Context
 * @return []dtos.TableData, error
 */
func (s *configService) getLargestTables(ctx context.Context) ([]dtos.TableData, error) {
	query := `
		SELECT
			schemaname AS schema_name,
			tablename AS table_name,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
		FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
		LIMIT 4;
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dtos.TableData
	for rows.Next() {
		var table dtos.TableData
		err := rows.Scan(&table.Schema, &table.Name, &table.Size)
		if err != nil {
			return nil, err
		}
		result = append(result, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

/**
 * Retrieve 4 largest files by size
 * @param ctx context.Context
 * @return []dtos.FileData, error
 */
func (s *configService) getLargestFiles(ctx context.Context) ([]dtos.FileData, error) {
	data, err := s.repo.GetLargestCheatsheets(ctx)
	if err != nil {
		return nil, err
	}

	var result []dtos.FileData
	for _, row := range data {
		result = append(result, dtos.FileData{
			Title:    row.Title,
			Category: row.Category,
			Size:     row.Size,
		})
	}

	return result, nil
}
