package geo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGeoSdk(t *testing.T) {
	tests := []struct {
		name             string
		basePath         string
		apiKey           string
		setEnv           bool
		expectConfigured bool
	}{
		{
			name:             "creates SDK with valid config",
			basePath:         "https://ipinfo.io",
			apiKey:           "test_token_123",
			setEnv:           true,
			expectConfigured: true,
		},
		{
			name:             "creates empty SDK when base path is not set",
			basePath:         "",
			apiKey:           "test_token",
			setEnv:           true,
			expectConfigured: false,
		},
		{
			name:             "creates empty SDK when api key is not set",
			basePath:         "https://ipinfo.io",
			apiKey:           "",
			setEnv:           true,
			expectConfigured: false,
		},
		{
			name:             "creates empty SDK when nothing is set",
			setEnv:           false,
			expectConfigured: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				if tt.basePath != "" {
					t.Setenv("IP_INFO_BASE_PATH", tt.basePath)
				}

				if tt.apiKey != "" {
					t.Setenv("IP_INFO_TOKEN", tt.apiKey)
				}
			}

			sdk := NewGeoSdk()

			if tt.expectConfigured {
				if sdk.apiKey == "" || sdk.basePath == "" {
					t.Errorf("Expected SDK to be configured, but basePath=%q, apiKey=%q", sdk.basePath, sdk.apiKey)
				}
			} else {
				if sdk.apiKey != "" || sdk.basePath != "" {
					t.Errorf("Expected SDK to be empty, but basePath=%q, apiKey=%q", sdk.basePath, sdk.apiKey)
				}
			}
		})
	}
}

func TestFetchCountry(t *testing.T) {
	tests := []struct {
		name            string
		ip              string
		mockResponse    IpInfoResponse
		mockStatusCode  int
		sdkConfigured   bool
		expectedCountry string
		expectError     bool
	}{
		{
			name: "successfully fetches country",
			ip:   "8.8.8.8",
			mockResponse: IpInfoResponse{
				Country: "United States",
			},
			mockStatusCode:  http.StatusOK,
			sdkConfigured:   true,
			expectedCountry: "United States",
			expectError:     false,
		},
		{
			name: "handles different country",
			ip:   "1.1.1.1",
			mockResponse: IpInfoResponse{
				Country: "Australia",
			},
			mockStatusCode:  http.StatusOK,
			sdkConfigured:   true,
			expectedCountry: "Australia",
			expectError:     false,
		},
		{
			name:          "returns error when SDK is not configured",
			ip:            "8.8.8.8",
			sdkConfigured: false,
			expectError:   true,
		},
		{
			name:           "handles API error status",
			ip:             "8.8.8.8",
			mockStatusCode: http.StatusUnauthorized,
			sdkConfigured:  true,
			expectError:    true,
		},
		{
			name:           "handles API not found",
			ip:             "invalid",
			mockStatusCode: http.StatusNotFound,
			sdkConfigured:  true,
			expectError:    true,
		},
		{
			name:           "handles API rate limit",
			ip:             "8.8.8.8",
			mockStatusCode: http.StatusTooManyRequests,
			sdkConfigured:  true,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates a mock HTTP server */
			var server *httptest.Server

			if tt.sdkConfigured {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					/* Verify the request */
					expectedPath := "/" + tt.ip
					if r.URL.Path != expectedPath {
						t.Errorf("Expected path %q, got %q", expectedPath, r.URL.Path)
					}

					/* Check token is in query params */
					token := r.URL.Query().Get("token")
					if token == "" {
						t.Error("Expected token in query params")
					}

					/* Write mock response */
					w.WriteHeader(tt.mockStatusCode)
					if tt.mockStatusCode == http.StatusOK {
						json.NewEncoder(w).Encode(tt.mockResponse)
					}
				}))

				defer server.Close()
			}

			/* Creates a Geo SDK instance pointing to the mock server */
			var sdk *IpInfoSdk
			if tt.sdkConfigured {
				sdk = &IpInfoSdk{
					basePath: server.URL,
					apiKey:   "test_token",
				}
			} else {
				sdk = &IpInfoSdk{}
			}

			country, err := sdk.FetchCountry(tt.ip)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError && country != tt.expectedCountry {
				t.Errorf("FetchCountry(%q) = %q; want %q", tt.ip, country, tt.expectedCountry)
			}
		})
	}
}
