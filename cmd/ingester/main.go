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
const runFinalizationTimeout = 15 * time.Second

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
	runRepository := postgres.NewIngestionRunRepository(database)

	run, err := runRepository.Create(ctx, cfg.sourceURL)
	if err != nil {
		log.Fatalf("create ingestion run: %v", err)
	}

	log.Printf("ingestion started run_id=%d source=%s", run.ID, cfg.sourceURL)

	stats, err := service.IngestWithStats(ctx)
	if err != nil {
		markFailedErr := markRunFailed(runRepository, run.ID, stats, err)
		if markFailedErr != nil {
			log.Printf("mark ingestion run failed run_id=%d err=%v", run.ID, markFailedErr)
		}
		log.Fatalf("ingestion failed run_id=%d rows_processed=%d rows_inserted=%d rows_updated=%d err=%v", run.ID, stats.RowsProcessed, stats.RowsInserted, stats.RowsUpdated, err)
	}

	err = markRunSuccess(runRepository, run.ID, stats)
	if err != nil {
		log.Fatalf("mark ingestion run success run_id=%d: %v", run.ID, err)
	}

	log.Printf("ingestion completed run_id=%d rows_processed=%d rows_inserted=%d rows_updated=%d", run.ID, stats.RowsProcessed, stats.RowsInserted, stats.RowsUpdated)
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

type ingestionRunFinalizer interface {
	MarkSuccess(ctx context.Context, runID int64, rowsProcessed int, rowsInserted int, rowsUpdated int) error
	MarkFailed(ctx context.Context, runID int64, rowsProcessed int, rowsInserted int, rowsUpdated int, errorMessage string) error
}

func markRunFailed(runRepository ingestionRunFinalizer, runID int64, stats usecase.IngestionStats, ingestionErr error) error {
	return finalizeRun(func(ctx context.Context) error {
		return runRepository.MarkFailed(ctx, runID, stats.RowsProcessed, stats.RowsInserted, stats.RowsUpdated, ingestionErr.Error())
	})
}

func markRunSuccess(runRepository ingestionRunFinalizer, runID int64, stats usecase.IngestionStats) error {
	return finalizeRun(func(ctx context.Context) error {
		return runRepository.MarkSuccess(ctx, runID, stats.RowsProcessed, stats.RowsInserted, stats.RowsUpdated)
	})
}

func finalizeRun(action func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), runFinalizationTimeout)
	defer cancel()

	return action(ctx)
}
