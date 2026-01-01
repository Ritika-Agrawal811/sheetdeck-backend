package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/entities"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/mocks"
	"github.com/jackc/pgx/v5"
)

func TestGetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                       string
		mockTotalCount             int64
		mockTotalCountErr          error
		mockCategoryDetails        []repository.GetCategoryDetailsRow
		mockCategoryDetailsErr     error
		mockSubcategoryCounts      []repository.CountCheatsheetsByCategoryAndSubcategoryRow
		mockSubcategoryCountsErr   error
		mockCategories             []string
		mockCategoriesErr          error
		mockSubcategories          []string
		mockSubcategoriesErr       error
		mockTraffic                repository.GetTotalViewsAndVisitorsRow
		mockTrafficErr             error
		mockAnalytics              repository.GetTotalClicksAndDownloadsRow
		mockAnalyticsErr           error
		expectError                bool
		expectErrorMsg             string
		expectedTotalCheatsheets   int64
		expectedTotalViews         int64
		expectedTotalVisitors      int64
		expectedCategoriesCount    int
		expectedSubcategoriesCount int
		expectedCategoryStatsCount int
	}{
		{
			name:           "successfully fetches complete config",
			mockTotalCount: 100,
			mockCategoryDetails: []repository.GetCategoryDetailsRow{
				{
					Category:        "html",
					CheatsheetCount: 50,
					Subcategories:   []string{"concepts"},
				},
				{
					Category:        "css",
					CheatsheetCount: 50,
					Subcategories:   []string{"attributes"},
				},
			},
			mockSubcategoryCounts: []repository.CountCheatsheetsByCategoryAndSubcategoryRow{
				{
					Category:        "html",
					Subcategory:     "concepts",
					CheatsheetCount: 50,
				},
				{
					Category:        "css",
					Subcategory:     "attributes",
					CheatsheetCount: 50,
				},
			},
			mockCategories:    []string{"html", "css"},
			mockSubcategories: []string{"concepts", "attributes"},
			mockTraffic: repository.GetTotalViewsAndVisitorsRow{
				TotalViews:    5000,
				TotalVisitors: 2500,
			},
			mockAnalytics: repository.GetTotalClicksAndDownloadsRow{
				TotalClicks:    1000,
				TotalDownloads: 500,
			},
			expectError:                false,
			expectedTotalCheatsheets:   100,
			expectedTotalViews:         5000,
			expectedTotalVisitors:      2500,
			expectedCategoriesCount:    2,
			expectedSubcategoriesCount: 2,
			expectedCategoryStatsCount: 2,
		},
		{
			name:              "handles error of fetching total count",
			mockTotalCountErr: fmt.Errorf("database error"),
			expectError:       true,
			expectErrorMsg:    "failed to get total cheat sheet count",
		},
		{
			name:              "handles error of fetching categories",
			mockTotalCount:    100,
			mockCategoriesErr: fmt.Errorf("failed to fetch categories"),
			expectError:       true,
			expectErrorMsg:    "failed to get all categories",
		},
		{
			name:                 "handles error of fetching subcategories",
			mockTotalCount:       100,
			mockCategories:       []string{"html", "css"},
			mockSubcategoriesErr: fmt.Errorf("failed to fetch subcategories"),
			expectError:          true,
			expectErrorMsg:       "failed to get all subcategories",
		},
		{
			name:              "handles error of fetching traffic",
			mockTotalCount:    100,
			mockCategories:    []string{"html", "css"},
			mockSubcategories: []string{"concepts", "attributes"},
			mockTrafficErr:    fmt.Errorf("failed to fetch traffic"),
			expectError:       true,
			expectErrorMsg:    "failed to fetch total view and visitors",
		},
		{
			name:              "handles error of fetching analytics",
			mockTotalCount:    100,
			mockCategories:    []string{"html", "css"},
			mockSubcategories: []string{"concepts", "attributes"},
			mockTraffic: repository.GetTotalViewsAndVisitorsRow{
				TotalViews:    5000,
				TotalVisitors: 2500,
			},
			mockAnalyticsErr: fmt.Errorf("failed to fetch analytics"),
			expectError:      true,
			expectErrorMsg:   "failed to fetch total clicks and downloads",
		},
		{
			name:                  "handles empty category details gracefully",
			mockTotalCount:        10,
			mockCategoryDetails:   []repository.GetCategoryDetailsRow{},
			mockSubcategoryCounts: []repository.CountCheatsheetsByCategoryAndSubcategoryRow{},
			mockCategories:        []string{},
			mockSubcategories:     []string{},
			mockTraffic: repository.GetTotalViewsAndVisitorsRow{
				TotalViews:    100,
				TotalVisitors: 50,
			},
			mockAnalytics: repository.GetTotalClicksAndDownloadsRow{
				TotalClicks:    20,
				TotalDownloads: 10,
			},
			expectError:                false,
			expectedTotalCheatsheets:   10,
			expectedTotalViews:         100,
			expectedTotalVisitors:      50,
			expectedCategoriesCount:    0,
			expectedSubcategoriesCount: 0,
			expectedCategoryStatsCount: 0,
		},
		{
			name:                   "handles category details error in fetchCategoryStats",
			mockTotalCount:         100,
			mockCategoryDetailsErr: fmt.Errorf("db error"),
			mockSubcategoryCounts:  []repository.CountCheatsheetsByCategoryAndSubcategoryRow{},
			mockCategories:         []string{"html"},
			mockSubcategories:      []string{"concepts"},
			mockTraffic: repository.GetTotalViewsAndVisitorsRow{
				TotalViews:    5000,
				TotalVisitors: 2500,
			},
			mockAnalytics: repository.GetTotalClicksAndDownloadsRow{
				TotalClicks:    1000,
				TotalDownloads: 500,
			},
			expectError:                false,
			expectedTotalCheatsheets:   100,
			expectedCategoryStatsCount: 0,
			expectedCategoriesCount:    1,
			expectedSubcategoriesCount: 1,
		},
		{
			name:           "handles subcategory counts error in fetchCategoryStats",
			mockTotalCount: 100,
			mockCategoryDetails: []repository.GetCategoryDetailsRow{
				{
					Category: "html", CheatsheetCount: 100, Subcategories: []string{"concepts"},
				},
			},
			mockSubcategoryCountsErr: fmt.Errorf("db error"),
			mockCategories:           []string{"html"},
			mockSubcategories:        []string{"concepts"},
			mockTraffic: repository.GetTotalViewsAndVisitorsRow{
				TotalViews:    5000,
				TotalVisitors: 2500,
			},
			mockAnalytics: repository.GetTotalClicksAndDownloadsRow{
				TotalClicks:    1000,
				TotalDownloads: 500,
			},
			expectError:                false,
			expectedTotalCheatsheets:   100,
			expectedCategoryStatsCount: 0,
			expectedCategoriesCount:    1,
			expectedSubcategoriesCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockQuerier{
				GetTotalCheasheetsCountFunc: func(ctx context.Context) (int64, error) {
					if tt.mockTotalCountErr != nil {
						return 0, tt.mockTotalCountErr
					}

					return tt.mockTotalCount, nil
				},

				GetCategoryDetailsFunc: func(ctx context.Context) ([]repository.GetCategoryDetailsRow, error) {
					if tt.mockCategoryDetailsErr != nil {
						return nil, tt.mockCategoryDetailsErr
					}

					return tt.mockCategoryDetails, nil
				},

				CountCheatsheetsByCategoryAndSubcategoryFunc: func(ctx context.Context) ([]repository.CountCheatsheetsByCategoryAndSubcategoryRow, error) {
					if tt.mockSubcategoryCountsErr != nil {
						return nil, tt.mockSubcategoryCountsErr
					}

					return tt.mockSubcategoryCounts, nil
				},

				GetCategoriesFunc: func(ctx context.Context) ([]string, error) {
					if tt.mockCategoriesErr != nil {
						return nil, tt.mockCategoriesErr
					}

					return tt.mockCategories, nil
				},

				GetSubcategoriesFunc: func(ctx context.Context) ([]string, error) {
					if tt.mockSubcategoriesErr != nil {
						return nil, tt.mockSubcategoriesErr
					}

					return tt.mockSubcategories, nil
				},

				GetTotalViewsAndVisitorsFunc: func(ctx context.Context) (repository.GetTotalViewsAndVisitorsRow, error) {
					if tt.mockTrafficErr != nil {
						return repository.GetTotalViewsAndVisitorsRow{}, tt.mockTrafficErr
					}

					return tt.mockTraffic, nil
				},

				GetTotalClicksAndDownloadsFunc: func(ctx context.Context) (repository.GetTotalClicksAndDownloadsRow, error) {
					if tt.mockAnalyticsErr != nil {
						return repository.GetTotalClicksAndDownloadsRow{}, tt.mockAnalyticsErr
					}

					return tt.mockAnalytics, nil
				},
			}

			/* Create service with the mock repo */
			service := NewConfigService(mockRepo, nil)

			result, err := service.GetConfig(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tt.expectErrorMsg {
					t.Errorf("Expected error messaged = %q, but got %q", tt.expectErrorMsg, err.Error())
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Stats.TotalCheatsheets != tt.expectedTotalCheatsheets {
				t.Errorf("Expected TotalCheatsheets %d, but got %d", tt.expectedTotalCheatsheets, result.Stats.TotalCheatsheets)
			}

			if tt.expectedTotalViews > 0 && result.Stats.TotalViews != tt.expectedTotalViews {
				t.Errorf("Expected TotalViews %d, but got %d", tt.expectedTotalViews, result.Stats.TotalViews)
			}

			if tt.expectedTotalVisitors > 0 && result.Stats.TotalUniqueVisitors != tt.expectedTotalVisitors {
				t.Errorf("Expected TotalUniqueVisitors %d, but got %d", tt.expectedTotalVisitors, result.Stats.TotalUniqueVisitors)
			}

			if len(result.Categories) != tt.expectedCategoriesCount {
				t.Errorf("Expected %d categories, but got %d", tt.expectedCategoriesCount, len(result.Categories))
			}

			if len(result.Subcategories) != tt.expectedSubcategoriesCount {
				t.Errorf("Expected %d subcategories, but got %d", tt.expectedSubcategoriesCount, len(result.Subcategories))
			}

			if len(result.CategoryStats) != tt.expectedCategoryStatsCount {
				t.Errorf("Expected %d category stats, but got %d", tt.expectedCategoryStatsCount, len(result.CategoryStats))
			}

			if !tt.expectError && tt.expectedCategoryStatsCount > 0 {
				/* Verify category stats structure */
				for _, catStat := range result.CategoryStats {
					if catStat.Category == "" {
						t.Error("Category name should not be empty")
					}

					if catStat.TotalCount < 0 {
						t.Errorf("Category total count should not be negative, got %d", catStat.TotalCount)
					}
				}
			}
		})
	}
}

func TestFetchCategoryStats(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                     string
		mockCategoryDetails      []repository.GetCategoryDetailsRow
		mockCategoryDetailsErr   error
		mockSubcategoryCounts    []repository.CountCheatsheetsByCategoryAndSubcategoryRow
		mockSubcategoryCountsErr error
		expectedStatsCount       int
		expectedFirstCategory    string
		expectedFirstTotalCount  int64
		expectedFirstSubcatCount int
	}{
		{
			name: "successfully fetches category stats with subcategories",
			mockCategoryDetails: []repository.GetCategoryDetailsRow{
				{
					Category:        "html",
					CheatsheetCount: 50,
					Subcategories:   []string{"concepts"},
				},
				{
					Category:        "css",
					CheatsheetCount: 50,
					Subcategories:   []string{"attributes"},
				},
			},
			mockSubcategoryCounts: []repository.CountCheatsheetsByCategoryAndSubcategoryRow{
				{
					Category:        "html",
					Subcategory:     "concepts",
					CheatsheetCount: 50,
				},
				{
					Category:        "css",
					Subcategory:     "attributes",
					CheatsheetCount: 50,
				},
			},
			expectedStatsCount:       2,
			expectedFirstCategory:    "html",
			expectedFirstTotalCount:  50,
			expectedFirstSubcatCount: 1,
		},
		{
			name:                   "returns empty array when category details fails",
			mockCategoryDetailsErr: fmt.Errorf("db error"),
			expectedStatsCount:     0,
		},
		{
			name: "returns stats without subcategory breakdown when subcategory count fails",
			mockCategoryDetails: []repository.GetCategoryDetailsRow{
				{
					Category:        "html",
					CheatsheetCount: 50,
					Subcategories:   []string{"concepts"},
				},
				{
					Category:        "css",
					CheatsheetCount: 50,
					Subcategories:   []string{"attributes"},
				},
			},
			mockSubcategoryCountsErr: fmt.Errorf("db error"),
			expectedStatsCount:       0,
		},
		{
			name:                  "returns empty array for empty category details",
			mockCategoryDetails:   []repository.GetCategoryDetailsRow{},
			mockSubcategoryCounts: []repository.CountCheatsheetsByCategoryAndSubcategoryRow{},
			expectedStatsCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates a mock repo */
			mockRepo := &mocks.MockQuerier{
				GetCategoryDetailsFunc: func(ctx context.Context) ([]repository.GetCategoryDetailsRow, error) {
					if tt.mockCategoryDetailsErr != nil {
						return nil, tt.mockCategoryDetailsErr
					}

					return tt.mockCategoryDetails, nil
				},

				CountCheatsheetsByCategoryAndSubcategoryFunc: func(ctx context.Context) ([]repository.CountCheatsheetsByCategoryAndSubcategoryRow, error) {
					if tt.mockSubcategoryCountsErr != nil {
						return nil, tt.mockSubcategoryCountsErr
					}

					return tt.mockSubcategoryCounts, nil
				},
			}

			/* Creates the config service */
			service := &configService{
				repo: mockRepo,
				db:   nil,
			}

			result := service.fetchCategoryStats(ctx)

			if len(result) != tt.expectedStatsCount {
				t.Errorf("Expected %d stats, but got %d", tt.expectedStatsCount, len(result))
			}

			if tt.expectedStatsCount > 0 && len(result) > 0 {
				firstStat := result[0]

				if firstStat.Category != tt.expectedFirstCategory {
					t.Errorf("Expected first category %q, but got %q", tt.expectedFirstCategory, firstStat.Category)
				}

				if firstStat.TotalCount != tt.expectedFirstTotalCount {
					t.Errorf("Expected first total count %d, but got %d", tt.expectedFirstTotalCount, firstStat.TotalCount)
				}

				if len(firstStat.SubcategoriesStats) != tt.expectedFirstSubcatCount {
					t.Errorf("Expected %d subcategory stats, but got %d", tt.expectedFirstSubcatCount, len(firstStat.SubcategoriesStats))
				}
			}
		})
	}
}

func TestGetUsage(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                    string
		mockDatabaseSize        entities.DatabaseSize
		mockDatabaseSizeErr     error
		mockLargestTablesRows   [][]interface{}
		mockLargestTablesErr    error
		mockStorageSize         repository.GetTotalImageSizeRow
		mockStorageSizeErr      error
		mockLargestFiles        []repository.GetLargestCheatsheetsRow
		mockLargestFilesErr     error
		expectError             bool
		expectErrorMsg          string
		expectedDatabaseBytes   int64
		expectedStorageBytes    int64
		expectedDatabasePercent float64
		expectedStoragePercent  float64
		expectedTablesCount     int
		expectedFilesCount      int
	}{
		{
			name: "successfully fetches complete usage",
			mockDatabaseSize: entities.DatabaseSize{
				SizeBytes:  314572800, // 300MB
				SizePretty: "300 MB",
			},
			mockLargestTablesRows: [][]interface{}{
				{"public", "cheatsheets", "100 MB"},
				{"public", "events", "50 MB"},
				{"public", "pageviews", "30 MB"},
			},
			mockStorageSize: repository.GetTotalImageSizeRow{
				TotalSize:       249561088, // 250MB
				TotalSizePretty: "250 MB",
			},
			mockLargestFiles: []repository.GetLargestCheatsheetsRow{
				{Title: "HTML Boilerplate", Category: "html", Size: "120 KB"},
				{Title: "CSS Rulesets", Category: "css", Size: "150 KB"},
			},
			expectError:             false,
			expectedDatabaseBytes:   314572800,
			expectedStorageBytes:    249561088,
			expectedDatabasePercent: 60.0, // 500MB is limit for Database
			expectedStoragePercent:  23.0, // 1GB is limit for Storage
			expectedTablesCount:     3,
			expectedFilesCount:      2,
		},
		{
			name:                "handles error of getting database size",
			mockDatabaseSizeErr: fmt.Errorf("database query failed"),
			expectError:         true,
			expectErrorMsg:      "failed to get database usage info",
		},
		{
			name: "handles error of getting largest tables",
			mockDatabaseSize: entities.DatabaseSize{
				SizeBytes:  300000000, // 300MB
				SizePretty: "300 MB",
			},
			mockLargestTablesErr: fmt.Errorf("query failed"),
			expectError:          true,
			expectErrorMsg:       "failed to get largest tables usage info",
		},
		{
			name: "handles error of getting storage size",
			mockDatabaseSize: entities.DatabaseSize{
				SizeBytes:  300000000, // 300MB
				SizePretty: "300 MB",
			},
			mockLargestTablesRows: [][]interface{}{
				{"public", "cheatsheets", "100 MB"},
				{"public", "events", "50 MB"},
				{"public", "pageviews", "30 MB"},
			},
			mockStorageSizeErr: fmt.Errorf("query failed"),
			expectError:        true,
			expectErrorMsg:     "failed to get storage usage info",
		},
		{
			name: "handles error of getting largest files",
			mockDatabaseSize: entities.DatabaseSize{
				SizeBytes:  300000000, // 300MB
				SizePretty: "300 MB",
			},
			mockLargestTablesRows: [][]interface{}{
				{"public", "cheatsheets", "100 MB"},
				{"public", "events", "50 MB"},
				{"public", "pageviews", "30 MB"},
			},
			mockStorageSize: repository.GetTotalImageSizeRow{
				TotalSize:       250000000, // 250MB
				TotalSizePretty: "250 MB",
			},
			mockLargestFilesErr: fmt.Errorf("files query failed"),
			expectError:         true,
			expectErrorMsg:      "failed to get largest tables usage info",
		},
		{
			name: "caps usage percent at 100",
			mockDatabaseSize: entities.DatabaseSize{
				SizeBytes:  2000000000, // 2GB (over limit)
				SizePretty: "2 GB",
			},
			mockLargestTablesRows: [][]interface{}{
				{"public", "cheatsheets", "100 MB"},
			},
			mockStorageSize: repository.GetTotalImageSizeRow{
				TotalSize:       1500000000, // 1.5GB (over limit)
				TotalSizePretty: "1.5 GB",
			},
			mockLargestFiles: []repository.GetLargestCheatsheetsRow{
				{Title: "HTML Boilerplate", Category: "html", Size: "120 KB"},
			},
			expectError:             false,
			expectedTablesCount:     1,
			expectedFilesCount:      1,
			expectedDatabasePercent: 100.0,
			expectedStoragePercent:  100.0,
		},
		{
			name: "handles empty largest tables and files",
			mockDatabaseSize: entities.DatabaseSize{
				SizeBytes:  100000000,
				SizePretty: "100 MB",
			},
			mockLargestTablesRows: [][]interface{}{},
			mockStorageSize: repository.GetTotalImageSizeRow{
				TotalSize:       50000000,
				TotalSizePretty: "50 MB",
			},
			mockLargestFiles:    []repository.GetLargestCheatsheetsRow{},
			expectError:         false,
			expectedTablesCount: 0,
			expectedFilesCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates a mock database pool */
			mockDB := &mocks.MockDatabasePool{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					/* Mock getDatabaseSize query */

					return &mocks.MockRow{
						ScanFunc: func(dest ...interface{}) error {

							if tt.mockDatabaseSizeErr != nil {
								return tt.mockDatabaseSizeErr
							}

							if tt.mockDatabaseSize.SizeBytes >= 0 {
								/* Scans into dest pointers */
								if sizeBytes, ok := dest[0].(*int64); ok {
									*sizeBytes = tt.mockDatabaseSize.SizeBytes
								}
								if sizePretty, ok := dest[1].(*string); ok {
									*sizePretty = tt.mockDatabaseSize.SizePretty
								}
							}
							return nil
						},
					}
				},

				QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
					/* Mock getLargestTables query */
					if tt.mockLargestTablesErr != nil {
						return nil, tt.mockLargestTablesErr
					}
					return mocks.NewMockRows(tt.mockLargestTablesRows), nil
				},
			}

			/* Creates a mock repo */
			mockRepo := &mocks.MockQuerier{
				GetTotalImageSizeFunc: func(ctx context.Context) (repository.GetTotalImageSizeRow, error) {
					if tt.mockStorageSizeErr != nil {
						return repository.GetTotalImageSizeRow{}, tt.mockStorageSizeErr
					}

					return tt.mockStorageSize, nil
				},

				GetLargestCheatsheetsFunc: func(ctx context.Context) ([]repository.GetLargestCheatsheetsRow, error) {
					if tt.mockLargestFilesErr != nil {
						return nil, tt.mockLargestFilesErr
					}

					return tt.mockLargestFiles, nil
				},
			}

			/* Creates the config service */
			service := NewConfigService(mockRepo, mockDB)

			result, err := service.GetUsage(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tt.expectErrorMsg {
					t.Errorf("Expected error messaged = %q, but got %q", tt.expectErrorMsg, err.Error())
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			/* Verify database usage */
			if tt.expectedDatabaseBytes > 0 && result.Database.SizeBytes != tt.expectedDatabaseBytes {
				t.Errorf("Expected database size %d, but got %d", tt.expectedDatabaseBytes, result.Database.SizeBytes)
			}

			if tt.expectedDatabasePercent > 0 && result.Database.UsagePercent != tt.expectedDatabasePercent {
				t.Errorf("Expected database usage percent %.2f, but got %.2f", tt.expectedDatabasePercent, result.Database.UsagePercent)
			}

			/* Verify storage usage */
			if tt.expectedStorageBytes > 0 && result.Storage.SizeBytes != tt.expectedStorageBytes {
				t.Errorf("Expected storage size %d, but got %d", tt.expectedStorageBytes, result.Storage.SizeBytes)
			}

			if tt.expectedStoragePercent > 0 && result.Storage.UsagePercent != tt.expectedStoragePercent {
				t.Errorf("Expected storage usage percent %.2f, but got %.2f", tt.expectedStoragePercent, result.Storage.UsagePercent)
			}

			/* Verify counts */
			if len(result.Database.LargestTables) != tt.expectedTablesCount {
				t.Errorf("Expected %d largest tables, got %d", tt.expectedTablesCount, len(result.Database.LargestTables))
			}

			if len(result.Storage.Files) != tt.expectedFilesCount {
				t.Errorf("Expected %d largest files, got %d", tt.expectedFilesCount, len(result.Storage.Files))
			}

			/* Verify timestamp is set */
			if result.Timestamp == "" {
				t.Error("Expected timestamp to be set")
			}
		})
	}
}
