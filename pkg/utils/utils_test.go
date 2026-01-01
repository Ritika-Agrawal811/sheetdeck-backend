package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestPgText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected pgtype.Text
	}{
		{
			name:     "empty string returns invalid",
			input:    "",
			expected: pgtype.Text{Valid: false},
		},
		{
			name:     "non-empty string returns valid",
			input:    "hello",
			expected: pgtype.Text{Valid: true, String: "hello"},
		},
		{
			name:     "whitespace string returns valid",
			input:    "  ",
			expected: pgtype.Text{Valid: true, String: "  "},
		},
		{
			name:     "string with special characters",
			input:    "hello@world.com",
			expected: pgtype.Text{Valid: true, String: "hello@world.com"},
		},
		{
			name:     "multiline string",
			input:    "line1 \n line2",
			expected: pgtype.Text{Valid: true, String: "line1 \n line2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgText(tt.input)

			if result.Valid != tt.expected.Valid {
				t.Errorf("Expect PgText(%q).Valid = %v; but got %v", tt.input, tt.expected.Valid, result.Valid)
			}

			if result.String != tt.expected.String {
				t.Errorf("Expect PgText(%q).String = %q; but got %q", tt.input, tt.expected.String, result.String)
			}
		})
	}
}

func TestPgInt8(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected pgtype.Int8
	}{
		{
			name:     "zero value",
			input:    0,
			expected: pgtype.Int8{Int64: 0, Valid: true},
		},
		{
			name:     "positive number",
			input:    42,
			expected: pgtype.Int8{Int64: 42, Valid: true},
		},
		{
			name:     "negative number",
			input:    -100,
			expected: pgtype.Int8{Int64: -100, Valid: true},
		},
		{
			name:     "max int64",
			input:    9223372036854775807,
			expected: pgtype.Int8{Int64: 9223372036854775807, Valid: true},
		},
		{
			name:     "min int64",
			input:    -9223372036854775808,
			expected: pgtype.Int8{Int64: -9223372036854775808, Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgInt8(tt.input)

			if result.Valid != tt.expected.Valid {
				t.Errorf("Expect PgInt8(%d).Valid = %v; but got %v", tt.input, tt.expected.Valid, result.Valid)
			}

			if result.Int64 != tt.expected.Int64 {
				t.Errorf("Expect PgInt8(%d).Int64 = %d; but got %d", tt.input, tt.expected.Int64, result.Int64)
			}
		})
	}
}

func TestStringToUUID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
		wantError bool
	}{
		{
			name:      "valid UUID",
			input:     "550e8400-e29b-41d4-a716-446655440000",
			wantValid: true,
			wantError: false,
		},
		{
			name:      "valid UUID without hyphens",
			input:     "550e8400e29b41d4a716446655440000",
			wantValid: true,
			wantError: false,
		},
		{
			name:      "valid UUID in uppercase",
			input:     "550E8400-E29B-41D4-A716-446655440000",
			wantValid: true,
			wantError: false,
		},
		{
			name:      "nil UUID (all zeros)",
			input:     "00000000-0000-0000-0000-000000000000",
			wantValid: true,
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "UUID with too short length",
			input:     "550e8400-e29b-41d4",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "UUID with invalid characters",
			input:     "550e8400-e29b-41d4-a716-44665544000g",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "UUID with wrong format",
			input:     "not-a-uuid-at-all",
			wantValid: false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StringToUUID(tt.input)

			// Check error expectation
			if tt.wantError && err == nil {
				t.Errorf("StringToUUID(%q) expected error but got none", tt.input)
			}

			if !tt.wantError && err != nil {
				t.Errorf("StringToUUID(%q) unexpected error: %v", tt.input, err)
			}

			// Check Valid field
			if result.Valid != tt.wantValid {
				t.Errorf("Expect StringToUUID(%q).Valid = %v, but got %v", tt.input, tt.wantValid, result.Valid)
			}

			// For valid UUIDs, verify the bytes were copied correctly
			if tt.wantValid && !tt.wantError {
				// Parse the original UUID to compare bytes
				originalUUID, _ := uuid.Parse(tt.input)

				for i := 0; i < 16; i++ {
					if result.Bytes[i] != originalUUID[i] {
						t.Errorf("StringToUUID(%q) byte mismatch at index %d: got %v; want %v",
							tt.input, i, result.Bytes[i], originalUUID[i])
						break
					}
				}
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		setEnv   bool
		expected string
	}{
		{
			name:     "returns env when set",
			key:      "TEST_KEY_1",
			fallback: "default",
			envValue: "from_env",
			setEnv:   true,
			expected: "from_env",
		},
		{
			name:     "returns fallback when env is not set",
			key:      "TEST_KEY_NOT_SET",
			fallback: "fallback_value",
			setEnv:   false,
			expected: "fallback_value",
		},
		{
			name:     "returns empty string env when explicitly set to empty",
			key:      "TEST_KEY_EMPTY",
			fallback: "fallback",
			envValue: "",
			setEnv:   true,
			expected: "",
		},
		{
			name:     "returns fallback when key is empty string",
			key:      "",
			fallback: "fallback",
			setEnv:   false,
			expected: "fallback",
		},
		{
			name:     "returns empty fallback when env is not set",
			key:      "TEST_KEY_2",
			fallback: "",
			setEnv:   false,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.key, tt.envValue)
			}

			result := GetEnv(tt.key, tt.fallback)

			if result != tt.expected {
				t.Errorf("Expect GetEnv(%q, %q) = %q; but got %q", tt.key, tt.fallback, tt.expected, result)
			}
		})
	}
}
