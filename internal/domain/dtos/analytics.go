package dtos

type PageviewRequest struct {
	Route     string `json:"route" binding:"required"`
	Referrer  string `json:"referrer"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}
