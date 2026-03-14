package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"football-analytics/internal/infra/postgres"
	"football-analytics/internal/infra/sources"
	"football-analytics/internal/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultFootballDataURL = "https://www.football-data.co.uk/mmz4281/2425/E0.csv"

type config struct {
	postgresHost     string
	postgresPort     string
	postgresDB       string
	postgresUser     string
	postgresPassword string
	postgresSSLMode  string
	sourceURL        string
}

func main() {
	cfg := loadConfig()

	database, err := sql.Open("pgx", buildPostgresDSN(cfg))
	if err != nil {
		log.Fatalf("open postgres connection: %v", err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	reader := sources.NewFootballDataSource(httpClient, cfg.sourceURL)

	service := usecase.NewIngestionService(
		reader,
		postgres.NewIngestionMatchRepository(database),
	)

	err = service.Ingest(ctx)
	if err != nil {
		log.Fatalf("run ingestion: %v", err)
	}

	log.Printf("ingestion completed")
}

func loadConfig() config {
	return config{
		postgresHost:     getenvDefault("POSTGRES_HOST", "localhost"),
		postgresPort:     getenvDefault("POSTGRES_PORT", "5432"),
		postgresDB:       getenvDefault("POSTGRES_DB", "football_analytics"),
		postgresUser:     getenvDefault("POSTGRES_USER", "postgres"),
		postgresPassword: getenvDefault("POSTGRES_PASSWORD", "postgres"),
		postgresSSLMode:  getenvDefault("POSTGRES_SSLMODE", "disable"),
		sourceURL:        getenvDefault("INGESTION_SOURCE_URL", defaultFootballDataURL),
	}
}

func buildPostgresDSN(cfg config) string {
	query := url.Values{}
	query.Set("sslmode", cfg.postgresSSLMode)

	connectionURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.postgresUser, cfg.postgresPassword),
		Host:     fmt.Sprintf("%s:%s", cfg.postgresHost, cfg.postgresPort),
		Path:     cfg.postgresDB,
		RawQuery: query.Encode(),
	}

	return connectionURL.String()
}

func getenvDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
