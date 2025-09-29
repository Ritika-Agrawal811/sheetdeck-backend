package dtos

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
