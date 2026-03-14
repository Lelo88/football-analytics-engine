package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	deliveryhttp "football-analytics/internal/delivery/http"
	"football-analytics/internal/infra/postgres"
	"football-analytics/internal/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type config struct {
	apiPort          string
	postgresHost     string
	postgresPort     string
	postgresDB       string
	postgresUser     string
	postgresPassword string
	postgresSSLMode  string
}

func main() {
	cfg := loadConfig()

	database, err := sql.Open("pgx", buildPostgresDSN(cfg))
	if err != nil {
		log.Fatalf("open postgres connection: %v", err)
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	queryRepository := postgres.NewTeamAnalyticsReadRepository(database)
	analyticsService := usecase.NewTeamAnalyticsService(queryRepository)
	handler := deliveryhttp.NewHandler(analyticsService)

	address := fmt.Sprintf(":%s", cfg.apiPort)
	log.Printf("api listening on %s", address)
	if err = http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("run api server: %v", err)
	}
}

func loadConfig() config {
	return config{
		apiPort:          getenvDefault("API_PORT", "8080"),
		postgresHost:     getenvDefault("POSTGRES_HOST", "localhost"),
		postgresPort:     getenvDefault("POSTGRES_PORT", "5432"),
		postgresDB:       getenvDefault("POSTGRES_DB", "football_analytics"),
		postgresUser:     getenvDefault("POSTGRES_USER", "postgres"),
		postgresPassword: getenvDefault("POSTGRES_PASSWORD", "postgres"),
		postgresSSLMode:  getenvDefault("POSTGRES_SSLMODE", "disable"),
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
