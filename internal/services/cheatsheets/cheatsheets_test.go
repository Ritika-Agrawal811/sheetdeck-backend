package cheatsheets

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/mocks"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateCheatsheet(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name              string
		slug              string
		title             string
		category          string
		subcategory       string
		imageSize         int64
		mockCheatsheetErr error
		mockUploadErr     error
		expectError       bool
		expectedErrorMsg  string
	}{
		{
			name:        "successfully creates cheat sheet",
			slug:        "css-target-pseudo-class",
			title:       ":target Pseudo Class",
			category:    "css",
			subcategory: "pseudo_classes",
			imageSize:   1024000,
			expectError: false,
		},
		{
			name:             "handles error of empty slug",
			slug:             "",
			title:            ":target Pseudo Class",
			category:         "css",
			subcategory:      "pseudo_classes",
			imageSize:        1024000,
			expectError:      true,
			expectedErrorMsg: "missing required fields",
		},
		{
			name:             "handles error of empty title",
			slug:             "css-target-pseudo-class",
			title:            "",
			category:         "css",
			subcategory:      "pseudo_classes",
			imageSize:        1024000,
			expectError:      true,
			expectedErrorMsg: "missing required fields",
		},
		{
			name:             "handles error of empty category",
			slug:             "css-target-pseudo-class",
			title:            ":target Pseudo Class",
			category:         "",
			subcategory:      "pseudo_classes",
			imageSize:        1024000,
			expectError:      true,
			expectedErrorMsg: "missing required fields",
		},
		{
			name:             "handles error of empty category",
			slug:             "css-target-pseudo-class",
			title:            ":target Pseudo Class",
			category:         "css",
			subcategory:      "",
			imageSize:        1024000,
			expectError:      true,
			expectedErrorMsg: "missing required fields",
		},
		{
			name:             "handles upload error",
			slug:             "css-target-pseudo-class",
			title:            ":target Pseudo Class",
			category:         "css",
			subcategory:      "pseudo_classes",
			imageSize:        1024000,
			mockUploadErr:    fmt.Errorf("network error"),
			expectError:      true,
			expectedErrorMsg: "failed to upload the cheat sheet in storage",
		},
		{
			name:              "handles database error for creating cheatsheet",
			slug:              "css-target-pseudo-class",
			title:             ":target Pseudo Class",
			category:          "css",
			subcategory:       "pseudo_classes",
			imageSize:         1024000,
			mockCheatsheetErr: fmt.Errorf("database error"),
			expectError:       true,
			expectedErrorMsg:  "failed to create cheatsheet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates mock Repo */
			mockRepo := &mocks.MockQuerier{
				CreateCheatsheetFunc: func(ctx context.Context, arg repository.CreateCheatsheetParams) error {
					if arg.Slug == "" || arg.Title == "" || arg.Category == "" || arg.Subcategory == "" {
						return fmt.Errorf("missing required fields")
					}

					if tt.mockCheatsheetErr != nil {
						return tt.mockCheatsheetErr
					}

					return nil
				},
			}

			mockStorageService := &mocks.MockStorageService{
				UploadFileFunc: func(fileName string, image multipart.File) (string, error) {
					if fileName == "" {
						return "", fmt.Errorf("filename cannot be empty")
					}

					if tt.mockUploadErr != nil {
						return "", tt.mockUploadErr
					}

					return "https://example.supabase.co/" + fileName, nil
				},
			}

			/* Creates service with mockRepo and mockStorageService */
			service := NewCheatsheetsService(mockRepo, mockStorageService)

			details := dtos.Cheatsheet{
				Slug:        tt.slug,
				Title:       tt.title,
				Category:    tt.category,
				SubCategory: tt.subcategory,
			}

			mockFile := mocks.CreateMockFile("test content")

			err := service.CreateCheatsheet(ctx, details, mockFile, tt.imageSize)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				if tt.expectedErrorMsg != "" {
					if !strings.Contains(err.Error(), tt.expectedErrorMsg) {
						t.Errorf("Expected error to contain %q, got %q", tt.expectedErrorMsg, err.Error())
					}
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}

func TestGetCheatsheetByID(t *testing.T) {
	ctx := context.Background()

	/* Valid UUID for testing */
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	/* Convert to pgtype.UUID for mock response */
	var pgUUID pgtype.UUID
	copy(pgUUID.Bytes[:], []byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00})
	pgUUID.Valid = true

	/* Test timestamps */
	now := time.Now().UTC()
	createdAt := pgtype.Timestamptz{Time: now.Add(-24 * time.Hour), Valid: true} // 1 day ago
	updatedAt := pgtype.Timestamptz{Time: now, Valid: true}

	tests := []struct {
		name              string
		id                string
		slug              string
		title             string
		category          string
		subcategory       string
		imageUrl          string
		imageSize         int64
		createdAt         pgtype.Timestamptz
		updatedAt         pgtype.Timestamptz
		mockCheatsheetErr error
		expectError       bool
		expectedErrorMsg  string
		validateResult    func(*testing.T, *repository.GetCheatsheetByIDRow)
	}{
		{
			name:        "successfully fetches cheat sheet details",
			id:          validUUID,
			slug:        "css-target-pseudo-class",
			title:       ":target Pseudo Class",
			category:    "css",
			subcategory: "pseudo_classes",
			imageUrl:    "https://example.supabase.co/storage/v1/object/public/cheatsheets/css-target-pseudo-class.webp",
			imageSize:   1024000,
			createdAt:   createdAt,
			updatedAt:   updatedAt,
			expectError: false,
			validateResult: func(t *testing.T, result *repository.GetCheatsheetByIDRow) {
				if result.Title != ":target Pseudo Class" {
					t.Errorf("Expected title ':target Pseudo Class', got %q", result.Title)
				}

				if result.Category != "css" {
					t.Errorf("Expected category 'css', got %q", result.Category)
				}

				if !result.ImageUrl.Valid {
					t.Error("Expected ImageUrl to be valid")
				}

				if !result.CreatedAt.Valid {
					t.Error("Expected CreatedAt to be valid")
				}

				if !result.UpdatedAt.Valid {
					t.Error("Expected UpdatedAt to be valid")
				}

				if !result.ImageSize.Valid {
					t.Error("Expected ImageSize to be valid")
				}
			},
		},
		{
			name:        "successfully handles cheatsheet with null optional fields",
			id:          validUUID,
			slug:        "css-target-pseudo-class",
			title:       ":target Pseudo Class",
			category:    "css",
			subcategory: "pseudo_classes",
			imageUrl:    "",
			imageSize:   0,
			createdAt:   createdAt,
			updatedAt:   updatedAt,
			expectError: false,
			validateResult: func(t *testing.T, result *repository.GetCheatsheetByIDRow) {
				if result.ImageUrl.Valid {
					t.Error("Expected ImageUrl to be invalid")
				}

				if result.ImageSize.Int64 != 0 {
					t.Errorf("Expected ImageSize to be 0, but got %d", result.ImageSize.Int64)
				}
			},
		},
		{
			name:             "returns error for invalid UUID format",
			id:               "invalid-uuid",
			expectError:      true,
			expectedErrorMsg: "invalid UUID format",
		},
		{
			name:             "returns error for empty string id",
			id:               "",
			expectError:      true,
			expectedErrorMsg: "invalid UUID format",
		},
		{
			name:             "returns error for malformed UUID",
			id:               "550e8400-e29b-41d4",
			expectError:      true,
			expectedErrorMsg: "invalid UUID format",
		},
		{
			name:              "handles error of cheatsheet not found",
			id:                validUUID,
			mockCheatsheetErr: fmt.Errorf("no rows in result set"),
			expectError:       true,
			expectedErrorMsg:  "failed to fetch cheatsheet",
		},
		{
			name:              "handles error of database connection error",
			id:                validUUID,
			mockCheatsheetErr: fmt.Errorf("database connection failed"),
			expectError:       true,
			expectedErrorMsg:  "failed to fetch cheatsheet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates a mock repo */
			mockRepo := &mocks.MockQuerier{
				GetCheatsheetByIDFunc: func(ctx context.Context, id pgtype.UUID) (repository.GetCheatsheetByIDRow, error) {
					if !id.Valid {
						t.Error("Expected valid UUID to be passed to repository")
					}

					if tt.mockCheatsheetErr != nil {
						return repository.GetCheatsheetByIDRow{}, tt.mockCheatsheetErr
					}

					return repository.GetCheatsheetByIDRow{
						ID:          pgUUID,
						Slug:        tt.slug,
						Title:       tt.title,
						Category:    repository.Category(tt.category),
						Subcategory: repository.Subcategory(tt.subcategory),
						ImageUrl:    utils.PgText(tt.imageUrl),
						CreatedAt:   tt.createdAt,
						UpdatedAt:   tt.updatedAt,
						ImageSize:   utils.PgInt8(tt.imageSize),
					}, nil
				},
			}

			/* Creates service with mockRepo */
			service := NewCheatsheetsService(mockRepo, nil)

			result, err := service.GetCheatsheetByID(ctx, tt.id)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				if tt.expectedErrorMsg != "" {
					if !strings.Contains(err.Error(), tt.expectedErrorMsg) {
						t.Errorf("Expected error to contain %q, got %q", tt.expectedErrorMsg, err.Error())
					}
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result to not be nil")
				return
			}

			/* Run custom validation if provided */
			if tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestGetCheatsheetBySlug(t *testing.T) {
	ctx := context.Background()

	/* Valid UUID for testing */
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	/* Convert to pgtype.UUID for mock response */
	var pgUUID pgtype.UUID
	copy(pgUUID.Bytes[:], []byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00})
	pgUUID.Valid = true

	/* Test timestamps */
	now := time.Now().UTC()
	createdAt := pgtype.Timestamptz{Time: now.Add(-24 * time.Hour), Valid: true} // 1 day ago
	updatedAt := pgtype.Timestamptz{Time: now, Valid: true}

	tests := []struct {
		name              string
		id                string
		slug              string
		title             string
		category          string
		subcategory       string
		imageUrl          string
		imageSize         int64
		createdAt         pgtype.Timestamptz
		updatedAt         pgtype.Timestamptz
		mockCheatsheetErr error
		expectError       bool
		expectedErrorMsg  string
		validateResult    func(*testing.T, *repository.GetCheatsheetBySlugRow)
	}{
		{
			name:        "successfully fetches cheat sheet details",
			id:          validUUID,
			slug:        "css-target-pseudo-class",
			title:       ":target Pseudo Class",
			category:    "css",
			subcategory: "pseudo_classes",
			imageUrl:    "https://example.supabase.co/storage/v1/object/public/cheatsheets/css-target-pseudo-class.webp",
			imageSize:   1024000,
			createdAt:   createdAt,
			updatedAt:   updatedAt,
			expectError: false,
			validateResult: func(t *testing.T, result *repository.GetCheatsheetBySlugRow) {
				if result.Title != ":target Pseudo Class" {
					t.Errorf("Expected title ':target Pseudo Class', got %q", result.Title)
				}

				if result.Category != "css" {
					t.Errorf("Expected category 'css', got %q", result.Category)
				}

				if !result.ImageUrl.Valid {
					t.Error("Expected ImageUrl to be valid")
				}

				if !result.CreatedAt.Valid {
					t.Error("Expected CreatedAt to be valid")
				}

				if !result.UpdatedAt.Valid {
					t.Error("Expected UpdatedAt to be valid")
				}

				if !result.ImageSize.Valid {
					t.Error("Expected ImageSize to be valid")
				}
			},
		},
		{
			name:        "successfully handles cheatsheet with null optional fields",
			id:          validUUID,
			slug:        "css-target-pseudo-class",
			title:       ":target Pseudo Class",
			category:    "css",
			subcategory: "pseudo_classes",
			imageUrl:    "",
			imageSize:   0,
			createdAt:   createdAt,
			updatedAt:   updatedAt,
			expectError: false,
			validateResult: func(t *testing.T, result *repository.GetCheatsheetBySlugRow) {
				if result.ImageUrl.Valid {
					t.Error("Expected ImageUrl to be invalid")
				}

				if result.ImageSize.Int64 != 0 {
					t.Errorf("Expected ImageSize to be 0, but got %d", result.ImageSize.Int64)
				}
			},
		},
		{
			name:              "handles error of empty string slug",
			slug:              "",
			mockCheatsheetErr: fmt.Errorf("no rows in result set"),
			expectError:       true,
			expectedErrorMsg:  "failed to fetch cheatsheet",
		},
		{
			name:              "handles error of cheatsheet not found",
			slug:              "css-target-pseudo-class",
			mockCheatsheetErr: fmt.Errorf("no rows in result set"),
			expectError:       true,
			expectedErrorMsg:  "failed to fetch cheatsheet",
		},
		{
			name:              "handles error of database connection error",
			slug:              "css-target-pseudo-class",
			mockCheatsheetErr: fmt.Errorf("database connection failed"),
			expectError:       true,
			expectedErrorMsg:  "failed to fetch cheatsheet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates a mock repo */
			mockRepo := &mocks.MockQuerier{
				GetCheatsheetBySlugFunc: func(ctx context.Context, slug string) (repository.GetCheatsheetBySlugRow, error) {
					if slug == "" {
						return repository.GetCheatsheetBySlugRow{}, fmt.Errorf("Expected non empty slug to be passed to repository")
					}

					if tt.mockCheatsheetErr != nil {
						return repository.GetCheatsheetBySlugRow{}, tt.mockCheatsheetErr
					}

					return repository.GetCheatsheetBySlugRow{
						ID:          pgUUID,
						Slug:        tt.slug,
						Title:       tt.title,
						Category:    repository.Category(tt.category),
						Subcategory: repository.Subcategory(tt.subcategory),
						ImageUrl:    utils.PgText(tt.imageUrl),
						CreatedAt:   tt.createdAt,
						UpdatedAt:   tt.updatedAt,
						ImageSize:   utils.PgInt8(tt.imageSize),
					}, nil
				},
			}

			/* Creates service with mockRepo */
			service := NewCheatsheetsService(mockRepo, nil)

			result, err := service.GetCheatsheetBySlug(ctx, tt.slug)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				if tt.expectedErrorMsg != "" {
					if !strings.Contains(err.Error(), tt.expectedErrorMsg) {
						t.Errorf("Expected error to contain %q, got %q", tt.expectedErrorMsg, err.Error())
					}
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result to not be nil")
				return
			}

			/* Run custom validation if provided */
			if tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}

}

func TestGetAllCheatsheets(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                     string
		category                 string
		subcategory              string
		sortBy                   string
		limit                    string
		mockCheatsheetErr        error
		expectError              bool
		expectedErrorMsg         string
		expectedCheatsheetsCount int64
		validateResult           func(*testing.T, []repository.ListCheatsheetsRow)
	}{
		{
			name:                     "successfully fetches cheatsheets",
			expectError:              false,
			expectedCheatsheetsCount: 15,
		},
		{
			name:                     "handles limit correctly",
			limit:                    "5",
			expectError:              false,
			expectedCheatsheetsCount: 5,
		},
		{
			name:                     "handles limit set as 0",
			limit:                    "0",
			expectError:              false,
			expectedCheatsheetsCount: 15,
		},
		{
			name:                     "handles sort by category",
			category:                 "html",
			limit:                    "12",
			expectError:              false,
			expectedCheatsheetsCount: 12,
			validateResult: func(t *testing.T, cheatsheets []repository.ListCheatsheetsRow) {
				for _, item := range cheatsheets {
					if item.Category != "html" {
						t.Errorf("Expected category to be 'html', but got %q", item.Category)
					}
				}
			},
		},
		{
			name:                     "handles sort by subcategory",
			subcategory:              "concepts",
			limit:                    "5",
			expectError:              false,
			expectedCheatsheetsCount: 5,
			validateResult: func(t *testing.T, cheatsheets []repository.ListCheatsheetsRow) {
				for _, item := range cheatsheets {
					if item.Subcategory != "concepts" {
						t.Errorf("Expected subcategory to be 'concepts', but got %q", item.Subcategory)
					}
				}
			},
		},
		{
			name:                     "handles sort by category and subcategory",
			category:                 "css",
			subcategory:              "pseudo_classes",
			limit:                    "3",
			expectError:              false,
			expectedCheatsheetsCount: 3,
			validateResult: func(t *testing.T, cheatsheets []repository.ListCheatsheetsRow) {
				for _, item := range cheatsheets {
					if item.Category != "css" {
						t.Errorf("Expected category to be 'css', but got %q", item.Category)
					}

					if item.Subcategory != "pseudo_classes" {
						t.Errorf("Expected subcategory to be 'pseudo_classes', but got %q", item.Subcategory)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates a mock repo */
			mockRepo := &mocks.MockQuerier{
				ListCheatsheetsFunc: func(ctx context.Context, arg repository.ListCheatsheetsParams) ([]repository.ListCheatsheetsRow, error) {
					if tt.mockCheatsheetErr != nil {
						return []repository.ListCheatsheetsRow{}, tt.mockCheatsheetErr
					}

					return getCheatsheets(arg), nil
				},
			}

			/* Creates service with mockRepo */
			service := NewCheatsheetsService(mockRepo, nil)

			result, err := service.GetAllCheatsheets(ctx, tt.category, tt.subcategory, tt.sortBy, tt.limit)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				if tt.expectedErrorMsg != "" {
					if !strings.Contains(err.Error(), tt.expectedErrorMsg) {
						t.Errorf("Expected error to contain %q, got %q", tt.expectedErrorMsg, err.Error())
					}
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result to not be nil")
				return
			}

			if len(result) != int(tt.expectedCheatsheetsCount) {
				t.Errorf("Expected %d cheatsheets, got %d", tt.expectedCheatsheetsCount, len(result))
			}

			if tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

/* Mock for ListCheatsheets */
func getCheatsheets(params repository.ListCheatsheetsParams) []repository.ListCheatsheetsRow {
	cheatsheets := []repository.ListCheatsheetsRow{
		{
			Slug:        "css-target-pseudo-class",
			Title:       ":target Pseudo Class",
			Category:    "css",
			Subcategory: "pseudo_classes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
			Downloads:   3,
			Views:       1,
		},
		{
			Slug:        "css-first-child-pseudo-class",
			Title:       ":first-child Pseudo Class",
			Category:    "css",
			Subcategory: "pseudo_classes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(1 * time.Minute), Valid: true},
			Downloads:   3,
			Views:       1,
		},
		{
			Slug:        "css-active-pseudo-class",
			Title:       ":active Pseudo Class",
			Category:    "css",
			Subcategory: "pseudo_classes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(2 * time.Minute), Valid: true},
			Downloads:   4,
			Views:       1,
		},
		{
			Slug:        "html-semantic-elements",
			Title:       "Semantic Elements",
			Category:    "html",
			Subcategory: "concepts",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(3 * time.Minute), Valid: true},
			Downloads:   3,
			Views:       0,
		},
		{
			Slug:        "html-inline-elements",
			Title:       "Inline Elements",
			Category:    "html",
			Subcategory: "concepts",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(4 * time.Minute), Valid: true},
			Downloads:   3,
			Views:       0,
		},
		{
			Slug:        "html-hide-decorative-icons",
			Title:       "Hide Decorative Icons",
			Category:    "html",
			Subcategory: "concepts",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(5 * time.Minute), Valid: true},
			Downloads:   3,
			Views:       2,
		},
		{
			Slug:        "html-headings-in-order",
			Title:       "Headings in Hierarchical Order",
			Category:    "html",
			Subcategory: "concepts",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(6 * time.Minute), Valid: true},
			Downloads:   3,
			Views:       1,
		},
		{
			Slug:        "html-clear-link-texts",
			Title:       "Meaningful Link Texts",
			Category:    "html",
			Subcategory: "concepts",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(7 * time.Minute), Valid: true},
			Downloads:   3,
			Views:       1,
		},
		{
			Slug:        "html-target-attribute",
			Title:       "target Attribute",
			Category:    "html",
			Subcategory: "attributes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(8 * time.Minute), Valid: true},
			Downloads:   1,
			Views:       0,
		},
		{
			Slug:        "html-download-attribute",
			Title:       "download Attribute",
			Category:    "html",
			Subcategory: "attributes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(9 * time.Minute), Valid: true},
			Downloads:   1,
			Views:       0,
		},
		{
			Slug:        "html-lang-attribute",
			Title:       "lang Attribute",
			Category:    "html",
			Subcategory: "attributes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(10 * time.Minute), Valid: true},
			Downloads:   1,
			Views:       0,
		},
		{
			Slug:        "html-inputmode-attribute",
			Title:       "inputmode Attribute",
			Category:    "html",
			Subcategory: "attributes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(11 * time.Minute), Valid: true},
			Downloads:   2,
			Views:       0,
		},
		{
			Slug:        "html-tabindex-attribute",
			Title:       "tabindex Attribute",
			Category:    "html",
			Subcategory: "attributes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(12 * time.Minute), Valid: true},
			Downloads:   1,
			Views:       0,
		},
		{
			Slug:        "html-disabled-attribute",
			Title:       "disabled Attribute",
			Category:    "html",
			Subcategory: "attributes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(13 * time.Minute), Valid: true},
			Downloads:   1,
			Views:       0,
		},
		{
			Slug:        "html-accept-attribute",
			Title:       "accept Attribute",
			Category:    "html",
			Subcategory: "attributes",
			CreatedAt:   pgtype.Timestamptz{Time: time.Now().UTC().Add(14 * time.Minute), Valid: true},
			Downloads:   1,
			Views:       0,
		},
	}

	var results []repository.ListCheatsheetsRow
	for _, item := range cheatsheets {
		// Filter by category if provided
		if params.Category.Valid && item.Category != string(params.Category.Category) {
			continue
		}

		// Filter by subcategory if provided
		if params.Subcategory.Valid && item.Subcategory != string(params.Subcategory.Subcategory) {
			continue
		}

		results = append(results, item)
	}

	return results[:params.Limit]
}
