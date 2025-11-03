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
	GetPageviewsStats(ctx context.Context, period string) (*dtos.PageviewStatsResponse, error)
	GetDeviceStats(ctx context.Context, period string) (*dtos.DeviceStatsResponse, error)
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

	var devicesData *entities.DeviceStats
	var err error

	switch period {
	case "24h":
		devicesData, err = s.getDevicesStatsForLast24Hours(ctx)
	case "7d", "30d", "3m", "6m", "12m":
		devicesData, err = s.getDevicesStatsByDay(ctx, int32(periodData.Days))
	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch pageviews data: %w", err)
	}

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
 * Get devices stats by days
 * @param days int32
 * @return *entities.DeviceStats, error
 */
func (s *analyticsService) getDevicesStatsByDay(ctx context.Context, days int32) (*entities.DeviceStats, error) {
	devicesData, err := s.repo.GetDevicesSummaryByDay(ctx, days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page views stats: %w", err)
	}

	devices := make([]dtos.DeviceStat, 0, len(devicesData))
	var totalMobileViews, totalMobileVisitors, totalDesktopViews, totalDesktopVisitors int64

	for _, data := range devicesData {
		devices = append(devices, dtos.DeviceStat{
			Date:            data.Date.Time,
			MobileViews:     data.MobileViews,
			MobileVisitors:  data.MobileVisitors,
			DesktopViews:    data.DesktopViews,
			DesktopVisitors: data.DesktopVisitors,
		})

		totalMobileViews += data.MobileViews
		totalMobileVisitors += data.MobileVisitors
		totalDesktopViews += data.DesktopViews
		totalDesktopVisitors += data.DesktopVisitors

	}

	response := &entities.DeviceStats{
		TotalMobileViews:     totalMobileViews,
		TotalMobileVisitors:  totalMobileVisitors,
		TotalDesktopViews:    totalDesktopViews,
		TotalDesktopVisitors: totalDesktopVisitors,
		Intervals:            devices,
	}

	return response, nil
}

/**
 * Get devices stats for last 24 hours
 * @return *entities.DeviceStats, error
 */
func (s *analyticsService) getDevicesStatsForLast24Hours(ctx context.Context) (*entities.DeviceStats, error) {
	devicesData, err := s.repo.GetDevicesSummaryForLast24Hours(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page views stats: %w", err)
	}

	devices := make([]dtos.DeviceStat, 0, len(devicesData))
	var totalMobileViews, totalMobileVisitors, totalDesktopViews, totalDesktopVisitors int64

	for _, data := range devicesData {
		devices = append(devices, dtos.DeviceStat{
			Date:            data.Hour.Time,
			MobileViews:     data.MobileViews,
			MobileVisitors:  data.MobileVisitors,
			DesktopViews:    data.DesktopViews,
			DesktopVisitors: data.DesktopVisitors,
		})

		totalMobileViews += data.MobileViews
		totalMobileVisitors += data.MobileVisitors
		totalDesktopViews += data.DesktopViews
		totalDesktopVisitors += data.DesktopVisitors

	}

	response := &entities.DeviceStats{
		TotalMobileViews:     totalMobileViews,
		TotalMobileVisitors:  totalMobileVisitors,
		TotalDesktopViews:    totalDesktopViews,
		TotalDesktopVisitors: totalDesktopVisitors,
		Intervals:            devices,
	}

	return response, nil
}

/**
 * Get page views metrics by period - 24h, 7d, 30d etc.
 * @param period string
 * @return *dtos.PageviewStatsResponse, error
 */
func (s *analyticsService) GetPageviewsStats(ctx context.Context, period string) (*dtos.PageviewStatsResponse, error) {

	// check if period is valid
	periodData, ok := entities.PeriodConfigs[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	var seriesData *entities.PageviewSeries
	var err error

	switch period {
	case "24h":
		seriesData, err = s.getPageviewSeriesForLast24Hours(ctx)
	case "7d", "30d", "3m", "6m", "12m":
		seriesData, err = s.getPageviewSeriesByDay(ctx, int32(periodData.Days))
	default:
		return nil, fmt.Errorf("invalid period")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch pageviews data: %w", err)
	}

	// Calculate start and end dates
	var startDate, endDate time.Time
	length := len(seriesData.Intervals)

	if length > 0 {
		startDate = seriesData.Intervals[0].Date
		endDate = seriesData.Intervals[length-1].Date
	}

	stats := &dtos.PageviewStatsResponse{
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
 * Get page views timeseries, total views and visitors for last 24 hours
 * @return *entities.PageviewSeries, error
 */
func (s *analyticsService) getPageviewSeriesForLast24Hours(ctx context.Context) (*entities.PageviewSeries, error) {
	timeseriesData, err := s.repo.GetPageviewTimeseriesForLast24Hours(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page views stats: %w", err)
	}

	intervals := make([]dtos.PageviewStat, 0, len(timeseriesData))
	var totalViews int64
	var uniqueVisitors int64

	for _, data := range timeseriesData {
		intervals = append(intervals, dtos.PageviewStat{
			Date:     data.Hour.Time,
			Views:    data.Views,
			Visitors: data.UniqueVisitors,
		})

		totalViews += data.Views
		uniqueVisitors += data.UniqueVisitors
	}

	response := &entities.PageviewSeries{
		TotalViews:          totalViews,
		TotalUniqueVisitors: uniqueVisitors,
		Intervals:           intervals,
	}

	return response, nil
}

/**
 * Get page views timeseries, total views and visitors by days
 * @param days int32
 * @return *entities.PageviewSeries, error
 */
func (s *analyticsService) getPageviewSeriesByDay(ctx context.Context, days int32) (*entities.PageviewSeries, error) {
	timeseriesData, err := s.repo.GetPageviewTimeseriesByDay(ctx, days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page views stats: %w", err)
	}

	intervals := make([]dtos.PageviewStat, 0, len(timeseriesData))
	var totalViews int64
	var uniqueVisitors int64

	for _, data := range timeseriesData {
		intervals = append(intervals, dtos.PageviewStat{
			Date:     data.Date.Time,
			Views:    data.Views,
			Visitors: data.UniqueVisitors,
		})

		totalViews += data.Views
		uniqueVisitors += data.UniqueVisitors
	}

	response := &entities.PageviewSeries{
		TotalViews:          totalViews,
		TotalUniqueVisitors: uniqueVisitors,
		Intervals:           intervals,
	}

	return response, nil
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
