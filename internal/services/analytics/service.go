package analytics

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/entities"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/geo"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	"github.com/mssola/user_agent"
)

type AnalyticsService interface {
	GetMetricsOverview(ctx context.Context, period string) (*dtos.MetricsOverviewResponse, error)
	GetDeviceStats(ctx context.Context, period string) (*dtos.DeviceStatsResponse, error)
	GetBrowserStats(ctx context.Context, period string) (*dtos.BrowserStatsResponse, error)
	GetOperatingSystemsStats(ctx context.Context, period string) (*dtos.OperatingSystemStatsResponse, error)
	GetReferrerStats(ctx context.Context, period string) (*dtos.ReferrerStatsResponse, error)
	GetRoutesStats(ctx context.Context, period string) (*dtos.RoutesStatsResponse, error)
	RecordPageView(ctx context.Context, details dtos.PageviewRequest) error
	RecordEvent(ctx context.Context, details dtos.EventRequest) error
}

type analyticsService struct {
	repo   *repository.Queries
	geoSdk *geo.IpInfoSdk
}

func NewAnalyticsService(repo *repository.Queries) AnalyticsService {
	geoSdk := geo.NewGeoSdk()

	return &analyticsService{
		repo:   repo,
		geoSdk: geoSdk,
	}
}

/**
 * Get routes metrics by period - 24h, 7d, 30d etc.
 * @param period string
 * @return *dtos.RoutesStatsResponse, error
 */
func (s *analyticsService) GetRoutesStats(ctx context.Context, period string) (*dtos.RoutesStatsResponse, error) {
	// check if period is valid
	periodData, ok := entities.PeriodConfigs[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var data interface{}
	var err error

	switch period {
	case "24h":
		data, err = s.repo.GetRoutesSummaryForLast24Hours(ctx)

	case "7d", "30d", "3m", "6m", "12m":
		data, err = s.repo.GetRoutesSummaryByDay(ctx, int32(periodData.Days))

	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch referrer stats: %w", err)
	}

	routes := buildRoutesStats(data)
	startDate, endDate := getStartAndEndDatesForPeriod(int64(periodData.Days))

	stats := &dtos.RoutesStatsResponse{
		Period:              period,
		StartDate:           startDate,
		EndDate:             endDate,
		TotalViews:          routes.TotalViews,
		TotalUniqueVisitors: routes.TotalUniqueVisitors,
		Routes:              routes.Routes,
	}

	return stats, nil
}

/**
 * Build routes stats from database rows
 * @param data interface{}
 * @return *entities.ReferrerStats
 */
func buildRoutesStats(data interface{}) *entities.RoutesStats {
	var routes []dtos.DataStat
	var totalViews, uniqueVisitors int64

	switch rows := data.(type) {
	case []repository.GetRoutesSummaryByDayRow:
		routes = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			routes = append(routes, dtos.DataStat{
				Name:     r.Pathname,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})

			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}
	case []repository.GetRoutesSummaryForLast24HoursRow:
		routes = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			routes = append(routes, dtos.DataStat{
				Name:     r.Pathname,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})

			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}
	}

	return &entities.RoutesStats{
		TotalViews:          totalViews,
		TotalUniqueVisitors: uniqueVisitors,
		Routes:              routes,
	}

}

/**
 * Get referrer metrics by period - 24h, 7d, 30d etc.
 * @param period string
 * @return *dtos.ReferrerStatsResponse, error
 */
func (s *analyticsService) GetReferrerStats(ctx context.Context, period string) (*dtos.ReferrerStatsResponse, error) {
	// check if period is valid
	periodData, ok := entities.PeriodConfigs[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var data interface{}
	var err error

	switch period {
	case "24h":
		data, err = s.repo.GetReferrerSummaryForLast24Hours(ctx)

	case "7d", "30d", "3m", "6m", "12m":
		data, err = s.repo.GetReferrerSummaryByDay(ctx, int32(periodData.Days))

	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch referrer stats: %w", err)
	}

	referrers := buildReferrerStats(data)
	startDate, endDate := getStartAndEndDatesForPeriod(int64(periodData.Days))

	stats := &dtos.ReferrerStatsResponse{
		Period:              period,
		StartDate:           startDate,
		EndDate:             endDate,
		TotalViews:          referrers.TotalViews,
		TotalUniqueVisitors: referrers.TotalUniqueVisitors,
		Referrers:           referrers.Referrers,
	}

	return stats, nil
}

/**
 * Build referrer stats from database rows
 * @param data interface{}
 * @return *entities.ReferrerStats
 */
func buildReferrerStats(data interface{}) *entities.ReferrerStats {
	var referrers []dtos.DataStat
	var totalViews, uniqueVisitors int64

	switch rows := data.(type) {
	case []repository.GetReferrerSummaryByDayRow:
		referrers = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			referrers = append(referrers, dtos.DataStat{
				Name:     r.Referrer.String,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})
			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}

	case []repository.GetReferrerSummaryForLast24HoursRow:
		referrers = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			referrers = append(referrers, dtos.DataStat{
				Name:     r.Referrer.String,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})
			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}
	}

	return &entities.ReferrerStats{
		TotalViews:          totalViews,
		TotalUniqueVisitors: uniqueVisitors,
		Referrers:           referrers,
	}
}

/**
 * Get operating systems metrics by period - 24h, 7d, 30d etc.
 * @param period string
 * @return *dtos.OperatingSystemStatsResponse, error
 */
func (s *analyticsService) GetOperatingSystemsStats(ctx context.Context, period string) (*dtos.OperatingSystemStatsResponse, error) {
	// check if period is valid
	periodData, ok := entities.PeriodConfigs[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var data interface{}
	var err error

	switch period {
	case "24h":
		data, err = s.repo.GetOSSummaryForLast24Hours(ctx)
	case "7d", "30d", "3m", "6m", "12m":
		data, err = s.repo.GetOSSummaryByDay(ctx, int32(periodData.Days))
	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch operating systems data: %w", err)
	}

	osData := buildOperatingSystemStats(data)
	startDate, endDate := getStartAndEndDatesForPeriod(int64(periodData.Days))

	stats := &dtos.OperatingSystemStatsResponse{
		Period:              period,
		StartDate:           startDate,
		EndDate:             endDate,
		TotalViews:          osData.TotalViews,
		TotalUniqueVisitors: osData.TotalUniqueVisitors,
		OperatingSystems:    osData.OperatingSystems,
	}

	return stats, nil
}

/**
 * Build os stats from database rows
 * @param data interface{}
 * @return *entities.OperatingSystemStats
 */
func buildOperatingSystemStats(data interface{}) *entities.OperatingSystemStats {
	var osData []dtos.DataStat
	var totalViews, uniqueVisitors int64

	switch rows := data.(type) {
	case []repository.GetOSSummaryByDayRow:
		osData = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			osData = append(osData, dtos.DataStat{
				Name:     r.OsGroup,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})

			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}

	case []repository.GetOSSummaryForLast24HoursRow:
		osData = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			osData = append(osData, dtos.DataStat{
				Name:     r.OsGroup,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})

			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}
	}

	return &entities.OperatingSystemStats{
		TotalViews:          totalViews,
		TotalUniqueVisitors: uniqueVisitors,
		OperatingSystems:    osData,
	}
}

/**
 * Get browsers metrics by period - 24h, 7d, 30d etc.
 * @param period string
 * @return *dtos.BrowserStatsResponse, error
 */
func (s analyticsService) GetBrowserStats(ctx context.Context, period string) (*dtos.BrowserStatsResponse, error) {
	// check if period is valid
	periodData, ok := entities.PeriodConfigs[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var data interface{}
	var err error

	switch period {
	case "24h":
		data, err = s.repo.GetBrowsersSummaryForLast24Hours(ctx)
	case "7d", "30d", "3m", "6m", "12m":
		data, err = s.repo.GetBrowsersSummaryByDay(ctx, int32(periodData.Days))
	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch pageviews data: %w", err)
	}

	browserData := buildBrowserStats(data)
	startDate, endDate := getStartAndEndDatesForPeriod(int64(periodData.Days))

	stats := &dtos.BrowserStatsResponse{
		Period:              period,
		StartDate:           startDate,
		EndDate:             endDate,
		TotalViews:          browserData.TotalViews,
		TotalUniqueVisitors: browserData.TotalUniqueVisitors,
		Browsers:            browserData.Browsers,
	}

	return stats, nil
}

/**
 * Build browser stats from database rows
 * @param data interface{}
 * @return *entities.BrowserStats
 */
func buildBrowserStats(data interface{}) *entities.BrowserStats {
	var browserData []dtos.DataStat
	var totalViews, uniqueVisitors int64

	switch rows := data.(type) {
	case []repository.GetBrowsersSummaryByDayRow:
		browserData = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			browserData = append(browserData, dtos.DataStat{
				Name:     r.Browser.String,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})
			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors

		}
	case []repository.GetBrowsersSummaryForLast24HoursRow:
		browserData = make([]dtos.DataStat, 0, len(rows))
		for _, r := range rows {
			browserData = append(browserData, dtos.DataStat{
				Name:     r.Browser.String,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})
			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors

		}
	}

	return &entities.BrowserStats{
		TotalUniqueVisitors: uniqueVisitors,
		TotalViews:          totalViews,
		Browsers:            browserData,
	}

}

/**
 * Get devices metrics by period - 24h, 7d, 30d etc.
 * @param period string
 * @return *dtos.DeviceStatsResponse, error
 */
func (s *analyticsService) GetDeviceStats(ctx context.Context, period string) (*dtos.DeviceStatsResponse, error) {
	// check if period is valid
	periodData, ok := entities.PeriodConfigs[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var data interface{}
	var err error

	switch period {
	case "24h":
		data, err = s.repo.GetBrowsersSummaryForLast24Hours(ctx)
	case "7d", "30d", "3m", "6m", "12m":
		data, err = s.repo.GetDevicesSummaryByDay(ctx, int32(periodData.Days))
	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch pageviews data: %w", err)
	}

	devicesData := buildDevicesStats(data)

	// Calculate start and end dates
	var startDate, endDate time.Time
	length := len(devicesData.Intervals)

	if length > 0 {
		startDate = devicesData.Intervals[0].Date
		endDate = devicesData.Intervals[length-1].Date
	}

	stats := &dtos.DeviceStatsResponse{
		Period:               period,
		StartDate:            startDate,
		EndDate:              endDate,
		TotalMobileViews:     devicesData.TotalMobileViews,
		TotalMobileVisitors:  devicesData.TotalMobileVisitors,
		TotalDesktopViews:    devicesData.TotalDesktopViews,
		TotalDesktopVisitors: devicesData.TotalDesktopVisitors,
		Intervals:            devicesData.Intervals,
	}

	return stats, nil
}

/**
 * Build devices stats from database rows
 * @param data interface{}
 * @return *entities.BrowserStats
 */
func buildDevicesStats(data interface{}) *entities.DeviceStats {
	var devicesData []dtos.DeviceStat
	var totalMobileViews, totalMobileVisitors, totalDesktopViews, totalDesktopVisitors int64

	switch rows := data.(type) {
	case []repository.GetDevicesSummaryByDayRow:
		devicesData = make([]dtos.DeviceStat, 0, len(rows))
		for _, r := range rows {
			devicesData = append(devicesData, dtos.DeviceStat{
				Date:            r.Date.Time,
				MobileViews:     r.MobileViews,
				MobileVisitors:  r.MobileVisitors,
				DesktopViews:    r.DesktopViews,
				DesktopVisitors: r.DesktopVisitors,
			})
			totalMobileViews += r.MobileViews
			totalMobileVisitors += r.MobileVisitors
			totalDesktopViews += r.DesktopViews
			totalDesktopVisitors += r.DesktopVisitors
		}

	case []repository.GetDevicesSummaryForLast24HoursRow:
		devicesData = make([]dtos.DeviceStat, 0, len(rows))
		for _, r := range rows {
			devicesData = append(devicesData, dtos.DeviceStat{
				Date:            r.Hour.Time,
				MobileViews:     r.MobileViews,
				MobileVisitors:  r.MobileVisitors,
				DesktopViews:    r.DesktopViews,
				DesktopVisitors: r.DesktopVisitors,
			})
			totalMobileViews += r.MobileViews
			totalMobileVisitors += r.MobileVisitors
			totalDesktopViews += r.DesktopViews
			totalDesktopVisitors += r.DesktopVisitors
		}
	}

	return &entities.DeviceStats{
		TotalMobileViews:     totalMobileViews,
		TotalMobileVisitors:  totalMobileVisitors,
		TotalDesktopViews:    totalDesktopViews,
		TotalDesktopVisitors: totalDesktopVisitors,
		Intervals:            devicesData,
	}
}

/**
 * Get page views metrics by period - 24h, 7d, 30d etc.
 * @param period string
 * @return *dtos.MetricsOverviewResponse, error
 */
func (s *analyticsService) GetMetricsOverview(ctx context.Context, period string) (*dtos.MetricsOverviewResponse, error) {

	// check if period is valid
	periodData, ok := entities.PeriodConfigs[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var data interface{}
	var err error

	switch period {
	case "24h":
		data, err = s.repo.GetMetricsTimeseriesForLast24Hours(ctx)
	case "7d", "30d", "3m", "6m", "12m":
		data, err = s.repo.GetMetricsTimeseriesByDay(ctx, int32(periodData.Days))
	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch pageviews data: %w", err)
	}

	seriesData := buildMetricsOverview(data)

	// Calculate start and end dates
	var startDate, endDate time.Time
	length := len(seriesData.Intervals)

	if length > 0 {
		startDate = seriesData.Intervals[0].Date
		endDate = seriesData.Intervals[length-1].Date
	}

	stats := &dtos.MetricsOverviewResponse{
		Period:              period,
		StartDate:           startDate,
		EndDate:             endDate,
		TotalViews:          seriesData.TotalViews,
		TotalUniqueVisitors: seriesData.TotalUniqueVisitors,
		Intervals:           seriesData.Intervals,
	}

	return stats, nil
}

/**
 * Build metrics stats from database rows
 * @param data interface{}
 * @return *entities.BrowserStats
 */
func buildMetricsOverview(data interface{}) *entities.MetricswSeries {
	var metricsData []dtos.MetricsOverview
	var totalViews, uniqueVisitors int64

	switch rows := data.(type) {
	case []repository.GetMetricsTimeseriesByDayRow:
		metricsData = make([]dtos.MetricsOverview, 0, len(rows))
		for _, r := range rows {
			metricsData = append(metricsData, dtos.MetricsOverview{
				Date:     r.Date.Time,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})

			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}
	case []repository.GetMetricsTimeseriesForLast24HoursRow:
		metricsData = make([]dtos.MetricsOverview, 0, len(rows))
		for _, r := range rows {
			metricsData = append(metricsData, dtos.MetricsOverview{
				Date:     r.Hour.Time,
				Views:    r.Views,
				Visitors: r.UniqueVisitors,
			})

			totalViews += r.Views
			uniqueVisitors += r.UniqueVisitors
		}
	}

	return &entities.MetricswSeries{
		TotalViews:          totalViews,
		TotalUniqueVisitors: uniqueVisitors,
		Intervals:           metricsData,
	}
}

/**
 * Record a page view
 * @param details dtos.PageviewRequest
 * @return error
 */
func (s *analyticsService) RecordPageView(ctx context.Context, details dtos.PageviewRequest) error {
	if details.Route == "" || details.UserAgent == "" || details.IpAddress == "" {
		return nil
	}

	browser, os, device := parseUserAgent(details.UserAgent)

	country, err := s.geoSdk.FetchCountry(details.IpAddress)
	if err != nil {
		return fmt.Errorf("failed to lookup geo info: %w", err)
	}

	pageViewParams := repository.StorePageviewParams{
		Pathname:  details.Route,
		Browser:   utils.PgText(browser),
		Os:        utils.PgText(os),
		Device:    utils.PgText(device),
		HashedIp:  hashIP(details.IpAddress),
		UserAgent: details.UserAgent,
		Country:   utils.PgText(country),
		Referrer:  utils.PgText(details.Referrer),
	}

	if err := s.repo.StorePageview(ctx, pageViewParams); err != nil {
		return fmt.Errorf("failed to record pageview: %w", err)
	}

	return nil
}

/**
 * Record an event - click, download etc
 * @param details dtos.EventRequest
 * @return error
 */
func (s *analyticsService) RecordEvent(ctx context.Context, details dtos.EventRequest) error {
	if details.IpAddress == "" {
		return nil
	}

	/* get cheatsheet id from database */
	data, err := s.repo.GetCheatsheetBySlug(ctx, details.CheatsheetSlug)
	if err != nil {
		return fmt.Errorf("failed to fetch cheatsheet id")
	}

	eventsParams := repository.StoreEventParams{
		CheatsheetID: data.ID,
		EventType:    repository.EventType(details.EventType),
		Pathname:     details.Route,
		HashedIp:     hashIP(details.IpAddress),
	}

	if err := s.repo.StoreEvent(ctx, eventsParams); err != nil {
		return fmt.Errorf("failed to record event: %w", err)
	}

	return nil
}

/**
 * Parse user agent string to extract browser, OS, and device type
 * @param uaString string
 * @return browser, os, device string
 */
func parseUserAgent(uaString string) (browser, os, device string) {
	ua := user_agent.New(uaString)

	// Browser
	browserName, _ := ua.Browser()

	// OS
	os = ua.OS()

	// Device
	if ua.Mobile() {
		device = "mobile"
	} else {
		device = "desktop"
	}

	return browserName, os, device
}

/**
 * Hash IP address with a salt from environment variable
 * @param ip string
 * @return hashed IP string
 */
func hashIP(ip string) string {
	salt := utils.GetEnv("IP_HASH_SALT", "")
	if salt == "" {
		log.Fatal("IP_HASH_SALT is not set in environment")
	}

	h := sha256.New()
	h.Write([]byte(ip + salt))
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * Get start and end dates for a given period in days
 * @param days int64
 * @return startDate, endDate time.Time, error
 */
func getStartAndEndDatesForPeriod(days int64) (time.Time, time.Time) {
	now := time.Now().UTC()
	var startDate, endDate time.Time

	if days == 0 {
		endDate = now.Truncate(time.Hour)
		startDate = endDate.Add(-24 * time.Hour)
	} else {
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		startDate = endDate.AddDate(0, 0, -int(days))
	}

	return startDate, endDate
}
