package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("PORT")
	os.Unsetenv("SORT_CODE")
	os.Unsetenv("COMPANY_NAME")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.SortCode != "987654" {
		t.Errorf("expected default sort code 987654, got %s", cfg.SortCode)
	}
	if cfg.CompanyName != "CICS Bank Sample Application" {
		t.Errorf("expected default company name, got %s", cfg.CompanyName)
	}
	if cfg.DatabaseURL == "" {
		t.Error("expected non-empty database URL")
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("SORT_CODE", "123456")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("SORT_CODE")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.SortCode != "123456" {
		t.Errorf("expected sort code 123456, got %s", cfg.SortCode)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envVal   string
		fallback string
		expected string
	}{
		{"env set", "TEST_KEY_A", "value_a", "fallback", "value_a"},
		{"env empty", "TEST_KEY_B", "", "fallback", "fallback"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envVal != "" {
				os.Setenv(tc.envKey, tc.envVal)
				defer os.Unsetenv(tc.envKey)
			}
			got := getEnv(tc.envKey, tc.fallback)
			if got != tc.expected {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tc.envKey, tc.fallback, got, tc.expected)
			}
		})
	}
}
