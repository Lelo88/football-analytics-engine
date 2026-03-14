package main

import "testing"

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "")
	t.Setenv("POSTGRES_PORT", "")
	t.Setenv("POSTGRES_DB", "")
	t.Setenv("POSTGRES_USER", "")
	t.Setenv("POSTGRES_PASSWORD", "")
	t.Setenv("POSTGRES_SSLMODE", "")
	t.Setenv("INGESTION_SOURCE_URL", "")

	cfg := loadConfig()

	if cfg.postgresHost != "localhost" {
		t.Fatalf("expected default host localhost, got %q", cfg.postgresHost)
	}
	if cfg.postgresPort != "5432" {
		t.Fatalf("expected default port 5432, got %q", cfg.postgresPort)
	}
	if cfg.postgresDB != "football_analytics" {
		t.Fatalf("expected default db football_analytics, got %q", cfg.postgresDB)
	}
	if cfg.postgresUser != "postgres" {
		t.Fatalf("expected default user postgres, got %q", cfg.postgresUser)
	}
	if cfg.postgresPassword != "postgres" {
		t.Fatalf("expected default password postgres, got %q", cfg.postgresPassword)
	}
	if cfg.postgresSSLMode != "disable" {
		t.Fatalf("expected default sslmode disable, got %q", cfg.postgresSSLMode)
	}
	if cfg.sourceURL != defaultFootballDataURL {
		t.Fatalf("expected default source url %q, got %q", defaultFootballDataURL, cfg.sourceURL)
	}
}

func TestLoadConfigOverrides(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "db.internal")
	t.Setenv("POSTGRES_PORT", "5544")
	t.Setenv("POSTGRES_DB", "analytics")
	t.Setenv("POSTGRES_USER", "fa_user")
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("POSTGRES_SSLMODE", "require")
	t.Setenv("INGESTION_SOURCE_URL", "https://example.test/source.csv")

	cfg := loadConfig()

	if cfg.postgresHost != "db.internal" || cfg.postgresPort != "5544" || cfg.postgresDB != "analytics" || cfg.postgresUser != "fa_user" || cfg.postgresPassword != "secret" || cfg.postgresSSLMode != "require" || cfg.sourceURL != "https://example.test/source.csv" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestBuildPostgresDSN(t *testing.T) {
	t.Parallel()

	dsn := buildPostgresDSN(config{
		postgresHost:     "db.internal",
		postgresPort:     "5544",
		postgresDB:       "analytics",
		postgresUser:     "fa_user",
		postgresPassword: "secret",
		postgresSSLMode:  "require",
	})

	expected := "postgres://fa_user:secret@db.internal:5544/analytics?sslmode=require"
	if dsn != expected {
		t.Fatalf("expected DSN %q, got %q", expected, dsn)
	}
}
