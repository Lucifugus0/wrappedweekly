package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	JWTExpiry       time.Duration
	AIProvider      string
	AppEnv          string
	CookieSecure    bool
	CookieDomain    string
	// CookieCrossSite=true sends SameSite=None (required when frontend and
	// backend are on different origins/domains, e.g. Vercel + Render).
	// Browsers reject SameSite=None without Secure, so this always implies
	// CookieSecure=true regardless of the COOKIE_SECURE value.
	// Default false -> SameSite=Lax, for same-origin setups (e.g. behind Nginx).
	CookieCrossSite bool
	FrontendBaseURL string
}

func Load() Config {
	crossSite := getEnv("COOKIE_CROSS_SITE", "false") == "true"
	return Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		JWTSecret:       getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiry:       parseDurationHours(getEnv("JWT_EXPIRY_HOURS", "24")),
		AIProvider:      getEnv("AI_PROVIDER", "mock"),
		AppEnv:          getEnv("APP_ENV", "development"),
		CookieSecure:    getEnv("COOKIE_SECURE", "false") == "true" || crossSite,
		CookieDomain:    getEnv("COOKIE_DOMAIN", ""),
		CookieCrossSite: crossSite,
		FrontendBaseURL: getEnv("FRONTEND_BASE_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDurationHours(hoursStr string) time.Duration {
	hours := 24
	if hoursStr != "" {
		if h, err := time.ParseDuration(hoursStr + "h"); err == nil {
			return h
		}
	}
	return time.Duration(hours) * time.Hour
}
