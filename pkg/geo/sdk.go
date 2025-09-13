package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type IpInfoSdk struct {
	basePath string
	apiKey   string
}

func NewGeoSdk() *IpInfoSdk {
	basePath := utils.GetEnv("IP_INFO_BASE_PATH", "")
	apiKey := utils.GetEnv("IP_INFO_TOKEN", "")

	if basePath == "" || apiKey == "" {
		log.Info().Msg("IP Info base path or API key is not set. Geo lookup will be disabled.")
		return &IpInfoSdk{}
	}

	return &IpInfoSdk{
		basePath: basePath,
		apiKey:   apiKey,
	}
}

func (s *IpInfoSdk) FetchCountry(ip string) (string, error) {
	if s.basePath == "" || s.apiKey == "" {
		return "", fmt.Errorf("IP Info SDK not configured")
	}

	url := fmt.Sprintf("%s/%s?token=%s", s.basePath, ip, s.apiKey)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to call ipinfo API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ipinfo API returned status %d", resp.StatusCode)
	}

	var data IpInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode ipinfo response: %w", err)
	}

	return data.Country, nil
}
