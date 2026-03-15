package usecase

import (
	"context"
	"fmt"
	"strings"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

type Reader interface {
	Fetch(ctx context.Context) ([]domain.IngestionMatch, error)
}

type MatchRepository interface {
	Save(ctx context.Context, match domain.IngestionMatch) (ports.IngestionSaveResult, error)
}

type IngestionStats struct {
	RowsProcessed int
	RowsInserted  int
	RowsUpdated   int
}

type IngestionService struct {
	reader     Reader
	repository MatchRepository
}

func NewIngestionService(
	reader Reader,
	repository MatchRepository,
) *IngestionService {
	return &IngestionService{
		reader:     reader,
		repository: repository,
	}
}

func (service *IngestionService) Ingest(ctx context.Context) error {
	_, err := service.IngestWithStats(ctx)
	return err
}

func (service *IngestionService) IngestWithStats(ctx context.Context) (IngestionStats, error) {
	stats := IngestionStats{}

	matches, err := service.reader.Fetch(ctx)
	if err != nil {
		return stats, fmt.Errorf("fetch matches: %w", err)
	}

	for index, match := range matches {
		stats.RowsProcessed++

		err := validateMatch(match)
		if err != nil {
			return stats, fmt.Errorf("validate match %d: %w", index, err)
		}

		saveResult, err := service.repository.Save(ctx, match)
		if err != nil {
			return stats, fmt.Errorf("save match %d: %w", index, err)
		}

		if saveResult.Inserted {
			stats.RowsInserted++
		}
		if saveResult.Updated {
			stats.RowsUpdated++
		}
	}

	return stats, nil
}

func validateMatch(match domain.IngestionMatch) error {
	if strings.TrimSpace(match.CompetitionName) == "" {
		return fmt.Errorf("competition name is required")
	}
	if strings.TrimSpace(match.SeasonLabel) == "" {
		return fmt.Errorf("season label is required")
	}
	if match.MatchDate.IsZero() {
		return fmt.Errorf("match date is required")
	}
	if strings.TrimSpace(match.HomeTeamName) == "" {
		return fmt.Errorf("home team is required")
	}
	if strings.TrimSpace(match.AwayTeamName) == "" {
		return fmt.Errorf("away team is required")
	}
	if strings.EqualFold(strings.TrimSpace(match.HomeTeamName), strings.TrimSpace(match.AwayTeamName)) {
		return fmt.Errorf("home and away teams must be different")
	}

	return nil
}
