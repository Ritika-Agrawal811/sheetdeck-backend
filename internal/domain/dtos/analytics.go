package dtos

import "time"

type PageviewRequest struct {
	Route     string `json:"route" binding:"required"`
	Referrer  string `json:"referrer"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

type EventRequest struct {
	Route          string `json:"route" binding:"required"`
	CheatsheetSlug string `json:"cheatsheet_slug" binding:"required"`
	EventType      string `json:"event_type" binding:"required"`
	IpAddress      string `json:"ip_address"`
}

type PageviewStatsResponse struct {
	Period              string         `json:"period"`
	StartDate           time.Time      `json:"start_date"`
	EndDate             time.Time      `json:"end_date"`
	TotalViews          int64          `json:"total_views"`
	TotalUniqueVisitors int64          `json:"total_unique_visitors"`
	Intervals           []PageviewStat `json:"intervals"`
}

type PageviewStat struct {
	Date     time.Time `json:"date"`
	Views    int64     `json:"views"`
	Visitors int64     `json:"visitors"`
}

type DeviceStatsResponse struct {
	Period               string       `json:"period"`
	StartDate            time.Time    `json:"start_date"`
	EndDate              time.Time    `json:"end_date"`
	TotalMobileViews     int64        `json:"total_mobile_views"`
	TotalMobileVisitors  int64        `json:"total_mobile_visitors"`
	TotalDesktopViews    int64        `json:"total_desktop_views"`
	TotalDesktopVisitors int64        `json:"total_desktop_visitors"`
	Intervals            []DeviceStat `json:"intervals"`
}

type DeviceStat struct {
	Date            time.Time `json:"date"`
	MobileViews     int64     `json:"mobile_views"`
	MobileVisitors  int64     `json:"mobile_visitors"`
	DesktopViews    int64     `json:"desktop_views"`
	DesktopVisitors int64     `json:"desktop_visitors"`
}

type BrowserStatsResponse struct {
	Period              string        `json:"period"`
	StartDate           time.Time     `json:"start_date"`
	EndDate             time.Time     `json:"end_date"`
	TotalViews          int64         `json:"total_views"`
	TotalUniqueVisitors int64         `json:"total_unique_visitors"`
	Browsers            []BrowserStat `json:"browsers"`
}

type BrowserStat struct {
	Browser  string `json:"browser"`
	Views    int64  `json:"views"`
	Visitors int64  `json:"visitors"`
}
