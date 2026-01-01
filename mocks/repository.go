package mocks

import (
	"context"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

/** Mock Implementation for Database Repository **/

/**
 * This is a compile-time interface compliance check in Go.
 * Ensures MockQuerier implements repository.Querier interface
 */
var _ repository.Querier = (*MockQuerier)(nil)

type MockQuerier struct {
	CountCheatsheetsByCategoryAndSubcategoryFunc func(ctx context.Context) ([]repository.CountCheatsheetsByCategoryAndSubcategoryRow, error)
	CreateCheatsheetFunc                         func(ctx context.Context, arg repository.CreateCheatsheetParams) error
	GetBrowsersSummaryByDayFunc                  func(ctx context.Context, days int32) ([]repository.GetBrowsersSummaryByDayRow, error)
	GetBrowsersSummaryForLast24HoursFunc         func(ctx context.Context) ([]repository.GetBrowsersSummaryForLast24HoursRow, error)
	GetCategoriesFunc                            func(ctx context.Context) ([]string, error)
	GetCategoryDetailsFunc                       func(ctx context.Context) ([]repository.GetCategoryDetailsRow, error)
	GetCheatsheetByIDFunc                        func(ctx context.Context, id pgtype.UUID) (repository.GetCheatsheetByIDRow, error)
	GetCheatsheetBySlugFunc                      func(ctx context.Context, slug string) (repository.GetCheatsheetBySlugRow, error)
	GetCountriesSummaryByDayFunc                 func(ctx context.Context, days int32) ([]repository.GetCountriesSummaryByDayRow, error)
	GetCountriesSummaryForLast24HoursFunc        func(ctx context.Context) ([]repository.GetCountriesSummaryForLast24HoursRow, error)
	GetDevicesSummaryByDayFunc                   func(ctx context.Context, days int32) ([]repository.GetDevicesSummaryByDayRow, error)
	GetDevicesSummaryForLast24HoursFunc          func(ctx context.Context) ([]repository.GetDevicesSummaryForLast24HoursRow, error)
	GetLargestCheatsheetsFunc                    func(ctx context.Context) ([]repository.GetLargestCheatsheetsRow, error)
	GetMetricsTimeseriesByDayFunc                func(ctx context.Context, days int32) ([]repository.GetMetricsTimeseriesByDayRow, error)
	GetMetricsTimeseriesForLast24HoursFunc       func(ctx context.Context) ([]repository.GetMetricsTimeseriesForLast24HoursRow, error)
	GetOSSummaryByDayFunc                        func(ctx context.Context, days int32) ([]repository.GetOSSummaryByDayRow, error)
	GetOSSummaryForLast24HoursFunc               func(ctx context.Context) ([]repository.GetOSSummaryForLast24HoursRow, error)
	GetReferrerSummaryByDayFunc                  func(ctx context.Context, days int32) ([]repository.GetReferrerSummaryByDayRow, error)
	GetReferrerSummaryForLast24HoursFunc         func(ctx context.Context) ([]repository.GetReferrerSummaryForLast24HoursRow, error)
	GetRoutesSummaryByDayFunc                    func(ctx context.Context, days int32) ([]repository.GetRoutesSummaryByDayRow, error)
	GetRoutesSummaryForLast24HoursFunc           func(ctx context.Context) ([]repository.GetRoutesSummaryForLast24HoursRow, error)
	GetSubcategoriesFunc                         func(ctx context.Context) ([]string, error)
	GetTotalCheasheetsCountFunc                  func(ctx context.Context) (int64, error)
	GetTotalClicksAndDownloadsFunc               func(ctx context.Context) (repository.GetTotalClicksAndDownloadsRow, error)
	GetTotalImageSizeFunc                        func(ctx context.Context) (repository.GetTotalImageSizeRow, error)
	GetTotalViewsAndVisitorsFunc                 func(ctx context.Context) (repository.GetTotalViewsAndVisitorsRow, error)
	ListCheatsheetsFunc                          func(ctx context.Context, arg repository.ListCheatsheetsParams) ([]repository.ListCheatsheetsRow, error)
	StoreEventFunc                               func(ctx context.Context, arg repository.StoreEventParams) error
	StorePageviewFunc                            func(ctx context.Context, arg repository.StorePageviewParams) error
	UpdateCheatsheetFunc                         func(ctx context.Context, arg repository.UpdateCheatsheetParams) error
}

/* Implement all interface methods */

func (m *MockQuerier) CountCheatsheetsByCategoryAndSubcategory(ctx context.Context) ([]repository.CountCheatsheetsByCategoryAndSubcategoryRow, error) {
	if m.CountCheatsheetsByCategoryAndSubcategoryFunc != nil {
		return m.CountCheatsheetsByCategoryAndSubcategoryFunc(ctx)
	}

	return []repository.CountCheatsheetsByCategoryAndSubcategoryRow{}, nil
}

func (m *MockQuerier) CreateCheatsheet(ctx context.Context, arg repository.CreateCheatsheetParams) error {
	if m.CreateCheatsheetFunc != nil {
		return m.CreateCheatsheetFunc(ctx, arg)
	}

	return nil
}

func (m *MockQuerier) GetBrowsersSummaryByDay(ctx context.Context, days int32) ([]repository.GetBrowsersSummaryByDayRow, error) {
	if m.GetBrowsersSummaryByDayFunc != nil {
		return m.GetBrowsersSummaryByDayFunc(ctx, days)
	}

	return []repository.GetBrowsersSummaryByDayRow{}, nil
}

func (m *MockQuerier) GetBrowsersSummaryForLast24Hours(ctx context.Context) ([]repository.GetBrowsersSummaryForLast24HoursRow, error) {
	if m.GetBrowsersSummaryForLast24HoursFunc != nil {
		return m.GetBrowsersSummaryForLast24HoursFunc(ctx)
	}

	return []repository.GetBrowsersSummaryForLast24HoursRow{}, nil
}

func (m *MockQuerier) GetCategories(ctx context.Context) ([]string, error) {
	if m.GetCategoriesFunc != nil {
		return m.GetCategoriesFunc(ctx)
	}

	return []string{}, nil
}

func (m *MockQuerier) GetCategoryDetails(ctx context.Context) ([]repository.GetCategoryDetailsRow, error) {
	if m.GetCategoryDetailsFunc != nil {
		return m.GetCategoryDetailsFunc(ctx)
	}

	return []repository.GetCategoryDetailsRow{}, nil
}

func (m *MockQuerier) GetCheatsheetByID(ctx context.Context, id pgtype.UUID) (repository.GetCheatsheetByIDRow, error) {
	if m.GetCheatsheetByIDFunc != nil {
		return m.GetCheatsheetByIDFunc(ctx, id)
	}

	return repository.GetCheatsheetByIDRow{}, nil
}

func (m *MockQuerier) GetCheatsheetBySlug(ctx context.Context, slug string) (repository.GetCheatsheetBySlugRow, error) {
	if m.GetCheatsheetBySlugFunc != nil {
		return m.GetCheatsheetBySlugFunc(ctx, slug)
	}

	return repository.GetCheatsheetBySlugRow{}, nil
}

func (m *MockQuerier) GetCountriesSummaryByDay(ctx context.Context, days int32) ([]repository.GetCountriesSummaryByDayRow, error) {
	if m.GetCountriesSummaryByDayFunc != nil {
		return m.GetCountriesSummaryByDayFunc(ctx, days)
	}
	return []repository.GetCountriesSummaryByDayRow{}, nil
}

func (m *MockQuerier) GetCountriesSummaryForLast24Hours(ctx context.Context) ([]repository.GetCountriesSummaryForLast24HoursRow, error) {
	if m.GetCountriesSummaryForLast24HoursFunc != nil {
		return m.GetCountriesSummaryForLast24HoursFunc(ctx)
	}

	return []repository.GetCountriesSummaryForLast24HoursRow{}, nil
}

func (m *MockQuerier) GetDevicesSummaryByDay(ctx context.Context, days int32) ([]repository.GetDevicesSummaryByDayRow, error) {
	if m.GetDevicesSummaryByDayFunc != nil {
		return m.GetDevicesSummaryByDayFunc(ctx, days)
	}

	return []repository.GetDevicesSummaryByDayRow{}, nil
}

func (m *MockQuerier) GetDevicesSummaryForLast24Hours(ctx context.Context) ([]repository.GetDevicesSummaryForLast24HoursRow, error) {
	if m.GetDevicesSummaryForLast24HoursFunc != nil {
		return m.GetDevicesSummaryForLast24HoursFunc(ctx)
	}

	return []repository.GetDevicesSummaryForLast24HoursRow{}, nil
}

func (m *MockQuerier) GetLargestCheatsheets(ctx context.Context) ([]repository.GetLargestCheatsheetsRow, error) {
	if m.GetLargestCheatsheetsFunc != nil {
		return m.GetLargestCheatsheetsFunc(ctx)
	}

	return []repository.GetLargestCheatsheetsRow{}, nil
}

func (m *MockQuerier) GetMetricsTimeseriesByDay(ctx context.Context, days int32) ([]repository.GetMetricsTimeseriesByDayRow, error) {
	if m.GetMetricsTimeseriesByDayFunc != nil {
		return m.GetMetricsTimeseriesByDayFunc(ctx, days)
	}

	return []repository.GetMetricsTimeseriesByDayRow{}, nil
}

func (m *MockQuerier) GetMetricsTimeseriesForLast24Hours(ctx context.Context) ([]repository.GetMetricsTimeseriesForLast24HoursRow, error) {
	if m.GetMetricsTimeseriesForLast24HoursFunc != nil {
		return m.GetMetricsTimeseriesForLast24HoursFunc(ctx)
	}

	return []repository.GetMetricsTimeseriesForLast24HoursRow{}, nil
}

func (m *MockQuerier) GetOSSummaryByDay(ctx context.Context, days int32) ([]repository.GetOSSummaryByDayRow, error) {
	if m.GetOSSummaryByDayFunc != nil {
		return m.GetOSSummaryByDayFunc(ctx, days)
	}

	return []repository.GetOSSummaryByDayRow{}, nil
}

func (m *MockQuerier) GetOSSummaryForLast24Hours(ctx context.Context) ([]repository.GetOSSummaryForLast24HoursRow, error) {
	if m.GetOSSummaryForLast24HoursFunc != nil {
		return m.GetOSSummaryForLast24HoursFunc(ctx)
	}

	return []repository.GetOSSummaryForLast24HoursRow{}, nil
}

func (m *MockQuerier) GetReferrerSummaryByDay(ctx context.Context, days int32) ([]repository.GetReferrerSummaryByDayRow, error) {
	if m.GetReferrerSummaryByDayFunc != nil {
		return m.GetReferrerSummaryByDayFunc(ctx, days)
	}

	return []repository.GetReferrerSummaryByDayRow{}, nil
}

func (m *MockQuerier) GetReferrerSummaryForLast24Hours(ctx context.Context) ([]repository.GetReferrerSummaryForLast24HoursRow, error) {
	if m.GetReferrerSummaryForLast24HoursFunc != nil {
		return m.GetReferrerSummaryForLast24HoursFunc(ctx)
	}

	return []repository.GetReferrerSummaryForLast24HoursRow{}, nil
}

func (m *MockQuerier) GetRoutesSummaryByDay(ctx context.Context, days int32) ([]repository.GetRoutesSummaryByDayRow, error) {
	if m.GetRoutesSummaryByDayFunc != nil {
		return m.GetRoutesSummaryByDayFunc(ctx, days)
	}

	return []repository.GetRoutesSummaryByDayRow{}, nil
}

func (m *MockQuerier) GetRoutesSummaryForLast24Hours(ctx context.Context) ([]repository.GetRoutesSummaryForLast24HoursRow, error) {
	if m.GetRoutesSummaryForLast24HoursFunc != nil {
		return m.GetRoutesSummaryForLast24HoursFunc(ctx)
	}

	return []repository.GetRoutesSummaryForLast24HoursRow{}, nil
}

func (m *MockQuerier) GetSubcategories(ctx context.Context) ([]string, error) {
	if m.GetSubcategoriesFunc != nil {
		return m.GetSubcategoriesFunc(ctx)
	}

	return []string{}, nil
}

func (m *MockQuerier) GetTotalCheasheetsCount(ctx context.Context) (int64, error) {
	if m.GetTotalCheasheetsCountFunc != nil {
		return m.GetTotalCheasheetsCountFunc(ctx)
	}

	return 0, nil
}

func (m *MockQuerier) GetTotalClicksAndDownloads(ctx context.Context) (repository.GetTotalClicksAndDownloadsRow, error) {
	if m.GetTotalClicksAndDownloadsFunc != nil {
		return m.GetTotalClicksAndDownloadsFunc(ctx)
	}

	return repository.GetTotalClicksAndDownloadsRow{}, nil
}

func (m *MockQuerier) GetTotalImageSize(ctx context.Context) (repository.GetTotalImageSizeRow, error) {
	if m.GetTotalImageSizeFunc != nil {
		return m.GetTotalImageSizeFunc(ctx)
	}

	return repository.GetTotalImageSizeRow{}, nil
}

func (m *MockQuerier) GetTotalViewsAndVisitors(ctx context.Context) (repository.GetTotalViewsAndVisitorsRow, error) {
	if m.GetTotalViewsAndVisitorsFunc != nil {
		return m.GetTotalViewsAndVisitorsFunc(ctx)
	}

	return repository.GetTotalViewsAndVisitorsRow{}, nil
}

func (m *MockQuerier) ListCheatsheets(ctx context.Context, arg repository.ListCheatsheetsParams) ([]repository.ListCheatsheetsRow, error) {
	if m.ListCheatsheetsFunc != nil {
		return m.ListCheatsheetsFunc(ctx, arg)
	}

	return []repository.ListCheatsheetsRow{}, nil
}

func (m *MockQuerier) StoreEvent(ctx context.Context, arg repository.StoreEventParams) error {
	if m.StoreEventFunc != nil {
		return m.StoreEventFunc(ctx, arg)
	}

	return nil
}

func (m *MockQuerier) StorePageview(ctx context.Context, arg repository.StorePageviewParams) error {
	if m.StorePageviewFunc != nil {
		return m.StorePageviewFunc(ctx, arg)
	}

	return nil
}

func (m *MockQuerier) UpdateCheatsheet(ctx context.Context, arg repository.UpdateCheatsheetParams) error {
	if m.UpdateCheatsheetFunc != nil {
		return m.UpdateCheatsheetFunc(ctx, arg)
	}

	return nil
}
