package entities

import (
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
)

type PeriodConfig struct {
	Days int
}

var PeriodConfigs = map[string]PeriodConfig{
	"24h": {Days: 0},
	"7d":  {Days: 6},
	"30d": {Days: 29},
	"3m":  {Days: 89},
	"6m":  {Days: 179},
	"12m": {Days: 364},
}

type PageviewSeries struct {
	TotalViews          int64               `json:"total_views"`
	TotalUniqueVisitors int64               `json:"total_unique_visitors"`
	Intervals           []dtos.PageviewStat `json:"intervals"`
}

type DeviceStats struct {
	TotalMobileViews     int64             `json:"total_mobile_views"`
	TotalMobileVisitors  int64             `json:"total_mobile_visitors"`
	TotalDesktopViews    int64             `json:"total_desktop_views"`
	TotalDesktopVisitors int64             `json:"total_desktop_visitors"`
	Intervals            []dtos.DeviceStat `json:"Intervals"`
}

type BrowserStats struct {
	TotalViews          int64              `json:"total_views"`
	TotalUniqueVisitors int64              `json:"total_unique_visitors"`
	Browsers            []dtos.BrowserStat `json:"browsers"`
}

type OperatingSystemStats struct {
	TotalViews          int64                      `json:"total_views"`
	TotalUniqueVisitors int64                      `json:"total_unique_visitors"`
	OperatingSystems    []dtos.OperatingSystemStat `json:"operating_systems"`
}
