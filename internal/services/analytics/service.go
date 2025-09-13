package analytics

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/domain/dtos"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/geo"
	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	"github.com/mssola/user_agent"
)

type AnalyticsService interface {
	RecordPageView(ctx context.Context, details dtos.PageviewRequest) error
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
